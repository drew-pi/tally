package queries

import (
	"database/sql"
	"tally/models"
)

func GetAllCSVFormats(db *sql.DB) ([]models.CSVFormatRow, error) {
	rows, err := db.Query(`
		SELECT id, bank_id, csv_column, column_type 
		FROM csv_formats 
		ORDER BY bank_id, csv_column
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var formats []models.CSVFormatRow
	for rows.Next() {
		var f models.CSVFormatRow
		if err := rows.Scan(&f.ID, &f.BankID, &f.CSVColumn, &f.ColumnType); err != nil {
			return nil, err
		}
		formats = append(formats, f)
	}
	return formats, nil
}

func GetCSVFormat(db *sql.DB, bankID int) (models.BankColumnMap, error) {
	rows, err := db.Query(`
        SELECT csv_column, column_type 
        FROM csv_formats 
        WHERE bank_id = $1
    `, bankID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columnMap := make(models.BankColumnMap)
	for rows.Next() {
		var csvColumn, columnType string
		if err := rows.Scan(&csvColumn, &columnType); err != nil {
			return nil, err
		}
		columnMap[csvColumn] = models.ColumnType(columnType)
	}
	return columnMap, nil
}

func InsertCSVFormat(db *sql.DB, bankID int, csvColumn string, columnType models.ColumnType) error {
	_, err := db.Exec(`
        INSERT INTO csv_formats (bank_id, csv_column, column_type)
        VALUES ($1, $2, $3)
        ON CONFLICT DO NOTHING
    `, bankID, csvColumn, string(columnType))
	return err
}
