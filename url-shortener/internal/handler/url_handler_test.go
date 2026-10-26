package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeURLShortener struct {
	code        string
	err         error
	called      bool
	receivedURL string

	originalURL string
	lookupErr   error
	lookupCode  string
}

func (f *fakeURLShortener) ShortenURL(originalURL string) (string, error) {
	f.called = true
	f.receivedURL = originalURL
	return f.code, f.err
}

func (f *fakeURLShortener) GetOriginalURL(shortCode string) (string, error) {
	f.lookupCode = shortCode
	if f.lookupErr != nil {
		return "", f.lookupErr
	}
	return f.originalURL, nil
}

func TestShortenValidRequest(t *testing.T) {
	fake := &fakeURLShortener{code: "abc123"}
	h := NewURLHandler(fake)

	body := `{"url":"https://example.com/products?id=123"}`
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.Shorten(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var response ShortenResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.ShortCode != "abc123" {
		t.Errorf("expected code abc123, got %q", response.ShortCode)
	}

	if response.ShortURL != "http://localhost:8080/abc123" {
		t.Errorf("unexpected short URL: %q", response.ShortURL)
	}

	if !fake.called {
		t.Error("expected service to be called")
	}

	if fake.receivedURL != "https://example.com/products?id=123" {
		t.Errorf("unexpected original URL: %q", fake.receivedURL)
	}
}

func TestShortenRejectsInvalidURL(t *testing.T) {
	fake := &fakeURLShortener{code: "abc123"}
	h := NewURLHandler(fake)

	req := httptest.NewRequest(
		http.MethodPost,
		"/shorten",
		strings.NewReader(`{"url":"not-a-valid-url"}`),
	)
	rec := httptest.NewRecorder()

	h.Shorten(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if fake.called {
		t.Error("service should not be called for an invalid URL")
	}
}

func TestShortenRejectsWrongMethod(t *testing.T) {
	fake := &fakeURLShortener{code: "abc123"}
	h := NewURLHandler(fake)

	req := httptest.NewRequest(http.MethodGet, "/shorten", nil)
	rec := httptest.NewRecorder()

	h.Shorten(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}

	if fake.called {
		t.Error("service should not be called for the wrong method")
	}
}

func TestRedirectValidShortCode(t *testing.T) {
	fake := &fakeURLShortener{
		originalURL: "https://example.com/products?id=123",
	}
	h := NewURLHandler(fake)

	req := httptest.NewRequest(http.MethodGet, "/gJXpSZ", nil)
	rec := httptest.NewRecorder()

	h.Redirect(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected status 302, got %d", rec.Code)
	}

	want := "https://example.com/products?id=123"
	if got := rec.Header().Get("Location"); got != want {
		t.Errorf("expected Location %q, got %q", want, got)
	}

	if fake.lookupCode != "gJXpSZ" {
		t.Errorf("expected lookup code gJXpSZ, got %q", fake.lookupCode)
	}
}

func TestRedirectUnknownShortCode(t *testing.T) {
	fake := &fakeURLShortener{
		lookupErr: sql.ErrNoRows,
	}
	h := NewURLHandler(fake)

	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	rec := httptest.NewRecorder()

	h.Redirect(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
}

func TestRedirectRejectsWrongMethod(t *testing.T) {
	fake := &fakeURLShortener{}
	h := NewURLHandler(fake)

	req := httptest.NewRequest(http.MethodPost, "/gJXpSZ", nil)
	rec := httptest.NewRecorder()

	h.Redirect(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rec.Code)
	}
}

func TestShortenRejectsMalformedJSON(t *testing.T) {
	fake := &fakeURLShortener{code: "abc123"}
	h := NewURLHandler(fake)

	req := httptest.NewRequest(
		http.MethodPost,
		"/shorten",
		strings.NewReader(`{"url":`),
	)
	rec := httptest.NewRecorder()

	h.Shorten(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if fake.called {
		t.Error("service should not be called for malformed JSON")
	}
}

func TestShortenRejectsExtraJSONValue(t *testing.T) {
	fake := &fakeURLShortener{code: "abc123"}
	h := NewURLHandler(fake)

	body := `{"url":"https://example.com"} {"extra":true}`
	req := httptest.NewRequest(
		http.MethodPost,
		"/shorten",
		strings.NewReader(body),
	)
	rec := httptest.NewRecorder()

	h.Shorten(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if fake.called {
		t.Error("service should not be called for extra JSON values")
	}
}

func TestShortenRejectsOversizedBody(t *testing.T) {
	fake := &fakeURLShortener{code: "abc123"}
	h := NewURLHandler(fake)

	body := `{"url":"https://example.com","padding":"` +
		strings.Repeat("a", 1<<20) + `"}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/shorten",
		strings.NewReader(body),
	)
	rec := httptest.NewRecorder()

	h.Shorten(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413, got %d", rec.Code)
	}

	if fake.called {
		t.Error("service should not be called for an oversized body")
	}
}

func TestShortenReturnsInternalServerError(t *testing.T) {
	fake := &fakeURLShortener{
		err: errors.New("database unavailable"),
	}
	h := NewURLHandler(fake)

	req := httptest.NewRequest(
		http.MethodPost,
		"/shorten",
		strings.NewReader(`{"url":"https://example.com"}`),
	)
	rec := httptest.NewRecorder()

	h.Shorten(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}

	if !fake.called {
		t.Error("expected service to be called")
	}
}

func TestRedirectReturnsInternalServerError(t *testing.T) {
	fake := &fakeURLShortener{
		lookupErr: errors.New("database unavailable"),
	}
	h := NewURLHandler(fake)

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	rec := httptest.NewRecorder()

	h.Redirect(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
}
