package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/komalkowshikkanuri/distributed-job-system/internal/models"
	"github.com/komalkowshikkanuri/distributed-job-system/internal/service"
)

type JobHandler struct {
	service *service.JobService
}

type UpdateJobRequest struct {
	Status string `json:"status"`
}

func NewJobHandler(h *service.JobService) *JobHandler {
	return &JobHandler{
		service: h,
	}
}

func (h *JobHandler) GetJob(jobID int) (models.Job, error) {
	job, err := h.service.GetJob(jobID)

	if err != nil {
		return job, err
	}

	return job, nil
}

func (h *JobHandler) GetJobHTTP(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	jobID, err := strconv.Atoi(parts[len(parts)-1])

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	job, err := h.GetJob(jobID)

	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(job)
	}
}

func (h *JobHandler) UpdateJobHTTP(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	jobID, err := strconv.Atoi(parts[len(parts)-1])

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var request UpdateJobRequest

	err = json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateJob(jobID, request.Status)

	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *JobHandler) JobHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetJobHTTP(w, r)
	case http.MethodPatch:
		h.UpdateJobHTTP(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
