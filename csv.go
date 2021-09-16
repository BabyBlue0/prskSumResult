package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func writeCSV(filePath string, recodes []PRSKOutputFormatToCSV) error {

	if len(recodes) <= 0 {
		return nil
	}

	titles := recodes[0].Titles()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	cw := csv.NewWriter(file)
	defer cw.Flush()

	//write title
	cw.Write(titles)

	for _, ofcsv := range recodes {
		ofmap := ofcsv.ToMap()
		rec := []string{}
		for _, t := range titles {
			rec = append(rec, ofmap[t])
		}

		if err := cw.Write(rec); err != nil {
			fmt.Println(err)
			break
		}
	}

	return nil
}
