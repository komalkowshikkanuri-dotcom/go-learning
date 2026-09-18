package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Connect to PostgreSQL
	db, err := sql.Open(
		"postgres",
		"host=localhost port=5432 user=postgres dbname=backend_db sslmode=disable",
	)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Check database connection
	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSQL!")

	// --------------------------------------------------
	// 1. QueryRow() - Get one user
	// --------------------------------------------------

	id := 2

	var name string
	var age int

	err = db.QueryRow(
		"SELECT name, age FROM users WHERE id = $1",
		id,
	).Scan(&name, &age)

	if err == sql.ErrNoRows {
		log.Println("User Not Found")
		return
	} else if err != nil {
		log.Fatal(err)
	}

	log.Println("Name:", name)
	log.Println("Age:", age)

	// --------------------------------------------------
	// 2. Query() - Get multiple users
	// --------------------------------------------------

	age = 25

	rows, err := db.Query(
		"SELECT id, name, age FROM users WHERE age > $1 LIMIT 5",
		age,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var age int

		err := rows.Scan(&id, &name, &age)

		if err != nil {
			log.Fatal(err)
		}

		log.Println("ID:", id)
		log.Println("Name:", name)
		log.Println("Age:", age)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// --------------------------------------------------
	// 3. INSERT using Exec()
	// --------------------------------------------------

	result, err := db.Exec(
		"INSERT INTO users (id, name, age) VALUES ($1, $2, $3)",
		100007,
		"David",
		28,
	)

	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Rows affected:", rowsAffected)

	// --------------------------------------------------
	// 4. Verify INSERT
	// --------------------------------------------------

	err = db.QueryRow(
		"SELECT name, age FROM users WHERE id = $1",
		100007,
	).Scan(&name, &age)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserted user:", name)
	log.Println("Age:", age)

	// --------------------------------------------------
	// 5. UPDATE using Exec()
	// --------------------------------------------------

	result, err = db.Exec(
		"UPDATE users SET age = $1 WHERE id = $2",
		29,
		100006,
	)

	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err = result.RowsAffected()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Rows affected:", rowsAffected)

	// --------------------------------------------------
	// 6. Verify UPDATE
	// --------------------------------------------------

	err = db.QueryRow(
		"SELECT name, age FROM users WHERE id = $1",
		100006,
	).Scan(&name, &age)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Updated user:", name)
	log.Println("Age:", age)
	
	// --------------------------------------------------
	// 7. DELETE using Exec()
	// --------------------------------------------------

	result, err = db.Exec(
    		"DELETE FROM users WHERE id = $1",
    		100007,
	)

	if err != nil {
    		log.Fatal(err)
	}

	rowsAffected, err = result.RowsAffected()

	if err != nil {
    		log.Fatal(err)
	}

	log.Println("Rows deleted:", rowsAffected)

	// --------------------------------------------------
	// 8. Transactions
	// --------------------------------------------------


	_, err = db.Begin()

	if err != nil {
    		log.Fatal(err)
	}

	log.Println("Transaction started")
}