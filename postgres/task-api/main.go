package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	UserID      int    `json:"user_id"`
}

func getTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(r.URL.Path, "/")

		if len(parts) > 2 && parts[2] != "" {
			taskID, err := strconv.Atoi(parts[len(parts)-1])

			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			var printTask Task
			row := db.QueryRow(
				"SELECT id, title, description, completed, user_id FROM tasks WHERE id = $1",
				taskID,
			)

			err = row.Scan(
				&printTask.ID,
				&printTask.Title,
				&printTask.Description,
				&printTask.Completed,
				&printTask.UserID,
			)

			if err == sql.ErrNoRows {
				http.Error(w, "Task not found", http.StatusNotFound)
				return
			}

			if err != nil {
				http.Error(w, "Database Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(printTask)
			return
		}

		rows, err := db.Query("SELECT id, title, description, completed, user_id FROM tasks")

		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		var tasks []Task

		for rows.Next() {
			var task Task

			err := rows.Scan(
				&task.ID,
				&task.Title,
				&task.Description,
				&task.Completed,
				&task.UserID,
			)

			if err != nil {
				http.Error(w, "Database Error", http.StatusInternalServerError)
				return
			}

			tasks = append(tasks, task)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	}
}

func createTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var task Task
		err := json.NewDecoder(r.Body).Decode(&task)

		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if task.ID <= 0 ||
			task.Title == "" ||
			task.UserID <= 0 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("INSERT INTO tasks (id, title, description, user_id) VALUES ($1, $2, $3, $4)", task.ID, task.Title, task.Description, task.UserID)

		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code == "23503" {
					http.Error(w, "Bad Request", http.StatusBadRequest)
					return
				}

				if pqErr.Code == "23505" {
					http.Error(w, "Task ID already exists", http.StatusConflict)
					return
				}
			}

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(task)
	}
}

func updateTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(r.URL.Path, "/")
		taskID, err := strconv.Atoi(parts[len(parts)-1])

		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		var task Task
		err = json.NewDecoder(r.Body).Decode(&task)

		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		result, err := db.Exec("UPDATE tasks SET title = $1, description = $2, completed = $3 WHERE id = $4", task.Title, task.Description, task.Completed, taskID)

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Task Not Found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		task.ID = taskID

		json.NewEncoder(w).Encode(task)
		return
	}
}

func deleteTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(r.URL.Path, "/")
		taskID, err := strconv.Atoi(parts[len(parts)-1])

		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		result, err := db.Exec("DELETE FROM tasks WHERE id = $1", taskID)

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Task Not Found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func tasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getTasks(db)(w, r)
		} else if r.Method == http.MethodPost {
			createTask(db)(w, r)
		} else if r.Method == http.MethodPut {
			updateTask(db)(w, r)
		} else if r.Method == http.MethodDelete {
			deleteTask(db)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func main() {
	fmt.Println("Task API starting...")

	db, err := sql.Open(
		"postgres",
		"host=localhost port=5432 user=postgres dbname=backend_db sslmode=disable",
	)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to PostgreSQL")

	http.HandleFunc("/tasks", tasksHandler(db))

	http.HandleFunc("/tasks/", tasksHandler(db))

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}
