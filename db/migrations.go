package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

func RunMigrations(db *sql.DB) error {
	migrations, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		return fmt.Errorf("failed to find migrations: %w", err)
	}

	for _, migration := range migrations {
		sql, err := os.ReadFile(migration)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", migration, err)
		}

		_, err = db.Exec(string(sql))
		if err != nil {
			return fmt.Errorf("failed to run migration %s: %w", migration, err)
		}

		fmt.Println("ran migration:", migration)
	}

	return nil
}
