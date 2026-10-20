package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_SSLMODE", "")
	t.Setenv("REDIS_ADDR", "")

	cfg := Load()

	if cfg.DBHost != "127.0.0.1" {
		t.Errorf("DBHost = %q, want %q", cfg.DBHost, "127.0.0.1")
	}
	if cfg.DBPort != "5432" {
		t.Errorf("DBPort = %q, want %q", cfg.DBPort, "5432")
	}
	if cfg.DBName != "job_system" {
		t.Errorf("DBName = %q, want %q", cfg.DBName, "job_system")
	}
	if cfg.DBUser != "postgres" {
		t.Errorf("DBUser = %q, want %q", cfg.DBUser, "postgres")
	}
	if cfg.DBPassword != "" {
		t.Errorf("DBPassword = %q, want empty", cfg.DBPassword)
	}
	if cfg.DBSSLMode != "disable" {
		t.Errorf("DBSSLMode = %q, want %q", cfg.DBSSLMode, "disable")
	}
	if cfg.RedisAddr != "localhost:6379" {
		t.Errorf("RedisAddr = %q, want %q", cfg.RedisAddr, "localhost:6379")
	}
}

func TestLoadCustomValues(t *testing.T) {
	t.Setenv("DB_HOST", "db.example.com")
	t.Setenv("REDIS_ADDR", "redis.example.com:6379")

	cfg := Load()

	if cfg.DBHost != "db.example.com" {
		t.Errorf("DBHost = %q, want %q", cfg.DBHost, "db.example.com")
	}
	if cfg.RedisAddr != "redis.example.com:6379" {
		t.Errorf("RedisAddr = %q, want %q", cfg.RedisAddr, "redis.example.com:6379")
	}
}
