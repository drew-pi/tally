package models

import (
	"fmt"
	"strconv"
	"time"
)

type Transaction struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Card        string    `json:"card"`
}

func NewTransaction(id int, dateStr, description, amount, card string) (Transaction, error) {
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to parse date: %w", err)
	}

	parsedAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to parse amount: %w", err)
	}

	return Transaction{
		ID:          id,
		Date:        date,
		Description: description,
		Amount:      parsedAmount,
		Card:        card,
	}, nil
}

type ColumnType string

const (
	ColumnDate   ColumnType = "date"
	ColumnName   ColumnType = "name"
	ColumnAmount ColumnType = "amount"
)

// BankColumnMap maps CSV column names to standardized ColumnTypes
type BankColumnMap map[string]ColumnType

// CSVConfig maps bank names to their column mappings
type CSVConfig map[string]BankColumnMap

var BankConfigs = CSVConfig{
	"WF": {
		"date":        ColumnDate,
		"description": ColumnName,
		"amount":      ColumnAmount,
	},
	"Fidelity": {
		"date":   ColumnDate,
		"name":   ColumnName,
		"amount": ColumnAmount,
	},
}
