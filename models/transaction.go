package models

import (
	"fmt"
	"strconv"
	"time"
)

type ColumnType string

const (
	ColumnDate            ColumnType = "date"
	ColumnVendor          ColumnType = "vendor"
	ColumnDescription     ColumnType = "description"
	ColumnCategory        ColumnType = "category"
	ColumnAmount          ColumnType = "amount"
	ColumnPaymentMethodID ColumnType = "payment_method_id"
)

func (c ColumnType) IsValid() bool {
	switch c {
	case ColumnDate, ColumnVendor, ColumnAmount: // ColumnDescription, ColumnCategory, ColumnPaymentMethodID:
		return true
	}
	return false
}

// BankColumnMap maps CSV column names to standardized ColumnTypes
type BankColumnMap map[string]ColumnType

type Transaction struct {
	ID              int       `json:"id"`
	Date            time.Time `json:"date"`
	Vendor          string    `json:"vendor"`
	Description     string    `json:"description"`
	Category        string    `json:"category"`
	Amount          float64   `json:"amount"`
	PaymentMethodID int       `json:"payment_method_id"`
}

func NewTransaction(id int, dateStr, vendor, description, category, amount string, paymentMethodID int) (Transaction, error) {
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to parse date: %w", err)
	}

	parsedAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to parse amount: %w", err)
	}

	return Transaction{
		ID:              id,
		Date:            date,
		Vendor:          vendor,
		Description:     description,
		Category:        category,
		Amount:          parsedAmount,
		PaymentMethodID: paymentMethodID,
	}, nil
}
