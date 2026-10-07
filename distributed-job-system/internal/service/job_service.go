package service

import (
	"errors"

	"github.com/komalkowshikkanuri/distributed-job-system/internal/models"
	"github.com/komalkowshikkanuri/distributed-job-system/internal/repository"
)

type JobService struct {
	repo *repository.JobRepository
}

func NewJobService(pointer *repository.JobRepository) *JobService {
	return &JobService{
		repo: pointer,
	}
}

func (s *JobService) GetJob(jobID int) (models.Job, error) {
	job, err := s.repo.GetJob(jobID)

	if err != nil {
		return job, err
	}

	return job, nil
}

func (s *JobService) UpdateJob(jobID int, status string) error {
	validStatuses := map[string]bool{
		"pending":    true,
		"processing": true,
		"completed":  true,
		"failed":     true,
	}

	if !validStatuses[status] {
		return errors.New("invalid job status")
	}

	currentJob, err := s.repo.GetJob(jobID)

	if err != nil {
		return err
	}

	if currentJob.Status == "pending" && status != "processing" {
		return errors.New("invalid job status transition")
	}

	if currentJob.Status == "processing" && status != "completed" && status != "failed" {
		return errors.New("invalid job status transition")
	}

	if currentJob.Status == "failed" && status != "pending" {
		return errors.New("invalid job status transition")
	}

	if currentJob.Status == "completed" {
		return errors.New("invalid job status transition")
	}

	return s.repo.UpdateJob(jobID, status)
}
