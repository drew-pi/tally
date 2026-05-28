package queries

import (
	"database/sql"
	"fmt"
	"tally/models"
	"time"
)

func InsertTransaction(db *sql.DB, date time.Time, description string, amount float64, bank string) error {

	// make sure that this bank is valid
	if _, ok := models.BankConfigs[bank]; !ok {
		return fmt.Errorf("bank %q not found in config", bank)
	}
	_, err := db.Exec(
		"INSERT INTO transactions (date, description, amount, card) VALUES ($1, $2, $3, $4)",
		date, description, amount, bank,
	)
	return err
}

func GetAllTransactions(db *sql.DB) (*sql.Rows, error) {
	return db.Query("SELECT date, description, amount, card FROM transactions")
}
