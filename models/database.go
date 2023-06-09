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

	if db != nil {
		db.Close()
	}

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

		if globals.DBTYPE == "coeus-sample" {
			transaction(db, ddl_sample)
		} else {
			transaction(db, ddl_blank)
		}
	} else {
		// File/Database exists. Open it and return a connection
		db, _ = sql.Open("sqlite3", filename)
		if err != nil {
			log.Fatal(err.Error())
		}

		// The database exists, but it may not have the correct tables so DDL must be run again
		if globals.DBTYPE == "coeus-sample" {
			transaction(db, ddl_sample)
		} else {
			transaction(db, ddl_blank)
		}
	}
	return db
}

// ReseedSampleDB drops tables in the existing database and reseeds a new one with sample data
func ReseedSampleDB() {
	db := NewDB()

	// Pass the database connection to the transaction function
	transaction(db, ddl_sample)
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
