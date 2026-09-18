package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open(
		"postgres",
		"host=localhost port=5432 user=postgres dbname=backend_db sslmode=disable",
	)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	stats := db.Stats()
	
	log.Println("Open Connection: ", stats.OpenConnections)
	log.Println("In use:", stats.InUse)
	log.Println("Idle:", stats.Idle)

	log.Println("Connected to PostgreSQL!")

	tx, err := db.Begin()

	if err != nil {
		log.Fatal(err)
	}

	defer tx.Rollback()

	log.Println("Transaction Started")

	// Operation 1: Update Alice
	result, err := tx.Exec(
		"UPDATE users SET age = $1 WHERE id = $2",
		30,
		1,
	)

	if err != nil {
		log.Fatal(err)
	}

	// Operation 2: Update Bob
	_, err = tx.Exec(
		"UPDATE users SET age = $1 WHERE id = $2",
		35,
		2,
	)

	if err != nil {
		log.Fatal(err)
	}

	// Operation 3: Intentionally invalid INSERT
	_, err = tx.Exec(
		"INSERT INTO jobs (id, title, user_id) VALUES ($1, $2, $3)",
		999999,
		"Invalid Job",
		999999,
	)

	if err != nil {
		log.Println("Operation failed:", err)
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Rows affected:", rowsAffected)

	err = tx.Commit()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Transaction Committed")
}