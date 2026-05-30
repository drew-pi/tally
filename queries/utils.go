package queries

import (
	"database/sql"
	"tally/models"
)

func ScanTransactions(rows *sql.Rows) ([]models.Transaction, error) {
	var transactions []models.Transaction

	for rows.Next() {
		var id, paymentMethodID int
		var dateStr, vendor, description, category, amount string

		err := rows.Scan(&id, &dateStr, &vendor, &description, &category, &amount, &paymentMethodID)
		if err != nil {
			return nil, err
		}

		t, err := models.NewTransaction(id, dateStr, vendor, description, category, amount, paymentMethodID)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}
