package queries

import (
	"database/sql"
	"tally/models"
)

func GetAllPaymentMethods(db *sql.DB) (*sql.Rows, error) {
	return db.Query("SELECT id, type, bank_id FROM payment_methods ORDER BY id")
}

func InsertPaymentMethod(db *sql.DB, methodType models.PaymentMethodType, bankID int) (int, error) {
	var id int
	err := db.QueryRow(
		"INSERT INTO payment_methods (type, bank_id) VALUES ($1, $2) RETURNING id",
		methodType, bankID,
	).Scan(&id)
	return id, err
}

func ScanPaymentMethods(rows *sql.Rows) ([]models.PaymentMethod, error) {
	var methods []models.PaymentMethod
	for rows.Next() {
		var m models.PaymentMethod
		if err := rows.Scan(&m.ID, &m.Type, &m.BankID); err != nil {
			return nil, err
		}
		methods = append(methods, m)
	}
	return methods, nil
}
