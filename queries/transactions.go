package queries

import (
	"database/sql"
	"tally/models"
	"time"
)

func InsertTransaction(db *sql.DB, date time.Time, vendor, description, category string, amount float64, paymentMethodID int) error {
	_, err := db.Exec(
		"INSERT INTO transactions (date, vendor, description, category, amount, payment_method_id) VALUES ($1, $2, $3, $4, $5, $6)",
		date, vendor, description, category, amount, paymentMethodID,
	)
	return err
}

func GetAllTransactions(db *sql.DB) (*sql.Rows, error) {
	return db.Query("SELECT id, date, vendor, description, category, amount, payment_method_id FROM transactions ORDER BY date DESC")
}

func GetTransactionById(db *sql.DB, id int) (*sql.Rows, error) {
	return db.Query("SELECT id, date, vendor, description, category, amount, payment_method_id FROM transactions WHERE id = $1", id)
}

func GetTransactionsByAmount(db *sql.DB, amount float64) (*sql.Rows, error) {
	return db.Query("SELECT id, date, vendor, description, category, amount, payment_method_id FROM transactions WHERE amount = $1 ORDER BY date DESC", amount)
}

func GetTransactionsByAmountRange(db *sql.DB, amountFrom, amountTo float64) (*sql.Rows, error) {
	return db.Query("SELECT id, date, vendor, description, category, amount, payment_method_id FROM transactions WHERE amount BETWEEN $1 AND $2 ORDER BY date DESC", amountFrom, amountTo)
}

func GetTransactionsByDateRange(db *sql.DB, dateFrom, dateTo time.Time) (*sql.Rows, error) {
	return db.Query("SELECT id, date, vendor, description, category, amount, payment_method_id FROM transactions WHERE date BETWEEN $1 AND $2 ORDER BY date DESC", dateFrom, dateTo)
}

func GetTransactionsByBank(db *sql.DB, bank models.BankName) (*sql.Rows, error) {
	return db.Query(`
		SELECT id, date, vendor, description, category, amount, payment_method_id
		FROM transactions
		WHERE payment_method_id = (
			SELECT id FROM payment_methods
			WHERE bank_id = (SELECT id FROM banks WHERE name = $1)
		)
		ORDER BY date DESC
	`, string(bank))
}

func GetTransactions(
	db *sql.DB,
	id sql.NullInt64,
	amountFrom sql.NullFloat64,
	amountTo sql.NullFloat64,
	dateFrom sql.NullTime,
	dateTo sql.NullTime,
	bank sql.NullString,
) (*sql.Rows, error) {
	return db.Query(`
		SELECT id, date, vendor, description, category, amount, payment_method_id
		FROM transactions
		WHERE ($1::int IS NULL OR id = $1)
		AND ($2::numeric IS NULL OR amount >= $2)
		AND ($3::numeric IS NULL OR amount <= $3)
		AND ($4::date IS NULL OR date >= $4)
		AND ($5::date IS NULL OR date <= $5)
		AND ($6::text IS NULL OR payment_method_id = (
			SELECT id FROM payment_methods
			WHERE bank_id = (SELECT id FROM banks WHERE name = $6)
		))
		ORDER BY date DESC
	`, id, amountFrom, amountTo, dateFrom, dateTo, bank)
}
