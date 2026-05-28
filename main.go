package main

import (
	// "bufio"
	// "encoding/csv"
	"bufio"
	"encoding/csv"
	"fmt"
	"os"

	// "log"
	// "os"
	csvparser "tally/csv"
	// "tally/db"
	// "tally/queries"

	"log"
	"net/http"
	"tally/db"
)

func main() {

	fmt.Println("Hello, World!")

	card := "WF"
	file_name := card + "_04:2026.csv"

	file, err := os.Open(file_name)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	bufferedFile := bufio.NewReader(file)
	reader := csv.NewReader(bufferedFile)

	database, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer database.Close()

	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	csvparser.ImportTransactions(database, card, reader)

	rows, err := database.Query("SELECT date, description, amount, card FROM transactions")
	if err != nil {
		fmt.Println("Error querying db:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var date, description, amount, card string
		err := rows.Scan(&date, &description, &amount, &card)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		fmt.Println(date, description, amount, card)
	}

	router := setupRoutes(database)

	log.Println("server starting on :8080")
	http.ListenAndServe(":8080", router)
}

// func main() {
// 	bank := "WF"
// 	fileName := bank + "_04:2026.csv"

// 	file, err := os.Open(fileName)
// 	if err != nil {
// 		log.Fatalf("failed to open file: %v", err)
// 	}
// 	defer file.Close()

// 	reader := csv.NewReader(bufio.NewReader(file))

// 	database, err := db.Connect()
// 	if err != nil {
// 		log.Fatalf("failed to connect to db: %v", err)
// 	}
// 	defer database.Close()

// 	if err := db.RunMigrations(database); err != nil {
// 		log.Fatalf("failed to run migrations: %v", err)
// 	}

// 	if err := csvparser.ImportTransactions(database, bank, reader); err != nil {
// 		log.Fatalf("failed to import transactions: %v", err)
// 	}

// 	rows, err := queries.GetAllTransactions(database)
// 	if err != nil {
// 		log.Fatalf("error querying db: %v", err)
// 		return
// 	}

// 	for rows.Next() {
// 		var date, description, amount, card string
// 		err := rows.Scan(&date, &description, &amount, &card)
// 		if err != nil {
// 			fmt.Println("Error scanning row:", err)
// 			return
// 		}
// 		fmt.Println(date, description, amount, card)
// 	}

// 	fmt.Println("Import successful")

// }
