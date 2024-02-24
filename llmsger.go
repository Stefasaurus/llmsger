package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func main() {

	start := time.Now()

	//const csvPath string = "sample2.csv"
	const csvPath string = "wrongfile.txt"

	records := readCsvFile(csvPath)

	if len(csvPath) < 4 {
		fmt.Println("File specified must be a CSV file (.csv).")
		os.Exit(1)
	}
	if fileExt := csvPath[len(csvPath)-4:]; fileExt != ".csv" {

		fmt.Println("File specified must be a CSV file (.csv).")
		os.Exit(1)

	}
	if len(records) == 0 {
		fmt.Println("File specified is empty or wrong type.")
		os.Exit(1)
	}
	if records[0][0] != "varname" {
		fmt.Printf("First cell of the file must be \"varname\", but it is: \"%s\".", records[0][0])
		os.Exit(1)
	}

	fmt.Println(records)

	for _, rows := range records {
		for _, fields := range rows {
			fmt.Println(fields)
		}
	}

	filepath := csvPath
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	fmt.Println(string(data))

	fmt.Println("Hello, World!")

	elapsed := time.Since(start)
	log.Printf("Program took %s", elapsed)
}
