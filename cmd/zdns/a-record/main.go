package main

// before run this program, you should build zdns and put the binary file in the same directory as this program.

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/mattn/go-pipeline"
	"github.com/yagikota/danewebperf/utils"
)

const (
	defaultInputCSV  = "../../../dataset/hall-of-flame/hall-of-flame-websites-tlsa-usage3.csv"
	defaultOutputCSV = "../../../dataset/hall-of-flame/hall-of-flame-websites-tlsa-usage3-ip-address.csv"
)

type RecordSlice []Record

type Record struct {
	Domain   string
	Response DNSResponse
}

type DNSResponse struct {
	Data      Data   `json:"data"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

type Data struct {
	Additionals []Additional `json:"additionals"`
	Answers     []AAnswer    `json:"answers"`
	Protocol    string       `json:"protocol"`
	Resolver    string       `json:"resolver"`
}

type Additional struct {
	Flags   string `json:"flags"`
	Type    string `json:"type"`
	UDPSize int    `json:"udpsize"`
	Version int    `json:"version"`
}

type AAnswer struct {
	Answer string `json:"answer"`
	Class  string `json:"class"`
	Name   string `json:"name"`
	TTL    int    `json:"ttl"`
	Type   string `json:"type"`
}

func zdns(domain string, rrType string) ([]byte, error) {
	return pipeline.Output(
		[]string{"echo", domain},
		[]string{"../zdns", rrType, "--conf-file", "../resolv.conf"},
	)
}

func zdnsA(domain string) ([]byte, error) {
	return zdns(domain, "A")
}

func writeCSV(filePath string, recordSlice RecordSlice) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"domain", "ip"}
	writer.Write(header)

	for _, record := range recordSlice {
		if len(record.Response.Data.Answers) == 0 {
			continue
		}
		firstResponse := record.Response.Data.Answers[0]
		if firstResponse.Type != "A" {
			continue
		}
		record := []string{
			record.Domain,
			firstResponse.Answer,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	start := time.Now()
	log.Println("start time: ", start.Format("2006-01-02-15-04-05"))

	first := flag.Int("first", 1, "first index of Domain list")
	last := flag.Int("last", -1, "last index of Domain list. if -1, last index is last index of Domain list")
	inputCSV := flag.String("inputCSV", defaultInputCSV, "input CSV path")
	outputCSV := flag.String("outputCSV", defaultOutputCSV, "output CSV path")
	concurrency := flag.Int("concurrency", 10, "concurrency")
	flag.Parse()

	// domainList, err := utils.ReadDomainListCSV(*inputCSV)
	domainList, err := utils.ReadDomainListCSV(*inputCSV)
	if err != nil {
		log.Fatalln(err)
	}

	if *last == -1 {
		*last = len(domainList)
	}
	subsetDomainList := domainList[*first-1 : *last]

	var wg sync.WaitGroup
	sem := make(chan struct{}, *concurrency)
	recordChannel := make(chan Record, len(subsetDomainList))
	for index, record := range subsetDomainList {
		sem <- struct{}{}
		wg.Add(1)

		go func(index int, record utils.Record) {
			defer func() {
				<-sem
				wg.Done()
			}()

			response, err := zdnsA(record.Domain)
			if err != nil {
				log.Println(err)
			}

			var dnsResponse DNSResponse
			if err := json.Unmarshal(response, &dnsResponse); err != nil {
				log.Printf("json unmarshal error: %v", err)
			}

			recordChannel <- Record{record.Domain, dnsResponse}
			log.Printf("index: %d, domain: %s", index+1, subsetDomainList[index].Domain)
		}(index, record)
	}

	go func() {
		wg.Wait()
		close(recordChannel)
	}()

	var recordSlice RecordSlice
	for record := range recordChannel {
		recordSlice = append(recordSlice, record)
	}

	// sort recordSlice by domain
	sort.Slice(recordSlice, func(i, j int) bool {
		return recordSlice[i].Domain < recordSlice[j].Domain
	})

	if err := writeCSV(*outputCSV, recordSlice); err != nil {
		log.Printf("write csv error: %v", err)
	}

	log.Println("finish!")
	log.Println("elapsed time: ", time.Since(start))
}
