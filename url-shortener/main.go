package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"url-shortener/internal/cache"
	"url-shortener/internal/handler"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"

	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

func main() {
	connStr := os.Getenv("URL_SHORTENER_DB_URL")
	if connStr == "" {
		log.Fatal("URL_SHORTENER_DB_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Open database: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Ping database: ", err)
	}

	repo := repository.NewURLRepository(db)

	// Configure Redis.
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	urlCache := cache.NewURLCache(redisAddr, 15*time.Minute)

	// Verify Redis connectivity, but keep PostgreSQL available if Redis is down.
	redisCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := urlCache.Get(redisCtx, "__health_check__"); err != nil &&
		err != redis.Nil {
		log.Printf("Redis unavailable; continuing without reliable cache reads: %v", err)
	}

	// Inject Redis cache into the service.
	urlService := service.NewURLServiceWithCache(repo, urlCache)
	urlHandler := handler.NewURLHandler(urlService)

	defer urlCache.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", urlHandler.Shorten)
	mux.HandleFunc("/", urlHandler.Redirect)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("URL shortener listening on http://localhost:8080")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("HTTP server: ", err)
	}
}
