package main

import (
	"database/sql"
	"net/http"
	"tally/handlers"
	"tally/middleware"
)

func setupRoutes(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// API routes
	// mux.HandleFunc("/api/time", handlers.GetTime)
	// mux.HandleFunc("/api/transactions", handlers.GetAllTransactions(db))

	mux.HandleFunc("/api/time", handlers.GetTime)
	mux.HandleFunc("/api/transactions", handlers.GetAllTransactions(db))
	mux.HandleFunc("/api/transactions/{id}", handlers.GetTransactionByID(db))
	mux.HandleFunc("/api/transactions/amount", handlers.GetTransactionsByAmountRange(db))

	// Serve static frontend files
	mux.Handle("/", http.FileServer(http.Dir("frontend/")))

	return middleware.Logger(mux)
}
