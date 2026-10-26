package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
}

type URLShortener interface {
	ShortenURL(originalURL string) (string, error)
	GetOriginalURL(shortCode string) (string, error)
}

type URLHandler struct {
	service URLShortener
}

func NewURLHandler(service URLShortener) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req ShortenRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			return
		}

		http.Error(w, "invalid JSON request", http.StatusBadRequest)
		return
	}

	// Reject extra JSON values after the first object.

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			return
		}

		http.Error(w, "request must contain one JSON object", http.StatusBadRequest)
		return
	}

	originalURL := strings.TrimSpace(req.URL)
	parsedURL, err := url.ParseRequestURI(originalURL)
	if err != nil || parsedURL.Host == "" ||
		(parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		http.Error(w, "url must be a valid HTTP or HTTPS URL", http.StatusBadRequest)
		return
	}

	code, err := h.service.ShortenURL(originalURL)
	if err != nil {
		http.Error(w, "could not shorten URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(ShortenResponse{
		ShortCode: code,
		ShortURL:  "http://localhost:8080/" + code,
	})
}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	shortCode := strings.TrimPrefix(r.URL.Path, "/")
	if shortCode == "" || strings.Contains(shortCode, "/") {
		http.NotFound(w, r)
		return
	}

	originalURL, err := h.service.GetOriginalURL(shortCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "short URL not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}
