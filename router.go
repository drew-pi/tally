package main

import (
	"database/sql"
	"net/http"
	"tally/handlers"
	"tally/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func setupRoutes(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	// built-in chi middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Recoverer) // catches panics

	// your custom middleware
	r.Use(middleware.Logger)

	// api routes grouped under /api
	r.Route("/api", func(r chi.Router) {

		// time
		r.Get("/time", handlers.GetTime)

		// transactions
		r.Route("/transactions", func(r chi.Router) {
			r.Get("/", handlers.GetTransactions(db))
			r.Get("/all", handlers.GetAllTransactions(db))
			r.Get("/{id}", handlers.GetTransactionByID(db))
			r.Get("/amount", handlers.GetTransactionsByAmount(db))
			r.Get("/amount/range", handlers.GetTransactionsByAmountRange(db))
			r.Get("/date/range", handlers.GetTransactionsByDateRange(db))
			r.Get("/bank", handlers.GetTransactionsByBank(db))
		})

		// banks
		r.Route("/banks", func(r chi.Router) {
			r.Get("/", handlers.GetAllBanks(db))
			r.Post("/", handlers.CreateBank(db))
		})

		// payment methods
		r.Route("/payment-methods", func(r chi.Router) {
			r.Get("/", handlers.GetAllPaymentMethods(db))
			r.Post("/", handlers.CreatePaymentMethod(db))
		})

		// csv formats
		r.Route("/csv-formats", func(r chi.Router) {
			r.Get("/", handlers.GetAllCSVFormats(db))
			r.Get("/bank", handlers.GetCSVFormat(db))
			r.Post("/", handlers.CreateCSVFormat(db))
		})

	})

	// static frontend
	r.Handle("/*", http.FileServer(http.Dir("frontend/")))

	return r
}
