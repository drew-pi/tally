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
	return db.Query("SELECT id, date, description, amount, card FROM transactions ORDER BY date DESC")
}

func GetTransactionById(db *sql.DB, id int) (*sql.Rows, error) {
	return db.Query("SELECT id, date, description, amount, card FROM transactions WHERE id = $1 ORDER BY date DESC;", id)
}

func GetTransactionsByAmount(db *sql.DB, amount float64) (*sql.Rows, error) {
	return db.Query("SELECT id, date, description, amount, card  FROM transactions  WHERE amount = $1 ORDER BY date DESC;", amount)
}

func GetTransactionsByAmountRange(db *sql.DB, amountFrom, amountTo float64) (*sql.Rows, error) {
	return db.Query("SELECT id, date, description, amount, card  FROM transactions  WHERE amount BETWEEN $1 AND $2 ORDER BY date DESC;", amountFrom, amountTo)
}

func GetTransactionsByDateRange(db *sql.DB, dateFrom, dateTo time.Time) (*sql.Rows, error) {
	return db.Query("SELECT id, date, description, amount, card  FROM transactions  WHERE date BETWEEN $1 AND $2 ORDER BY date DESC;", dateFrom, dateTo)
}

func GetTransactionsByBank(db *sql.DB, card string) (*sql.Rows, error) {
	return db.Query("SELECT id, date, description, amount, card  FROM transactions  WHERE amount = $1;", card)
}

func GetTransactions(
	db *sql.DB,
	id sql.NullInt64,
	amount sql.NullFloat64,
	dateFrom sql.NullString,
	dateTo sql.NullString,
	card sql.NullString,
) (*sql.Rows, error) {
	return db.Query(`
	SELECT id, date, description, amount, card
		FROM transactions  
		WHERE ($1::int IS NULL OR id = $1) 
		AND ($2::numeric IS NULL OR amount = $2) 
		AND ($3::date IS NULL OR date >= $3) 
		AND ($4::date IS NULL OR date <= $4) 
		AND ($5::card_type IS NULL OR card = $5) 
		ORDER BY date DESC;
	`, id, amount, dateFrom, dateTo, card)
}
