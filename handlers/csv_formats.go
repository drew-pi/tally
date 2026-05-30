package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"tally/models"
	"tally/queries"
)

// GET /api/csv-formats — returns all csv formats across all banks
func GetAllCSVFormats(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		formats, err := queries.GetAllCSVFormats(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"count":   len(formats),
			"formats": formats,
		})
	}
}

// GET /api/csv-formats?bank_id=1 — returns csv formats for a specific bank
func GetCSVFormat(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bankIDStr := r.URL.Query().Get("bank_id")
		if bankIDStr == "" {
			http.Error(w, "bank_id is required", http.StatusBadRequest)
			return
		}

		bankID, err := strconv.Atoi(bankIDStr)
		if err != nil {
			http.Error(w, "invalid bank_id", http.StatusBadRequest)
			return
		}

		format, err := queries.GetCSVFormat(db, bankID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"bank_id": bankID,
			"format":  format,
		})
	}
}

// POST /api/csv-formats — creates a new csv format mapping
// body: {"bank_id": 1, "csv_column": "name", "column_type": "vendor"}
func CreateCSVFormat(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			BankID     int    `json:"bank_id"`
			CSVColumn  string `json:"csv_column"`
			ColumnType string `json:"column_type"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		if body.CSVColumn == "" {
			http.Error(w, "csv_column is required", http.StatusBadRequest)
			return
		}

		colType := models.ColumnType(body.ColumnType)
		if !colType.IsValid() {
			http.Error(w, "invalid column_type", http.StatusBadRequest)
			return
		}

		err := queries.InsertCSVFormat(db, body.BankID, body.CSVColumn, colType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"bank_id":     body.BankID,
			"csv_column":  body.CSVColumn,
			"column_type": body.ColumnType,
		})
	}
}
