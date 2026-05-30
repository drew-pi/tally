package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	csvparser "tally/csv"
	"tally/db"
	"tally/models"
	"tally/test"
)

func main() {

	// connect to db
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer database.Close()

	// run migrations to create tables
	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// seed banks from known constants
	if err := test.SeedBanks(database); err != nil {
		log.Fatalf("failed to seed banks: %v", err)
	}

	if err := test.SeedPaymentMethods(database); err != nil {
		log.Fatalf("failed to seed payment methods: %v", err)
	}

	// seed csv formats from BankConfigs
	if err := test.SeedCSVFormats(database); err != nil {
		log.Fatalf("failed to seed csv formats: %v", err)
	}

	// open and parse the csv
	fileName := string(models.WellsFargo) + "_04:2026.csv"
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))

	// import transactions from csv into db
	if err := csvparser.ImportTransactions(database, models.WellsFargo, reader); err != nil {
		log.Fatalf("failed to import transactions: %v", err)
	}

	// query and print all transactions to verify import
	rows, err := database.Query(`
		SELECT id, date, vendor, description, category, amount, payment_method_id 
		FROM transactions 
		ORDER BY date DESC
	`)
	if err != nil {
		log.Fatalf("failed to query transactions: %v", err)
	}
	defer rows.Close()

	fmt.Println("\n--- Imported Transactions ---")
	for rows.Next() {
		var id, paymentMethodID int
		var date, vendor, description, category, amount string
		err := rows.Scan(&id, &date, &vendor, &description, &category, &amount, &paymentMethodID)
		if err != nil {
			log.Fatalf("failed to scan row: %v", err)
		}
		fmt.Printf("id=%-4d date=%-25s vendor=%-30s amount=%s\n", id, date, vendor, amount)
	}
	fmt.Println("-----------------------------")

	// start http server
	router := setupRoutes(database)
	log.Println("server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
