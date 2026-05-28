package models

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
