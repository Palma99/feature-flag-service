package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	dbConn, err := sql.Open("postgres", "postgresql://root:root@localhost:5432/local_feature_flag?sslmode=disable")
	if err != nil {
		panic(err)
	}

	db = dbConn
}

func getExecutedMigrations() []string {
	rows, err := db.Query("SELECT name FROM migrations where executed_at is not null")

	if err != nil {
		panic(err)
	}

	migrations := []string{}

	for rows.Next() {
		var migration string
		if err := rows.Scan(&migration); err != nil {
			panic(err)
		}
		migrations = append(migrations, migration)
	}

	return migrations
}

func createMigrationTable() error {

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			executed_at TIMESTAMP
		)
	`)

	return err
}

func getLocalMigrationsName() []string {
	fileNames := []string{}

	dir := "migrations"

	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	// Stampa i nomi dei file
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			fileNames = append(fileNames, strings.Split(file.Name(), ".sql")[0])
		}
	}

	return fileNames
}

func isMigrationExecuted(executedMigrations []string, name string) bool {
	for _, executedMigration := range executedMigrations {
		if executedMigration == name {
			return true
		}
	}

	return false
}

func executeMigration(name string) error {
	fileName := fmt.Sprintf("migrations/%s.sql", name)

	file, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(file))

	return err
}

func updateMigrationTable(migrationName string) error {
	_, err := db.Exec(`INSERT INTO migrations (name, executed_at) VALUES ($1, now())`, migrationName)
	return err
}

func main() {
	if err := createMigrationTable(); err != nil {
		panic(err)
	}

	localMigrationNames := getLocalMigrationsName()
	executedMigrations := getExecutedMigrations()

	if len(localMigrationNames) == len(executedMigrations) {
		fmt.Println("\033[32mNo migrations to execute.")
		return
	}

	for _, migrationName := range localMigrationNames {
		if !isMigrationExecuted(executedMigrations, migrationName) {
			fmt.Printf("\033[33mExecuting migration %s...\n", migrationName)
			if err := executeMigration(migrationName); err != nil {
				panic(err)
			}
			if err := updateMigrationTable(migrationName); err != nil {
				panic(err)
			}

			fmt.Printf("\033[32mMigration %s executed.\n", migrationName)
		}
	}

	fmt.Println("\033[32mAll migrations executed successfully.")
}
