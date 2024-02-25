/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package process

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func readCsvFile(filePath string) (fields [][]string, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("reading failed: %w", err)
		}
	}()

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		err = fmt.Errorf("unable to get file information for %s %w", filePath, err)
		return nil, err
	}
	log.Println("File Name:", fileInfo.Name())
	//log.Println("Size ", fileInfo.Size(), " bytes")
	//log.Println("Permissions:", fileInfo.Mode())
	//log.Println("Last modified:", fileInfo.ModTime())
	//log.Println("Is Directory: ", fileInfo.IsDir())

	f, err := os.Open(filePath)
	if err != nil {
		err = fmt.Errorf("unable to open file %s %w", filePath, err)
		return nil, err
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		err = fmt.Errorf("unable to parse file as CSV for %s %w", filePath, err)
		return nil, err
	}

	return records, err
}

func parseRecords(records [][]string) (langMap map[string][]string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("parsing records failed: %w", err)
		}
	}()

	// First we want to get the languages defined
	// If the cell is empy, we use it to identify empty rows
	languages := records[0]

	retMap := make(map[string][]string, len(records))

	//Get rid of leading and trailing spaces in the language names, doing so easier identifies empty spaces
	for i := 0; i < len(languages); i++ {
		languages[i] = strings.TrimSpace(languages[i])
	}

	for langidx, lang := range languages {

		// Array of strings for the map later
		lang_strings := make([]string, len(records)-1)

		if lang == "" {
			fmt.Println("Field", langidx+1, "in CSV has no language name, skipping column!")
			continue
		}

		for idx, row := range records[1:] {
			lang_strings[idx] = row[langidx]
		}

		// Put it in the map
		retMap[lang] = lang_strings
	}

	langMap = retMap

	return langMap, err
}

func Csv(path string) (langMap map[string][]string, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("CSV Processing failed: %w", err)
		}
	}()

	start := time.Now()

	fileExtension := filepath.Ext(path)
	if fileExtension != ".csv" {
		err = errors.New("file must be a CSV file (.csv)")
		return nil, err
	}

	//Read the CSV file
	rows, err := readCsvFile(path)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		err = errors.New("CSV file is empty")
		return nil, err
	}

	if rows[0][0] != "var" {
		err = errors.New("first field of CSV file must be \"var\"")
		return nil, err
	}

	langMap, err = parseRecords(rows)
	if err != nil {
		return nil, err
	}

	elapsed := time.Since(start)
	log.Printf("process.Csv took %s", elapsed)

	return langMap, err
}
