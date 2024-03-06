package utils

import (
	"encoding/csv"
	"os"
	"sort"
)

type DomainList []Record

type Record struct {
	Domain string
}

func ReadDomainListCSV(path string) (DomainList, error) {
	var domainList DomainList
	file, err := os.Open(path)
	if err != nil {
		return domainList, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return domainList, err
	}
	for _, row := range rows {
		domainList = append(domainList, Record{
			Domain: row[0],
		})
	}
	return domainList, nil
}

func WriteDomainListCSV(path string, domainList DomainList) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range domainList {
		record := []string{record.Domain}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func WritePageLoadTimeCSV(path string, domainPageLoadMap map[string]string, cache, dane bool) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"domain", "pageLoadTime", "cache", "dane"}); err != nil {
		return err
	}

	domains := make([]string, 0, len(domainPageLoadMap))
	for domain := range domainPageLoadMap {
		domains = append(domains, domain)
	}
	sort.Slice(domains, func(i, j int) bool {
		return domains[i] < domains[j]
	})

	for _, domain := range domains {
		record := []string{domain, domainPageLoadMap[domain], "false", "false"}
		if cache {
			record[2] = "true"
		}
		if dane {
			record[3] = "true"
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
