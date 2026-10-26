package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type URLRepository interface {
	Create(shortCode, originalURL string) error
	GetOriginalURL(shortCode string) (string, error)
}

type URLCache interface {
	Get(ctx context.Context, shortCode string) (string, error)
	Set(ctx context.Context, shortCode, originalURL string) error
}

type URLService struct {
	repo  URLRepository
	cache URLCache
}

func NewURLService(repo URLRepository) *URLService {
	return &URLService{
		repo: repo,
	}
}

func NewURLServiceWithCache(repo URLRepository, cache URLCache) *URLService {
	return &URLService{
		repo:  repo,
		cache: cache,
	}
}

func (s *URLService) ShortenURL(originalURL string) (string, error) {
	const maxAttempts = 5

	for attempt := 0; attempt < maxAttempts; attempt++ {
		code, err := GenerateShortCode(6)
		if err != nil {
			return "", err
		}

		err = s.repo.Create(code, originalURL)
		if err == nil {
			return code, nil
		}

		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			continue
		}

		return "", fmt.Errorf("save URL mapping: %w", err)
	}

	return "", fmt.Errorf("could not generate a unique short code after retries")
}

func (s *URLService) GetOriginalURL(shortCode string) (string, error) {
	ctx := context.Background()

	if s.cache != nil {
		originalURL, err := s.cache.Get(ctx, shortCode)
		if err == nil {
			return originalURL, nil
		}
		// A cache miss or Redis error falls back to PostgreSQL.
	}

	originalURL, err := s.repo.GetOriginalURL(shortCode)
	if err != nil {
		return "", fmt.Errorf("get original URL: %w", err)
	}

	if s.cache != nil {
		// A cache write failure must not break a successful lookup.
		_ = s.cache.Set(ctx, shortCode, originalURL)
	}

	return originalURL, nil
}
