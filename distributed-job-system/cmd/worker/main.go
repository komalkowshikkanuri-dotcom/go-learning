package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/komalkowshikkanuri/distributed-job-system/internal/redis"
	"github.com/komalkowshikkanuri/distributed-job-system/internal/repository"
	"github.com/komalkowshikkanuri/distributed-job-system/internal/service"

	_ "github.com/lib/pq"
)

func main() {
	redis.NewClient()

	if err := redis.Ping(); err != nil {
		fmt.Println("Redis Connection failed", err)
		return
	}

	fmt.Println("Worker connected to Redis")

	cxt := context.Background()

	jobID, err := redis.Client.RPop(cxt, "job_queue").Result()

	if err != nil {
		fmt.Println("Failed to get a job", err)
		return
	}

	jobIDINT, err := strconv.Atoi(jobID)

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(jobIDINT)

	connStr := "host=127.0.0.1 port=5432 database=job_system user=postgres password=" + os.Getenv("DB_PASSWORD") + " sslmode=disable"

	fmt.Println("DB password set:", os.Getenv("DB_PASSWORD") != "")

	fmt.Println("Host: 127.0.0.1")
	fmt.Println("Port: 5432")
	fmt.Println("Database: job_system")
	fmt.Println("User: postgres")
	fmt.Println("Password set:", os.Getenv("DB_PASSWORD") != "")
	fmt.Println("SSL mode: disable")

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
		return
	}

	if err := db.Ping(); err != nil {
		log.Fatal("PostgreSQL connection failed:", err)
	}

	fmt.Println("Worker Connected to PostgreSQL")

	repo := repository.NewJobRepository(db)

	jobService := service.NewJobService(repo)

	job, err := jobService.GetJob(jobIDINT)

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Job Found:", job)

	err = jobService.UpdateJob(job.JobID, "processing")

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Job is now Processing")
}
