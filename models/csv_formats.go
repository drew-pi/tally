package models

// CSVFormatRow represents a single row from the csv_formats table
type CSVFormatRow struct {
	ID         int        `json:"id"`
	BankID     int        `json:"bank_id"`
	CSVColumn  string     `json:"csv_column"`
	ColumnType ColumnType `json:"column_type"`
}
