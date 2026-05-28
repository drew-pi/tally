package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type ColumnType string

const (
    Date        ColumnType = "date"
    Description ColumnType = "description"
    Amount      ColumnType = "amount"
)

type BankConfig map[string]ColumnType

type CsvConfig map[string]BankConfig

var csvConfig = CsvConfig{
    "WF": {
        "date":        Date,
        "description": Description,
        "amount":      Amount,
    },
    "Fidelity": {
        "date":   Date,
        "name":   Description,
        "amount": Amount,
    },
}


func printTable(list []int) {
	var header, values strings.Builder

	for idx, value := range list {
		fmt.Fprintf(&header, "%d | ", idx)
		fmt.Fprintf(&values, "%d | ", value)
	}

	splitter := strings.Repeat("-", header.Len()-1)

	fmt.Printf("%s\n%s\n%s\n",header.String(), splitter, values.String())
}

func printCsvTable(reader *csv.Reader) {

	var table strings.Builder

	header, header_err := reader.Read()
	if header_err != nil { // catches EOF too
		fmt.Println("Empty reader")
		return
	}

	fmt_header := strings.Join(header, " | ")
	splitter := strings.Repeat("-", len(fmt_header))

	for {
		row, err := reader.Read()
		if err != nil { // catches EOF too
			break
		}

		fmt_row := strings.Join(row, " | ")
		fmt.Fprintf(&table, "%s\n", fmt_row)
	}

	fmt.Printf("%s\n%s\n%s\n", fmt_header, splitter, table.String())

	for {
		row, err := reader.Read()
		if err != nil { // catches EOF too
			break
		}
		fmt.Println(row)
	}	
}

func printCleanTable(reader *csv.Reader) {

	valid_column_names := map[string]bool{
		"date":        true,
		"description": true,
		"name":        true,
		"amount":      true,
	}

	valid_column_idx := make(map[int]bool)

	header, header_err := reader.Read()
	if header_err != nil {
		fmt.Println("Empty reader")
		return
	}

	clean_header := make([]string, 0, len(header))

	for idx, value := range header {
		if _, ok := valid_column_names[strings.ToLower(value)]; ok {
			valid_column_idx[idx] = true
			clean_header = append(clean_header, value)
		}
	}

	fmt_header := strings.Join(clean_header, " | ")
	splitter := strings.Repeat("-", len(fmt_header)-1)

	clean_row := make([]string, 0, len(valid_column_idx))

	var table strings.Builder

	for {
		row, err := reader.Read()
		if err != nil { // catches EOF too
			break
		}

		clean_row = clean_row[:0]
		for idx, value := range row {
			if valid_column_idx[idx] {
				clean_row = append(clean_row, value)
			}
		}

		fmt_row := strings.Join(clean_row, " | ")
		fmt.Fprintf(&table, "%s\n", fmt_row)
	}

	fmt.Printf("%s\n%s\n%s\n", fmt_header, splitter, table.String())
}


func printCleanTableWithConfig(card string, reader *csv.Reader) {

	valid_column_names := csvConfig[card]

	valid_column_idx := make(map[int]bool)

	header, header_err := reader.Read()
	if header_err != nil {
		fmt.Println("Empty reader")
		return
	}

	clean_header := make([]string, 0, len(header))

	for idx, value := range header {
		clean_value, ok := valid_column_names[strings.ToLower(value)]
		if ok {
			valid_column_idx[idx] = true
			clean_header = append(clean_header, string(clean_value))
		}
	}

	fmt_header := strings.Join(clean_header, " | ")
	splitter := strings.Repeat("-", len(fmt_header)-1)

	clean_row := make([]string, 0, len(valid_column_idx))

	var table strings.Builder

	for {
		row, err := reader.Read()
		if err != nil { // catches EOF too
			break
		}

		clean_row = clean_row[:0]
		for idx, value := range row {
			if valid_column_idx[idx] {
				clean_row = append(clean_row, value)
			}
		}

		fmt_row := strings.Join(clean_row, " | ")
		fmt.Fprintf(&table, "%s\n", fmt_row)
	}

	fmt.Printf("%s\n%s\n%s\n", fmt_header, splitter, table.String())
}


func main() {
	fmt.Println("Hello, World!")

	card := "Fidelity"
	file_name := card + "_04:2026.csv"

	file, err := os.Open(file_name)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
	defer file.Close()

	bufferedFile := bufio.NewReader(file)
	reader := csv.NewReader(bufferedFile)

	// printCsvTable(reader)

	// printCleanTable(reader)

	printCleanTableWithConfig(card, reader)




    






}