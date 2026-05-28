package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


func Connect() *sql.DB {
	err := godotenv.Load()
	if err != nil {
        log.Fatal("Error loading .env file")
    }

	connStr := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
    )

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to db:", err)
	}

	err = db.Ping()
	if err != nil {
        log.Fatal("Error pinging db:", err)
    }

	return db
}

func main() {
	fmt.Println("Hello World")
}