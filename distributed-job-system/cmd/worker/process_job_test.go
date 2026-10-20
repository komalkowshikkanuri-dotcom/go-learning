package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	//"strconv"
	"sync"
	"testing"

	//goredis "github.com/redis/go-redis/v9"

	"github.com/komalkowshikkanuri/distributed-job-system/internal/redis"
	"github.com/komalkowshikkanuri/distributed-job-system/internal/repository"
	"github.com/komalkowshikkanuri/distributed-job-system/internal/service"

	_ "github.com/lib/pq"
)

func TestProcessJob(t *testing.T) {
	redis.NewClient()
	defer redis.Client.Close()

	if err := redis.Ping(); err != nil {
		fmt.Println("Redis Connection failed", err)
		t.Fatal(err)
	}

	fmt.Println("Worker connected to Redis")

	cxt := context.Background()

	connStr := "host=127.0.0.1 port=5432 database=job_system user=postgres password=" + os.Getenv("DB_PASSWORD") + " sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	repo := repository.NewJobRepository(db)
	jobService := service.NewJobService(repo)

	var jobID int

	err = db.QueryRow(
		`INSERT INTO jobs (type, status)
		VALUES ($1, $2)
		RETURNING id`,
		"week4_day1_test", "pending",
	).Scan(&jobID)

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_, err := db.Exec(
			`DELETE FROM job_attempts WHERE job_id = $1`,
			jobID,
		)

		if err != nil {
			t.Errorf("failed to delete test attempts: %v", err)
		}

		_, err = db.Exec(
			`DELETE FROM jobs WHERE id = $1`,
			jobID,
		)

		if err != nil {
			t.Errorf("failed to delete test job: %v", err)
		}
	})

	t.Logf("Created job with ID: %d", jobID)

	shouldRetry, err := processJob(jobID, jobService, cxt, redis.Client)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(shouldRetry)

	var status string

	err = db.QueryRow(
		`SELECT status FROM jobs WHERE id = $1`,
		jobID,
	).Scan(&status)

	if err != nil {
		t.Fatal(err)
	}

	if status != "pending" {
		t.Fatalf("expected job status pending, got %q", status)
	}

	var attemptCount int

	err = db.QueryRow(
		`SELECT COUNT(*) FROM job_attempts WHERE job_id = $1`,
		jobID,
	).Scan(&attemptCount)

	if err != nil {
		t.Fatal(err)
	}

	if attemptCount != 1 {
		t.Fatalf("expected 1 attempt, got %d", attemptCount)
	}

	var attemptStatus string

	err = db.QueryRow(
		`SELECT status FROM job_attempts WHERE job_id = $1`,
		jobID,
	).Scan(&attemptStatus)

	if err != nil {
		t.Fatal(err)
	}

	if attemptStatus != "failed" {
		t.Fatalf("expect 'failed', got %q", attemptStatus)
	}

	if !shouldRetry {
		t.Fatal("expected job to be eligible for retry")
	}

	//queuedJobID, err := redis.Client.RPop(cxt, "job_queue").Result()
	//if err != nil {
	//	t.Fatalf("expected job to be requeued in Redis: %v", err)
	//}

	//if queuedJobID != strconv.Itoa(jobID) {
	//	t.Fatalf("expected queued job ID %q, got %q", strconv.Itoa(jobID), queuedJobID)
	//}

}

func TestUpdateAttemptStatus(t *testing.T) {
	connStr := "host=127.0.0.1 port=5432 database=job_system user=postgres password=" + os.Getenv("DB_PASSWORD") + " sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	repo := repository.NewJobRepository(db)
	jobService := service.NewJobService(repo)

	var jobID int

	err = db.QueryRow(
		`INSERT INTO jobs (type, status)
		VALUES ($1, $2)
		RETURNING id`,
		"week4_day2_test", "pending",
	).Scan(&jobID)

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_, err := db.Exec(
			`DELETE FROM job_attempts WHERE job_id = $1`,
			jobID,
		)

		if err != nil {
			t.Errorf("failed to delete test attempts: %v", err)
		}

		_, err = db.Exec(
			`DELETE FROM jobs WHERE id = $1`,
			jobID,
		)

		if err != nil {
			t.Errorf("failed to delete test job: %v", err)
		}
	})

	attemptID, err := jobService.CreateAttempt(jobID)

	if err != nil {
		t.Fatal(err)
	}

	err = jobService.UpdateAttemptStatus(attemptID, "failed")

	if err != nil {
		t.Fatal(err)
	}

	var finished_at sql.NullTime
	var status string

	err = db.QueryRow(
		`SELECT status, finished_at
		FROM job_attempts
		WHERE id = $1`,
		attemptID,
	).Scan(&status, &finished_at)

	if err != nil {
		t.Fatal(err)
	}

	if status != "failed" || finished_at.Valid != true {
		t.Fatalf("Expected 'failed', got %q. Expected time, got %v", status, finished_at)
	}

	processingAttemptID, err := jobService.CreateAttempt(jobID)
	if err != nil {
		t.Fatal(err)
	}

	err = jobService.UpdateAttemptStatus(-1, "failed")
	if err == nil {
		t.Fatal(err)
	}

	err = db.QueryRow(
		`SELECT status, finished_at
		FROM job_attempts
		WHERE id = $1`,
		processingAttemptID,
	).Scan(&status, &finished_at)

	if err != nil {
		t.Fatal(err)
	}

	if status != "processing" || finished_at.Valid {
		t.Fatalf("Expected 'processing', got %q. Expected time, got %v", status, finished_at)
	}
}

