package letsdane

import (
	"encoding/csv"
	"io"
)

var DANEValidationResultsChan = make(chan DANEValidationResult, 1)

type DANEValidationResult struct {
	Host            string
	IsDANEValidated bool
	Err             error
}

func writeEachLine(line []string, w io.Writer) error {
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()
	if err := csvWriter.Write(line); err != nil {
		return err
	}
	return nil
}

func WriteToCSV(resultChan <-chan DANEValidationResult, w io.Writer) error {

	if err := writeEachLine([]string{"Host", "DANE Validated", "Error"}, w); err != nil {
		return err
	}

	var daneValidated string
	// log.Println(<-resultChan)
	for result := range resultChan {
		if result.IsDANEValidated {
			daneValidated = "true"
		} else {
			daneValidated = "false"
		}

		if result.Err != nil {
			if err := writeEachLine([]string{result.Host, daneValidated, result.Err.Error()}, w); err != nil {
				return err
			}
			continue
		}

		if err := writeEachLine([]string{result.Host, daneValidated, ""}, w); err != nil {
			return err
		}
	}

	return nil
}
