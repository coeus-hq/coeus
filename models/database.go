package models

import (
	globals "coeus/globals"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// NewDB opens the configured database and returns a connection to it.
// If the database does not exist, it creates a new one with sample data
func NewDB() *sql.DB {
	var db *sql.DB

	// TODO load config from .env using embedFS if running as a single binary

	// Check to see the file exists
	var filename string = globals.DBNAME + ".db"
	_, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Creating database…")
		file, _ := os.Create(filename)
		fmt.Printf("Database %s created\n", filename)
		file.Close()

		db, err = sql.Open("sqlite3", filename)
		if err != nil {
			log.Fatal(err.Error())
		}
		transaction(db, ddl)
		return db
	}

	// File/Database exists. Open it and return a connection
	db, _ = sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}

// transaction executes an array of SQL statememts in a single transaction.
func transaction(db *sql.DB, statements []string) {
	var ctx = context.Background()
	var tx, err = db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("Couldn't begin transaction")
		panic(err)
	}
	for _, statement := range statements {
		_, err = tx.ExecContext(ctx, statement)
		if err != nil {
			tx.Rollback()
			if err != nil {
				fmt.Println(statement)
				panic(err)
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println("Failed to commit transaction")
		panic(err)
	}
}
