package queries

import (
	"database/sql"
	"fmt"
	"tally/models"
)

func GetAllBanks(db *sql.DB) (*sql.Rows, error) {
	return db.Query("SELECT id, name FROM banks ORDER BY name")
}

func InsertBank(db *sql.DB, name models.BankName) (int, error) {
	var id int
	err := db.QueryRow(
		"INSERT INTO banks (name) VALUES ($1) RETURNING id", string(name),
	).Scan(&id)
	return id, err
}

func ScanBanks(rows *sql.Rows) ([]models.Bank, error) {
	var banks []models.Bank
	for rows.Next() {
		var b models.Bank
		if err := rows.Scan(&b.ID, &b.Name); err != nil {
			return nil, err
		}
		banks = append(banks, b)
	}
	return banks, nil
}

func GetBankByName(db *sql.DB, name models.BankName) (models.Bank, error) {
	var b models.Bank
	err := db.QueryRow("SELECT id, name FROM banks WHERE name = $1", name).Scan(&b.ID, &b.Name)
	if err != nil {
		return models.Bank{}, fmt.Errorf("bank %s not found: %w", name, err)
	}
	return b, nil
}
