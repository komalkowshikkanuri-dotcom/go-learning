package repository

import (
	"database/sql"
	"errors"

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
