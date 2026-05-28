package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"tally/queries"
)

func GetAllTransactions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		rows, err := queries.GetAllTransactions(db)
		if err != nil {
			http.Error(w, "failed to get transactions", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		transactions, err := queries.ScanTransactions(rows)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(transactions)
		json.NewEncoder(w).Encode(map[string]any{
			"count":        len(transactions),
			"transactions": transactions,
		})
	}
}

// GET /api/transactions/{id}
func GetTransactionByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		rows, err := queries.GetTransactionById(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		transactions, err := queries.ScanTransactions(rows)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(transactions)
		json.NewEncoder(w).Encode(map[string]any{
			"count":        len(transactions),
			"transactions": transactions,
		})
	}
}

// GET /api/transactions/amount?amount_from=10.00&amount_to=50.00
func GetTransactionsByAmountRange(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		amountFromStr := r.URL.Query().Get("amount_from")
		amountToStr := r.URL.Query().Get("amount_to")

		var amountFrom, amountTo float64

		if amountFromStr != "" {
			v, err := strconv.ParseFloat(amountFromStr, 64)
			if err != nil {
				http.Error(w, "invalid amount_from", http.StatusBadRequest)
				return
			}
			amountFrom = float64(v)
		}

		if amountToStr != "" {
			v, err := strconv.ParseFloat(amountToStr, 64)
			if err != nil {
				http.Error(w, "invalid amount_to", http.StatusBadRequest)
				return
			}
			amountTo = float64(v)
		}

		rows, err := queries.GetTransactionsByAmountRange(db, amountFrom, amountTo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		transactions, err := queries.ScanTransactions(rows)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(transactions)
		json.NewEncoder(w).Encode(map[string]any{
			"count":        len(transactions),
			"transactions": transactions,
		})
	}
}

// func GetTransactionsByAmount(db *sql.DB, amount float64) (*sql.Rows, error) {
// 	return db.Query("SELECT id, date, description, amount, card  FROM transactions  WHERE amount = $1 ORDER BY date DESC;", amount)
// }

// func GetTransactionsByDateRange(db *sql.DB, dateFrom, dateTo time.Time) (*sql.Rows, error) {
// 	return db.Query("SELECT id, date, description, amount, card  FROM transactions  WHERE date BETWEEN $1 AND $2 ORDER BY date DESC;", dateFrom, dateTo)
// }

// func GetTransactionsByBank(db *sql.DB, card string) (*sql.Rows, error) {
// 	return db.Query("SELECT id, date, description, amount, card  FROM transactions  WHERE amount = $1;", card)
// }

// func GetTransactions(
// 	db *sql.DB,
// 	id sql.NullInt64,
// 	amount sql.NullFloat64,
// 	dateFrom sql.NullString,
// 	dateTo sql.NullString,
// 	card sql.NullString,
// ) (*sql.Rows, error) {
// 	return db.Query(`
// 	SELECT id, date, description, amount, card
// 		FROM transactions
// 		WHERE ($1::int IS NULL OR id = $1)
// 		AND ($2::numeric IS NULL OR amount = $2)
// 		AND ($3::date IS NULL OR date >= $3)
// 		AND ($4::date IS NULL OR date <= $4)
// 		AND ($5::card_type IS NULL OR card = $5)
// 		ORDER BY date DESC;
// 	`, id, amount, dateFrom, dateTo, card)
// }

// GetTransaction, UpdateTransaction
