package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"tally/models"
	"tally/queries"
	"time"

	"github.com/go-chi/chi/v5"
)

// GET /api/transactions/all
func GetAllTransactions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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
		json.NewEncoder(w).Encode(map[string]any{
			"count":        len(transactions),
			"transactions": transactions,
		})
	}
}

// GET /api/transactions/{id}
func GetTransactionByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
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
		json.NewEncoder(w).Encode(map[string]any{
			"count":        len(transactions),
			"transactions": transactions,
		})
	}
}

// GET /api/transactions/amount/range?amount_from=10.00&amount_to=50.00
func GetTransactionsByAmountRange(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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
		json.NewEncoder(w).Encode(map[string]any{
			"count":        len(transactions),
			"transactions": transactions,
		})
	}
}

// GET /api/transactions/amount?amount=10.00
func GetTransactionsByAmount(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		amountStr := r.URL.Query().Get("amount")

		var amount float64

		if amountStr != "" {
			v, err := strconv.ParseFloat(amountStr, 64)
			if err != nil {
				http.Error(w, "invalid amount_from", http.StatusBadRequest)
				return
			}
			amount = float64(v)
		}

		rows, err := queries.GetTransactionsByAmount(db, amount)
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
		json.NewEncoder(w).Encode(map[string]any{
			"count":        len(transactions),
			"transactions": transactions,
		})
	}
}

// GET /api/transactions/date/range?date_from=2026-01-01&date_to=2026-05-01
func GetTransactionsByDateRange(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		dateFromStr := r.URL.Query().Get("date_from")
		dateToStr := r.URL.Query().Get("date_to")

		var dateFrom, dateTo time.Time

		if dateFromStr != "" {
			t, err := time.Parse("2006-01-02", dateFromStr)
			if err != nil {
				http.Error(w, "invalid date_from, use YYYY-MM-DD", http.StatusBadRequest)
				return
			}
			dateFrom = time.Time(t)
		}

		if dateToStr != "" {
			t, err := time.Parse("2006-01-02", dateToStr)
			if err != nil {
				http.Error(w, "invalid date_to, use YYYY-MM-DD", http.StatusBadRequest)
				return
			}
			dateTo = time.Time(t)
		}

		rows, err := queries.GetTransactionsByDateRange(db, dateFrom, dateTo)
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
		json.NewEncoder(w).Encode(map[string]any{
			"count":        len(transactions),
			"transactions": transactions,
		})
	}
}

// GET /api/transactions/bank?bank=WF
func GetTransactionsByBank(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bankStr := r.URL.Query().Get("bank")

		if bankStr == "" {
			http.Error(w, "bank parameter required", http.StatusBadRequest)
			return
		}

		bank := models.BankName(bankStr)
		if !bank.IsKnown() {
			http.Error(w, "unknown bank", http.StatusBadRequest)
			return
		}

		rows, err := queries.GetTransactionsByBank(db, bank)
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
		json.NewEncoder(w).Encode(map[string]any{
			"count":        len(transactions),
			"transactions": transactions,
		})
	}
}

// GET /api/transactions?id=1&amount_from=10.00&amount_to=50.00&date_from=2026-01-01&date_to=2026-05-01&bank=WF
func GetTransactions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idStr := r.URL.Query().Get("id")
		amountFromStr := r.URL.Query().Get("amount_from")
		amountToStr := r.URL.Query().Get("amount_to")
		dateFromStr := r.URL.Query().Get("date_from")
		dateToStr := r.URL.Query().Get("date_to")
		bankStr := r.URL.Query().Get("bank")

		var id sql.NullInt64
		var amountFrom, amountTo sql.NullFloat64
		var dateFrom, dateTo sql.NullTime
		var bank sql.NullString

		if idStr != "" {
			v, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				http.Error(w, "invalid amount_from", http.StatusBadRequest)
				return
			}
			id = sql.NullInt64{Int64: v, Valid: true}
		}

		if amountFromStr != "" {
			v, err := strconv.ParseFloat(amountFromStr, 64)
			if err != nil {
				http.Error(w, "invalid amount_from", http.StatusBadRequest)
				return
			}
			amountFrom = sql.NullFloat64{Float64: v, Valid: true}
		}

		if amountToStr != "" {
			v, err := strconv.ParseFloat(amountToStr, 64)
			if err != nil {
				http.Error(w, "invalid amount_to", http.StatusBadRequest)
				return
			}
			amountTo = sql.NullFloat64{Float64: v, Valid: true}
		}

		if dateFromStr != "" {
			t, err := time.Parse("2006-01-02", dateFromStr)
			if err != nil {
				http.Error(w, "invalid date_from, use YYYY-MM-DD", http.StatusBadRequest)
				return
			}
			dateFrom = sql.NullTime{Time: t, Valid: true}
		}

		if dateToStr != "" {
			t, err := time.Parse("2006-01-02", dateToStr)
			if err != nil {
				http.Error(w, "invalid date_to, use YYYY-MM-DD", http.StatusBadRequest)
				return
			}
			dateTo = sql.NullTime{Time: t, Valid: true}
		}

		if bankStr != "" {
			bank = sql.NullString{String: bankStr, Valid: true}
		}

		rows, err := queries.GetTransactions(db, id, amountFrom, amountTo, dateFrom, dateTo, bank)
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
		json.NewEncoder(w).Encode(map[string]any{
			"count":        len(transactions),
			"transactions": transactions,
		})
	}
}

// UpdateTransaction
