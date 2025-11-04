package database

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func getEnvInt(key string, def int32) int32 {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return int32(n)
		}
	}
	return def
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}

// NewPool creates a pgx connection pool with sane defaults and env overrides.
// Env overrides (optional):
//
//	DB_MAX_CONNS (int), DB_MIN_CONNS (int), DB_MAX_CONN_LIFETIME (duration),
//	DB_MAX_CONN_IDLE_TIME (duration), DB_HEALTH_CHECK_PERIOD (duration)
func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	// Defaults with env overrides
	cfg.MaxConns = getEnvInt("DB_MAX_CONNS", 10)
	cfg.MinConns = getEnvInt("DB_MIN_CONNS", 0)
	cfg.MaxConnLifetime = getEnvDuration("DB_MAX_CONN_LIFETIME", 30*time.Minute)
	cfg.MaxConnIdleTime = getEnvDuration("DB_MAX_CONN_IDLE_TIME", 5*time.Minute)
	cfg.HealthCheckPeriod = getEnvDuration("DB_HEALTH_CHECK_PERIOD", 30*time.Second)

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// Ensure we can connect now rather than failing on first request
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}
