package initDB

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func CheckDB(dbName string) error {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	if err = createTables(db); err != nil {
		return err
	}
	return nil
}

func createTables(db *sql.DB) error {
	tableUsers := `CREATE TABLE IF NOT EXISTS storage_passwords (
		     login TEXT,
		     service_login TEXT,
		     service_name TEXT,
		     password TEXT unique,
		     create_date  TEXT
	)`

	tableStorage := `CREATE TABLE IF NOT EXISTS auth (
		     login TEXT,
		     hash TEXT
	)`

	if _, err := db.Exec(tableUsers); err != nil {
		return err
	}

	if _, err := db.Exec(tableStorage); err != nil {
		return err
	}
	return nil
}
