package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	logger      *slog.Logger
	concurrency = 10
)

const (
	defaultInputCSV  = "../../dataset/hall-of-flame/hall-of-flame-websites-tlsa-usage3-ip-address.csv"
	defaultOutputCSV = "../../dataset/hall-of-flame/hall-of-flame-websites-tlsa-usage3-ip-address-country-code.csv"
	// https://api.iplocation.net/
	apiURL = "https://api.iplocation.net/"
)

type Response struct {
	IP              string `json:"ip"`
	IPNumber        string `json:"ip_number"`
	IPVersion       int    `json:"ip_version"`
	CountryName     string `json:"country_name"`
	CountryCode2    string `json:"country_code2"`
	Isp             string `json:"isp"`
	ResponseCode    string `json:"response_code"`
	ResponseMessage string `json:"response_message"`
}

func getCountryCodeFromIP(ip string) ([]byte, error) {

	res, err := http.Get(apiURL + "?ip=" + ip)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

type DomainIPs []DomainIP

type DomainIP struct {
	Domain string
	IP     string
}

func readDomainIPCSV(filePath string) (DomainIPs, error) {
	var domainIPs DomainIPs
	file, err := os.Open(filePath)
	if err != nil {
		return domainIPs, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return domainIPs, err
	}
	for _, row := range rows[1:] {
		domainIPs = append(domainIPs, DomainIP{
			Domain: row[0],
			IP:     row[1],
		})
	}
	return domainIPs, nil
}

type Results []Result

type Result struct {
	Domain string
	IP     string
}

func writeResultsCSV(filePath string, results Results) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"domain", "country_code"}
	writer.Write(header)

	for _, result := range results {
		record := []string{
			result.Domain,
			result.IP,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	start := time.Now()
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	domainIPs, err := readDomainIPCSV(defaultInputCSV)
	if err != nil {
		logger.Error(err.Error())
	}

	results := make([]Result, len(domainIPs))
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)
	for idx, domainIP := range domainIPs {
		sem <- struct{}{}
		wg.Add(1)
		go func(idx int, domainIP DomainIP) {
			defer func() {
				<-sem
				wg.Done()
			}()

			buf, err := getCountryCodeFromIP(domainIP.IP)
			if err != nil {
				logger.Error(err.Error())
			}

			var response Response
			err = json.Unmarshal(buf, &response)
			if err != nil {
				logger.Error(err.Error())
			}

			results[idx] = Result{
				Domain: domainIP.Domain,
				IP:     response.CountryCode2,
			}
			logger.Info(fmt.Sprintf("index: %d, domain: %s, ip: %s", idx+1, domainIP.Domain, response.CountryCode2))

		}(idx, domainIP)
	}

	wg.Wait()

	if err := writeResultsCSV(defaultOutputCSV, results); err != nil {
		logger.Error(err.Error())
	}

	logger.Info(fmt.Sprintf("elapsed time: %v", time.Since(start)))

}
