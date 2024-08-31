package model

type LetsDANERecords []LetsDANERecord

type LetsDANERecord struct {
	Host          string
	DANEValidated string
	Error         string
}

func ConvertLetsDANERecords(records [][]string) LetsDANERecords {
	var r LetsDANERecords
	if len(records) == 0 {
		return r
	}

	for _, record := range records[1:] {

		r = append(r, LetsDANERecord{
			Host:          record[0],
			DANEValidated: record[1],
			Error:         record[2],
		})

	}
	return r
}

func FilterDANEValidated(records LetsDANERecords) LetsDANERecords {
	var r LetsDANERecords
	for _, record := range records {
		if record.DANEValidated == "true" {
			r = append(r, record)
		}
	}
	return r
}

type Records []HarRecord

type HarRecord struct {
	Status string
	Method string
	Domain string
}

func ConvertHarStruct(records [][]string) Records {
	var r Records
	if len(records) == 0 {
		return r
	}

	for _, record := range records[1:] {
		r = append(r, HarRecord{
			Status: record[0],
			Method: record[1],
			Domain: record[2],
		})
	}
	return r
}
