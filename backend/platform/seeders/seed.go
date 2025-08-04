package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/otterly-id/otterly/backend/platform/database"
	"golang.org/x/crypto/bcrypt"
)

const (
	seedDir = "./platform/seeders/data"
)

type Seed struct {
	Table   string   `json:"table"`
	Columns []string `json:"columns"`
	Values  [][]any  `json:"values"`
}

func main() {
	stopOnError := flag.Bool("stop-on-error", false, "Stop the seeder immediately if an error occurs")
	flag.Parse()

	fmt.Println("Seeding database started")

	if err := godotenv.Load(); err != nil {
		fmt.Println("no .env file found")
		return
	}

	db, err := database.GetDBConnection()
	if err != nil {
		fmt.Printf("failed to open database connection: %v", err)
		return
	}

	files, err := os.ReadDir(seedDir)
	if err != nil {
		fmt.Printf("failed to read seed directory: %v\n", err)
		return
	}

	for _, file := range files {
		f := strings.Split(file.Name(), ".")

		if file.IsDir() || f[len(f)-1] != "json" {
			continue
		}

		content, err := os.ReadFile(filepath.Join(seedDir, file.Name()))
		if err != nil {
			fmt.Printf("error during reading file %s: %v\n", file.Name(), err)
			continue
		}

		var data Seed
		if err := json.Unmarshal(content, &data); err != nil {
			fmt.Printf("error during unmarshalling file content from %s: %v\n", file.Name(), err)
			continue
		}

		insertedCount, err := execQuery(data, db.DB, file.Name(), *stopOnError)
		if insertedCount > 0 {
			fmt.Printf("File '%s': successfully inserted %d rows.\n", file.Name(), insertedCount)
		}

		if err != nil {
			fmt.Printf("\nStopping seeder due to error: %v", err)
			return
		}
	}

	fmt.Println("Database seeded successfully")
}

func execQuery(data Seed, db *sqlx.DB, fileName string, stopOnError bool) (int, error) {
	if db == nil {
		return 0, fmt.Errorf("database connection is nil for file: %s", fileName)
	}

	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES %s`,
		data.Table,
		strings.Join(data.Columns, ","),
		prepareInsertQuery(data.Columns),
	)

	passwordIndex := -1
	if data.Table == "users" {
		for i, col := range data.Columns {
			if col == "password" {
				passwordIndex = i
				break
			}
		}
	}

	successfulInserts := 0
	for _, value := range data.Values {
		if data.Table == "users" && passwordIndex != -1 {
			password, ok := value[passwordIndex].(string)
			if !ok {
				err := fmt.Errorf("password field is not a string in file: %s", fileName)
				if stopOnError {
					return successfulInserts, err
				}
				fmt.Printf("%v, skipping row.\n", err)
				continue
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				err := fmt.Errorf("failed to hash password in file %s: %v", fileName, err)
				if stopOnError {
					return successfulInserts, err
				}
				fmt.Printf("%v, skipping row.\n", err)
				continue
			}
			value[passwordIndex] = string(hashedPassword)
		}

		if _, err := db.Exec(sqlx.Rebind(sqlx.DOLLAR, query), value...); err != nil {
			err := fmt.Errorf("error inserting row from file %s: %v", fileName, err)
			if stopOnError {
				return successfulInserts, err
			}
			fmt.Println(err)
		} else {
			successfulInserts++
		}
	}

	return successfulInserts, nil
}

func prepareInsertQuery(columns []string) string {
	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	return fmt.Sprintf("(%s)", strings.Join(placeholders, ","))
}
