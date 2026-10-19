package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/komalkowshikkanuri/distributed-job-system/internal/models"
)

var ErrJobNotFound = errors.New("job not found")

type JobRepository struct {
	db *sql.DB
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{
		db: db,
	}
}

func (r *JobRepository) GetJob(jobID int) (models.Job, error) {
	var job models.Job

	row := r.db.QueryRow(
		`SELECT id, type, status
		FROM jobs
		WHERE id = $1`,
		jobID,
	)

	err := row.Scan(
		&job.JobID,
		&job.JobType,
		&job.Status,
	)

	if err != nil {
		return job, err
	}

	return job, nil
}

func (r *JobRepository) UpdateJob(jobID int, status string) error {
	result, err := r.db.Exec(
		`Update jobs
		SET status = $1
		WHERE id = $2`,
		status, jobID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrJobNotFound
	}

	return nil
}

func (r *JobRepository) CreateAttempt(jobID int) (int, error) {
	var attemptID int

	err := r.db.QueryRow(
		`INSERT INTO job_attempts (job_id, status)
			VALUES ($1, $2)
			RETURNING id;`,
		jobID, "processing",
	).Scan(&attemptID)

	if err != nil {
		return 0, err
	}

	return attemptID, nil
}

func (r *JobRepository) UpdateAttemptStatus(attemptID int, status string) error {
	res, err := r.db.Exec(
		`UPDATE job_attempts
		SET status = $1,
			finished_at = CASE
			   WHEN $2::text = 'completed' THEN NOW()
			   WHEN $3::text = 'failed' THEN NOW()
			   ELSE null
			END
		WHERE id = $4`,
		status, status, status, attemptID,
	)

	if err != nil {
		return err
	}

	n, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if n == 0 {
		return fmt.Errorf("attempt %d not found", attemptID)
	}

	return nil
}

func (r *JobRepository) Attempts(jobID int) (int, error) {
	var currentAttempt int

	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM job_attempts
		WHERE job_id = $1`,
		jobID,
	).Scan(&currentAttempt)

	if err != nil {
		return 0, err
	}

	return currentAttempt, nil
}
