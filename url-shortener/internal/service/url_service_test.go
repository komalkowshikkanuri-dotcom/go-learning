package service

import (
	//"error"
	"context"
	"fmt"
	"testing"
	//"time"

	//"github.com/redis/go-redis/v9"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type fakeURLRepository struct {
	originalURL  string
	createErr    error
	savedCode    string
	savedURL     string
	calls        int
	createErrors []error
	getErr       error
	getCalls     int
}

type fakeURLCache struct {
	data     map[string]string
	getErr   error
	setErr   error
	getCalls int
	setCalls int
}

func (f *fakeURLCache) Get(ctx context.Context, shortCode string) (string, error) {
	f.getCalls++

	if f.getErr != nil {
		return "", f.getErr
	}

	value, ok := f.data[shortCode]
	if !ok {
		return "", redis.Nil
	}

	return value, nil
}

func (f *fakeURLCache) Set(ctx context.Context, shortCode, originalURL string) error {
	f.setCalls++

	if f.setErr != nil {
		return f.setErr
	}

	if f.data == nil {
		f.data = make(map[string]string)
	}

	f.data[shortCode] = originalURL
	return nil
}

func (f *fakeURLRepository) Create(shortCode, originalURL string) error {
	f.calls++
	f.savedCode = shortCode
	f.savedURL = originalURL
	if len(f.createErrors) >= f.calls {
		return f.createErrors[f.calls-1]
	}

	return f.createErr
}

func (f *fakeURLRepository) GetOriginalURL(shortCode string) (string, error) {
	f.getCalls++
	if f.getErr != nil {
		return "", f.getErr
	}
	return f.originalURL, nil
}

func TestShortenURLSuccess(t *testing.T) {
	repo := &fakeURLRepository{}
	svc := NewURLService(repo)

	originalURL := "https://example.com/products"

	code, err := svc.ShortenURL(originalURL)
	if err != nil {
		t.Fatalf("ShortenURL returned an error: %v", err)
	}

	if len(code) != 6 {
		t.Fatalf("expected code length 6, got %d", len(code))
	}

	if repo.savedCode != code {
		t.Fatalf("saved code %q does not match returned code %q", repo.savedCode, code)
	}

	if repo.savedURL != originalURL {
		t.Fatalf("saved URL %q does not match original URL %q", repo.savedURL, originalURL)
	}

	if repo.calls != 1 {
		t.Fatalf("expected 1 repository call, got %d", repo.calls)
	}
}

func TestShortenURLRetriesDuplicateCode(t *testing.T) {
	duplicateErr := &pq.Error{Code: "23505"}

	repo := &fakeURLRepository{
		createErrors: []error{
			duplicateErr,
			nil,
		},
	}

	svc := NewURLService(repo)

	code, err := svc.ShortenURL("https://example.com/products")
	if err != nil {
		t.Fatalf("expected retry to succeed, got error: %v", err)
	}

	if len(code) != 6 {
		t.Fatalf("expected code length 6, got %d", len(code))
	}

	if repo.calls != 2 {
		t.Fatalf("expected 2 repository calls, got %d", repo.calls)
	}
}

func TestShortenURLReturnsRepositoryError(t *testing.T) {
	repo := &fakeURLRepository{
		createErr: fmt.Errorf("database unavailable"),
	}
	svc := NewURLService(repo)

	_, err := svc.ShortenURL("https://example.com/products")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if repo.calls != 1 {
		t.Fatalf("expected 1 repository call, got %d", repo.calls)
	}
}

func TestShortenURLFailsAfterMaxRetries(t *testing.T) {
	duplicateErr := &pq.Error{Code: "23505"}

	repo := &fakeURLRepository{
		createErr: duplicateErr,
	}
	svc := NewURLService(repo)

	_, err := svc.ShortenURL("https://example.com/products")
	if err == nil {
		t.Fatal("expected an error after exhausting retries, got nil")
	}

	if repo.calls != 5 {
		t.Fatalf("expected 5 repository calls, got %d", repo.calls)
	}
}

func TestGetOriginalURLReturnsRepositoryError(t *testing.T) {
	repo := &fakeURLRepository{
		getErr: fmt.Errorf("database unavailable"),
	}
	svc := NewURLService(repo)

	_, err := svc.GetOriginalURL("abc123")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestGetOriginalURLCacheHit(t *testing.T) {
	const shortCode = "abc123"
	const originalURL = "https://example.com/products"

	repo := &fakeURLRepository{}
	cache := &fakeURLCache{
		data: map[string]string{
			shortCode: originalURL,
		},
	}

	svc := NewURLServiceWithCache(repo, cache)

	got, err := svc.GetOriginalURL(shortCode)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got != originalURL {
		t.Fatalf("expected %q, got %q", originalURL, got)
	}

	if repo.getCalls != 0 {
		t.Fatalf("expected no PostgreSQL lookups, got %d", repo.getCalls)
	}

	if cache.getCalls != 1 {
		t.Fatalf("expected 1 cache lookup, got %d", cache.getCalls)
	}

	if cache.setCalls != 0 {
		t.Fatalf("expected no cache writes, got %d", cache.setCalls)
	}
}

func TestGetOriginalURLCacheMiss(t *testing.T) {
	const shortCode = "abc123"
	const originalURL = "https://example.com/products"

	repo := &fakeURLRepository{
		originalURL: originalURL,
	}
	cache := &fakeURLCache{
		data: make(map[string]string),
	}

	svc := NewURLServiceWithCache(repo, cache)

	got, err := svc.GetOriginalURL(shortCode)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got != originalURL {
		t.Fatalf("expected %q, got %q", originalURL, got)
	}

	if cache.getCalls != 1 {
		t.Fatalf("expected 1 cache lookup, got %d", cache.getCalls)
	}

	if repo.getCalls != 1 {
		t.Fatalf("expected 1 PostgreSQL lookup, got %d", repo.getCalls)
	}

	if cache.setCalls != 1 {
		t.Fatalf("expected 1 cache write, got %d", cache.setCalls)
	}

	if cache.data[shortCode] != originalURL {
		t.Fatalf("expected URL to be cached, got %q", cache.data[shortCode])
	}
}

func TestGetOriginalURLFallsBackWhenCacheFails(t *testing.T) {
	const shortCode = "abc123"
	const originalURL = "https://example.com/products"

	repo := &fakeURLRepository{
		originalURL: originalURL,
	}
	cache := &fakeURLCache{
		data:   make(map[string]string),
		getErr: fmt.Errorf("redis unavailable"),
	}

	svc := NewURLServiceWithCache(repo, cache)

	got, err := svc.GetOriginalURL(shortCode)
	if err != nil {
		t.Fatalf("expected fallback to succeed, got %v", err)
	}

	if got != originalURL {
		t.Fatalf("expected %q, got %q", originalURL, got)
	}

	if repo.getCalls != 1 {
		t.Fatalf("expected 1 repository lookup, got %d", repo.getCalls)
	}
}
