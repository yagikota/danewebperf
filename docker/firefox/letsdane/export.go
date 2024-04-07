package letsdane

import (
	"encoding/csv"
	"io"
	"log"
)

type DANEValidationResultSlice []DANEValidationResult

type DANEValidationResult struct {
	Host            string
	IsDANEValidated bool
	Err             error
}

func ExportAsCSV(results DANEValidationResultSlice, w io.Writer) error {

	log.Println("results: ", results)

	if len(results) == 0 {
		log.Println("No results to export")
		return nil
	}

	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	if err := csvWriter.Write([]string{"Host", "DANEValidated", "Error"}); err != nil {
		return err
	}

	log.Println(len(results), "results to export")

	for _, result := range results {
		var daneValidated string
		if result.IsDANEValidated {
			daneValidated = "true"
		} else {
			daneValidated = "false"
		}

		if result.Err != nil {
			if err := csvWriter.Write([]string{result.Host, daneValidated, result.Err.Error()}); err != nil {
				return err
			}
			continue
		}

		if err := csvWriter.Write([]string{result.Host, daneValidated, ""}); err != nil {
			return err
		}

	}

	return nil
}
