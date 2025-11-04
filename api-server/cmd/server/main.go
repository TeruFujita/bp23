package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"line_bot/api-server/pkg/database"
)

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	// Load environment variables from .env files if present (development use)
	_ = godotenv.Load(".env.local")
	_ = godotenv.Load("../.env.local")

	r := gin.Default()

	allowedOrigin := getEnv("ALLOWED_ORIGIN", "http://localhost:3000")
	cfg := cors.Config{
		AllowOrigins:     []string{allowedOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: false,
	}
	r.Use(cors.New(cfg))

	// Database
	dbURL := os.Getenv("DATABASE_URL")
	var pool *pgxpool.Pool
	if dbURL != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		p, err := database.NewPool(ctx, dbURL)
		if err == nil {
			pool = p
		} else {
			// If DB init failed, we still start the server but healthz will show degraded
		}
	}

	r.GET("/healthz", func(c *gin.Context) {
		// Base status
		resp := gin.H{"status": "ok"}

		// DB check if configured
		if pool != nil {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
			defer cancel()
			var one int
			if err := pool.QueryRow(ctx, "select 1").Scan(&one); err != nil || one != 1 {
				if err == nil {
					err = pgx.ErrNoRows
				}
				c.JSON(http.StatusServiceUnavailable, gin.H{"status": "degraded", "db": err.Error()})
				return
			}
			resp["db"] = "ok"
		} else {
			resp["db"] = "disabled"
		}
		c.JSON(http.StatusOK, resp)
	})

	_ = r.Run(":8080")
}
