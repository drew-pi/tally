package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"tally/models"
	"tally/queries"
)

func GetAllBanks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := queries.GetAllBanks(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		banks, err := queries.ScanBanks(rows)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"count": len(banks),
			"banks": banks,
		})
	}
}

func CreateBank(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		if body.Name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		id, err := queries.InsertBank(db, models.BankName(body.Name))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"id":   id,
			"name": body.Name,
		})
	}
}
