package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"tally/models"
	"tally/queries"
)

func GetAllPaymentMethods(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := queries.GetAllPaymentMethods(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		methods, err := queries.ScanPaymentMethods(rows)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"count":           len(methods),
			"payment_methods": methods,
		})
	}
}

func CreatePaymentMethod(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Type   string `json:"type"`
			BankID int    `json:"bank_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		methodType := models.PaymentMethodType(body.Type)
		if !methodType.IsValid() {
			http.Error(w, "invalid payment method type", http.StatusBadRequest)
			return
		}

		id, err := queries.InsertPaymentMethod(db, methodType, body.BankID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"id":      id,
			"type":    body.Type,
			"bank_id": body.BankID,
		})
	}
}
