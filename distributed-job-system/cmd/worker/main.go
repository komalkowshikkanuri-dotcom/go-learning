package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/komalkowshikkanuri/distributed-job-system/internal/config"
	"github.com/komalkowshikkanuri/distributed-job-system/internal/redis"
	"github.com/komalkowshikkanuri/distributed-job-system/internal/repository"
	"github.com/komalkowshikkanuri/distributed-job-system/internal/service"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()
	redis.NewClient(cfg.RedisAddr)
	defer redis.Client.Close()

	if err := redis.Ping(); err != nil {
		fmt.Println("Redis Connection failed", err)
		return
	}

	fmt.Println("Worker connected to Redis")

	cxt := context.Background()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)
	fmt.Println("DB password set:", os.Getenv("DB_PASSWORD") != "")

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("PostgreSQL connection failed:", err)
	}

	fmt.Println("Worker Connected to PostgreSQL")

	repo := repository.NewJobRepository(db)
	jobService := service.NewJobService(repo)

	var wg sync.WaitGroup
	errorChannel := make(chan error)
	receiverDone := make(chan struct{})
	activeJobs := &sync.Map{}

	go func() {
		defer close(receiverDone)

		for err := range errorChannel {
			fmt.Println("Worker error:", err)
		}
	}()

	semaphore := make(chan struct{}, 3)

	for {
		jobID, err := redis.Client.RPop(cxt, "job_queue").Result()

		if err != nil {
			if errors.Is(err, goredis.Nil) {
				time.Sleep(time.Second)
				continue
			}

			fmt.Println("Failed to fetch job from Redis:", err)
			break
		}

		fmt.Println("Got job:", jobID)

		jobIDINT, err := strconv.Atoi(jobID)

		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Println(jobIDINT)

		active := markJobActive(activeJobs, jobIDINT)
		if !active {
			continue
		}

		semaphore <- struct{}{}

		wg.Add(1)

		go func(jobIDINT int) {
			defer wg.Done()

			defer func() {
				<-semaphore
			}()

			shouldRetry, err := processJob(jobIDINT, jobService, cxt, redis.Client)

			if err != nil {
				errorChannel <- err
				markJobInactive(activeJobs, jobIDINT)
				return
			}

			// Release the active marker before requeueing.
			markJobInactive(activeJobs, jobIDINT)

			if shouldRetry {
				err := redis.Client.LPush(
					cxt,
					"job_queue",
					strconv.Itoa(jobIDINT),
				).Err()

				if err != nil {
					errorChannel <- fmt.Errorf("failed to requeue job %d: %w", jobIDINT, err)
				}
			}
		}(jobIDINT)
	}

	wg.Wait()
	close(errorChannel)
	<-receiverDone

	fmt.Println("All workers finished")
}

func markJobActive(activeJobs *sync.Map, jobID int) bool {
	_, loaded := activeJobs.LoadOrStore(jobID, struct{}{})
	return !loaded
}

func markJobInactive(activeJobs *sync.Map, jobID int) {
	activeJobs.Delete(jobID)
}

func updatePeak(peak *int32, current int32) {
	for {
		oldPeak := atomic.LoadInt32(peak)

		if current <= oldPeak {
			return
		}

		if atomic.CompareAndSwapInt32(peak, oldPeak, current) {
			return
		}
	}
}

func processJob(jobID int, jobService *service.JobService, cxt context.Context, client *goredis.Client) (bool, error) {
	job, err := jobService.GetJob(jobID)
	if err != nil {
		return false, err
	}

	fmt.Println("Job Found:", job)

	attemptID, err := jobService.CreateAttempt(job.JobID)

	if err != nil {
		return false, err
	}

	fmt.Println("Attempt ID: ", attemptID)

	err = jobService.UpdateJob(job.JobID, "processing")

	if err != nil {
		return false, err
	}

	fmt.Println("Job is now Processing")

	if job.JobType == "docker-success-test" {
		if err := jobService.UpdateAttemptStatus(attemptID, "completed"); err != nil {
			return false, err
		}

		if err := jobService.UpdateJob(job.JobID, "completed"); err != nil {
			return false, err
		}

		fmt.Println("Job completed successfully:", job.JobID)
		return false, nil
	}

	err = errors.New("simulated failure")

	if err != nil {
		err = helper(job.JobID, attemptID, jobService)

		if err != nil {
			return false, err
		}
	}

	currentAttempt, err := jobService.Attempts(job.JobID)
	if err != nil {
		return false, err
	}

	fmt.Println(currentAttempt)

	got := jobService.CanRetry(currentAttempt, 3)

	if got {
		err = jobService.UpdateJob(job.JobID, "pending")
		if err != nil {
			return false, err
		}

		return true, nil
	}

	fmt.Println("Maximum attempts reached")
	return false, nil
}

func helper(jobID int, attemptID int, jobService *service.JobService) error {
	err := jobService.UpdateAttemptStatus(attemptID, "failed")

	if err != nil {
		return err
	}

	err = jobService.UpdateJob(jobID, "failed")

	if err != nil {
		return err
	}

	return nil
}