func TestMarkJobActive(t *testing.T) {
	activeJobs := &sync.Map{}

	job1 := markJobActive(activeJobs, 101)

	if !job1 {
		t.Fatal("expected job 101 to be marked active")
	}

	job1 = markJobActive(activeJobs, 101)

	if job1 {
		t.Fatal("expected job 101 to be marked active")
	}

	markJobInactive(activeJobs, 101)

	job1 = markJobActive(activeJobs, 101)

	if !job1 {
		t.Fatal("expected job 101 to be marked active")
	}

	job2 := markJobActive(activeJobs, 202)

	if !job2 {
		t.Fatal("expected job 202 to be marked active")
	}
}

func TestProcessMultipleJobs(t *testing.T) {
	connStr := "host=127.0.0.1 port=5432 database=job_system user=postgres password=" + os.Getenv("DB_PASSWORD") + " sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	redis.NewClient()
	defer redis.Client.Close()

	if err := redis.Ping(); err != nil {
		t.Fatal(err)
	}

	fmt.Println("Worker connected to Redis")

	cxt := context.Background()

	repo := repository.NewJobRepository(db)
	jobService := service.NewJobService(repo)

	var jobID int
	var jobIDs []int

	t.Cleanup(func() {
		for _, id := range jobIDs {
			_, err := db.Exec(`DELETE FROM job_attempts WHERE job_id = $1`,
				id,
			)

			if err != nil {
				t.Errorf("failed to delete test attempts: %v", err)
			}

			_, err = db.Exec(`DELETE FROM jobs WHERE id = $1`,
				id,
			)

			if err != nil {
				t.Errorf("failed to delete test job: %v", err)
			}
		}
	})

	for range 3 {
		err = db.QueryRow(
			`INSERT INTO jobs (type, status)
        		VALUES ($1, $2)
        		RETURNING id`,
			"week4_day4_test", "pending",
		).Scan(&jobID)

		if err != nil {
			t.Fatal(err)
		}

		jobIDs = append(jobIDs, jobID)
	}

	shouldRetryByJob := make(map[int]bool)

	for _, id := range jobIDs {
		shouldRetry, err := processJob(id, jobService, cxt, redis.Client)

		if err != nil {
			t.Fatal(err)
		}

		shouldRetryByJob[id] = shouldRetry
	}

	for _, id := range jobIDs {
		var status string

		err := db.QueryRow(
			`SELECT status FROM jobs 
			WHERE id = $1`,
			id,
		).Scan(&status)

		if err != nil {
			t.Fatalf("Failed to load the status %v", err)
		}

		if status != "pending" {
			t.Errorf("Expected pending, got: %q", status)
		}
	}

	for _, id := range jobIDs {
		var attemptCount int

		err := db.QueryRow(
			`SELECT COUNT(*) 
			FROM job_attempts 
			WHERE job_id = $1`,
			id,
		).Scan(&attemptCount)

		if err != nil {
			t.Fatalf("Failed to count the job attempts %v", err)
		}

		if attemptCount != 1 {
			t.Errorf("expected 1 attempt, got %d", attemptCount)
		}

	}

	for _, id := range jobIDs {
		var attemptStatus string

		err := db.QueryRow(
			`SELECT status FROM job_attempts 
			WHERE job_id = $1`,
			id,
		).Scan(&attemptStatus)

		if err != nil {
			t.Fatalf("Failed to load the attempt status %v", err)
		}

		if attemptStatus != "failed" {
			t.Errorf("expected 'failed' attempt, got %q", attemptStatus)
		}

		if !shouldRetryByJob[id] {
			t.Errorf("job %d was expected to be eligible for retry", id)
		}
	}

}
