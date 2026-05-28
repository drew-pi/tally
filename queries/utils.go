package queries

import (
	"database/sql"
	"tally/models"
)

func ScanTransactions(rows *sql.Rows) ([]models.Transaction, error) {
	var transactions []models.Transaction

	for rows.Next() {
		var dateStr, description, amount, card string
		var id int
		err := rows.Scan(&id, &dateStr, &description, &amount, &card)
		if err != nil {
			return nil, err
		}

		t, err := models.NewTransaction(id, dateStr, description, amount, card)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}
