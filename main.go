package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"tally/db"
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

	fmt.Printf("%s\n%s\n%s\n", header.String(), splitter, values.String())
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

func insertCsvIntoTable(db *sql.DB, card string, reader *csv.Reader) {

	valid_column_names := csvConfig[card]

	valid_column_idx := make(map[int]bool)

	header, header_err := reader.Read()
	if header_err != nil {
		fmt.Println("Empty reader")
		return
	}

	for idx, value := range header {
		_, ok := valid_column_names[strings.ToLower(value)]
		if ok {
			valid_column_idx[idx] = true
			// clean_header = append(clean_header, string(clean_value))
		}
	}

	tx, err := db.Begin()
	if err != nil {
		return
	}

	stmt, err := tx.Prepare("INSERT INTO transactions (date, description, amount, card) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return
	}
	defer stmt.Close()

	clean_rows := make([]string, 0, len(valid_column_idx))

	for {
		row, err := reader.Read()
		if err != nil { // catches EOF too
			break
		}

		clean_rows = clean_rows[:0]
		for idx, value := range row {
			if valid_column_idx[idx] {
				clean_rows = append(clean_rows, value)
			}
		}

		_, err = stmt.Exec(clean_rows[0], clean_rows[1], clean_rows[2], card)
		if err != nil {
			tx.Rollback() // undo everything if one fails
			return
		}
	}

	tx.Commit()
}

func insertCsvIntoTableEfficient(db *sql.DB, card string, reader *csv.Reader) {

	// get csv header
	header, err := reader.Read()
	if err != nil {
		fmt.Println("Empty reader")
		return
	}

	// get the valid column names for this card
	valid_column_names := csvConfig[card]

	// create column map
	column_map := make(map[int]ColumnType)

	// populate column map to clean up csv columns
	for idx, value := range header {
		if colType, ok := valid_column_names[strings.ToLower(value)]; ok {
			column_map[idx] = colType
		}
	}

	// batch process all inserts into db
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO transactions (date, description, amount, card) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return
	}
	defer stmt.Close()

	// specific struct as db type
	type RowData struct {
		date        time.Time
		description string
		amount      float64
	}

outer:
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		var data RowData

		for idx, value := range row {
			switch column_map[idx] {
			case Date:
				date, err := time.Parse("01/02/2006", value)
				if err != nil {
					tx.Rollback()
					return
				}

				data.date = date
			case Description:
				data.description = value
			case Amount:
				amount, err := strconv.ParseFloat(value, 64)
				if err != nil {
					tx.Rollback()
					return
				}
				if amount > 0 {
					fmt.Println("Got a positive value", value)
					continue outer // skips to next row
				}
				data.amount = -amount
			}
		}

		_, err = stmt.Exec(data.date, data.description, data.amount, card)
		if err != nil {
			tx.Rollback()
			return
		}
	}

	tx.Commit() // nothing gets rolled back once commit is called
}

func main() {
	fmt.Println("Hello, World!")

	card := "WF"
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

	// printCleanTableWithConfig(card, reader)

	database := db.Connect()
	defer database.Close()

	_, err = database.Exec(`
    	CREATE TYPE card_type AS ENUM ('WF', 'Fidelity')
	`)

	_, err = database.Exec(`
		CREATE TABLE IF NOT EXISTS transactions (
			id          SERIAL PRIMARY KEY,
			date        DATE,
			description TEXT,
			amount      NUMERIC(10, 2),
			card        card_type
		)
	`)

	insertCsvIntoTableEfficient(database, card, reader)

	date, err := time.Parse("01/02/2006", "05/01/2026") // MM/DD/YYYY format
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}

	_, err = database.Exec(
		"INSERT INTO transactions (date, description, amount, card) VALUES ($1, $2, $3, $4)",
		time.Time(date), "Grocery Store", float64(-52.30), "Fidelity",
	)

	rows, err := database.Query("SELECT date, description, amount, card FROM transactions")
	if err != nil {
		fmt.Println("Error querying db:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var date, description, amount, card string
		err := rows.Scan(&date, &description, &amount, &card)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		fmt.Println(date, description, amount, card)
	}
}
