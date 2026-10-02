package models

type Job struct {
	JobID   int    `json:"job_id"`
	JobType string `json:"type"`
	Status  string `json:"status"`
}
