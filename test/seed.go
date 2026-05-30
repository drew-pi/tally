package test

import (
	"database/sql"
	"fmt"
	"tally/models"
	"tally/queries"
)

func SeedCSVFormats(db *sql.DB) error {
	seedData := map[models.BankName]models.BankColumnMap{
		models.WellsFargo: {
			"date":        models.ColumnDate,
			"description": models.ColumnVendor,
			"amount":      models.ColumnAmount,
		},
		models.Fidelity: {
			"date":   models.ColumnDate,
			"name":   models.ColumnVendor,
			"amount": models.ColumnAmount,
		},
	}

	for bankName, columnMap := range seedData {
		bank, err := queries.GetBankByName(db, bankName)
		if err != nil {
			return fmt.Errorf("bank %s not found: %w", bankName, err)
		}

		for csvColumn, columnType := range columnMap {
			err := queries.InsertCSVFormat(db, bank.ID, csvColumn, columnType)
			if err != nil {
				return fmt.Errorf("failed to seed csv format: %w", err)
			}
		}
	}
	return nil
}

func SeedBanks(db *sql.DB) error {
	for _, name := range models.KnownBanks {
		_, err := db.Exec(`
            INSERT INTO banks (name) 
            VALUES ($1) 
            ON CONFLICT (name) DO NOTHING
        `, name)
		if err != nil {
			return fmt.Errorf("failed to seed bank %s: %w", name, err)
		}
	}
	return nil
}

func SeedPaymentMethods(db *sql.DB) error {
	seedData := []struct {
		methodType models.PaymentMethodType
		bankName   models.BankName
	}{
		{models.Debit, models.WellsFargo},
		{models.Credit, models.Fidelity},
	}

	for _, s := range seedData {
		bank, err := queries.GetBankByName(db, s.bankName)
		if err != nil {
			return fmt.Errorf("bank %s not found: %w", s.bankName, err)
		}

		_, err = db.Exec(`
            INSERT INTO payment_methods (type, bank_id)
            VALUES ($1, $2)
            ON CONFLICT DO NOTHING
        `, s.methodType, bank.ID)
		if err != nil {
			return fmt.Errorf("failed to seed payment method: %w", err)
		}
	}
	return nil
}
