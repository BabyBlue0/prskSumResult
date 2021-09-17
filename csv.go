package main

import (
	"encoding/csv"
	"os"
)

func writeCSV(filePath string, records []PRSKOutputFormatToCSV) error {

	if len(records) <= 0 {
		return nil
	}

	titles := records[0].Titles()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	cw := csv.NewWriter(file)
	defer cw.Flush()

	//write title
	cw.Write(titles)

	for _, ofcsv := range records {
		ofmap := ofcsv.ToMap()
		rec := []string{}
		for _, t := range titles {
			rec = append(rec, ofmap[t])
		}

		if err := cw.Write(rec); err != nil {
			return err
		}
	}

	return nil
}
