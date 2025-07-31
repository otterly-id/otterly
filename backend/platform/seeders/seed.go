package main

import (
	"encoding/json"
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
	if err := godotenv.Load(); err != nil {
		fmt.Println("no .env file found")
		return
	}

	db, err := database.OpenDBConnection()
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

		execQuery(data, db, file.Name())
	}
}

func execQuery(data Seed, db *database.Queries, fileName string) {
	if db == nil {
		fmt.Printf("database connection is nil, cannot execute query for file: %s\n", fileName)
		return
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

	for _, value := range data.Values {
		if data.Table == "users" && passwordIndex != -1 {
			password, ok := value[passwordIndex].(string)
			if !ok {
				fmt.Printf("password field is not a string, skipping row in file: %s\n", fileName)
				continue
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				fmt.Printf("failed to hash password, skipping row in file %s: %v\n", fileName, err)
				continue
			}

			value[passwordIndex] = hashedPassword
		}

		if _, err := db.Exec(sqlx.Rebind(sqlx.DOLLAR, query), value...); err != nil {
			fmt.Printf("error in running seeder file %s: %v\n", fileName, err)
		}
	}
}

func prepareInsertQuery(columns []string) string {
	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	return fmt.Sprintf("(%s)", strings.Join(placeholders, ","))
}
