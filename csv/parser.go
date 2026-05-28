package csvparser

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"tally/models"
	"time"
)

// specific struct as db type
type transaction struct {
	date        time.Time
	description string
	amount      float64
}

func buildColumnMap(bank string, header []string) map[int]models.ColumnType {

	validBankColumns := models.BankConfigs[bank]
	columnMap := make(map[int]models.ColumnType, len(header))

	for idx, value := range header {
		if colType, ok := validBankColumns[strings.ToLower(value)]; ok {
			columnMap[idx] = colType
		}
	}
	return columnMap
}

func parseRow(row []string, columnMap map[int]models.ColumnType) (transaction, bool, error) {

	var entry transaction

	for idx, value := range row {
		switch columnMap[idx] {
		case models.ColumnDate:
			date, err := time.Parse("01/02/2006", value)
			if err != nil {
				return entry, false, fmt.Errorf("invalid date %q: %w", value, err)
			}

			entry.date = date
		case models.ColumnName:
			entry.description = value
		case models.ColumnAmount:
			amount, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return entry, false, fmt.Errorf("invalid amount %q: %w", value, err)
			}
			if amount > 0 {
				return entry, true, nil // skip credits
			}
			entry.amount = -amount
		}
	}

	return entry, false, nil
}

func ImportTransactions(db *sql.DB, bank string, reader *csv.Reader) error {

	// get csv header
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("empty or unreadable CSV: %w", err)
	}

	columnMap := buildColumnMap(bank, header)

	// batch process all inserts into db
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO transactions (date, description, amount, card) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

outer:
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		entry, skip, err := parseRow(row, columnMap)
		if err != nil {
			return err
		}
		if skip {
			continue outer
		}

		_, err = stmt.Exec(entry.date, entry.description, entry.amount, bank)
		if err != nil {
			return fmt.Errorf("failed to insert row: %w", err)
		}
	}

	return tx.Commit()

}
