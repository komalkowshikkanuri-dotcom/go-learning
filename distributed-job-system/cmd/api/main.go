package main

import (
	"fmt"
	"net/http"
	"log"
	"os"
	"encoding/json"

	"database/sql"
	_ "github.com/lib/pq"
)

type Job struct {
	JobID   int     `json:"job_id"`
	JobType string  `json:"type"`
	Status  string  `json:"status"`
}

type CreateJobRequest struct {
	UserID int    `json:"user_id"`
	Type   string `json:"type"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ROOT HANDLER:", r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, "Distributed Job System API")
}


func createJob(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CREATE JOB HANDLER:", r.Method, r.URL.Path)

		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var req CreateJobRequest
		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.UserID == 0 || req.Type == "" {
			http.Error(w, "Invalid job data", http.StatusBadRequest)
			return
		}

		var job Job
		err = db.QueryRow(
			`INSERT INTO jobs (user_id, type, status)
			VALUES ($1, $2, $3)
			RETURNING id, type, status`,
			req.UserID, req.Type, "pending",
		).Scan(&job.JobID, &job.JobType, &job.Status)

		if err != nil {
			fmt.Println("DATABASE ERROR:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(job)
	}
}

func getJobs(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		rows, err := db.Query(
			`SELECT id, type, status FROM jobs`,
		)
		
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		defer rows.Close()
		var jobs []Job
		
		for rows.Next() {
			var job Job	
			err := rows.Scan (
				&job.JobID,
				&job.JobType,
				&job.Status,
			)
			
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			jobs = append(jobs, job)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jobs)
	}
}


func main() {
	db, err := sql.Open(
		"postgres",
    		"host=127.0.0.1 port=5432 database=job_system user=postgres password=" + os.Getenv("DB_PASSWORD") + " sslmode=disable",
	)

	if err != nil {
		log.Fatal(err)
		return
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Database connected!")

	http.HandleFunc("/", handler)
	http.HandleFunc("/job", createJob(db))
	http.HandleFunc("/jobs", getJobs(db))
	
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}