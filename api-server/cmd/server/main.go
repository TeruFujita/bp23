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

	// ----- Users -----
	type createUserRequest struct {
		LineID string  `json:"line_id" binding:"required"`
		Name   *string `json:"name"`
	}

	r.POST("/api/users", func(c *gin.Context) {
		if pool == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "db not configured"})
			return
		}
		var req createUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()
		var id string
		// INSERT with upsert-like behavior: if line_id exists, return existing id
		// Simplicity: try insert, on conflict do nothing then select id
		if _, err := pool.Exec(ctx,
			"INSERT INTO users(line_id, name) VALUES ($1, $2) ON CONFLICT (line_id) DO NOTHING",
			req.LineID, req.Name,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := pool.QueryRow(ctx,
			"SELECT id FROM users WHERE line_id = $1",
			req.LineID,
		).Scan(&id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id})
	})

	r.GET("/api/users/:line_id", func(c *gin.Context) {
		if pool == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "db not configured"})
			return
		}
		lineID := c.Param("line_id")
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()
		var (
			id        string
			name      *string
			createdAt time.Time
		)
		if err := pool.QueryRow(ctx,
			"SELECT id, name, created_at FROM users WHERE line_id = $1",
			lineID,
		).Scan(&id, &name, &createdAt); err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id, "line_id": lineID, "name": name, "created_at": createdAt})
	})

	// ----- Routes -----
	type createRouteRequest struct {
		LineID   string  `json:"line_id" binding:"required"`
		StartLat float64 `json:"start_lat" binding:"required"`
		StartLng float64 `json:"start_lng" binding:"required"`
		EndLat   float64 `json:"end_lat" binding:"required"`
		EndLng   float64 `json:"end_lng" binding:"required"`
	}

	r.POST("/api/routes", func(c *gin.Context) {
		if pool == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "db not configured"})
			return
		}
		var req createRouteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()
		var routeID string
		// Insert route using subquery to fetch user_id by line_id
		sql := `INSERT INTO routes(user_id, start_lat, start_lng, end_lat, end_lng)
                VALUES ((SELECT id FROM users WHERE line_id = $1), $2, $3, $4, $5)
                RETURNING id`
		if err := pool.QueryRow(ctx, sql, req.LineID, req.StartLat, req.StartLng, req.EndLat, req.EndLng).Scan(&routeID); err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": routeID})
	})

	r.GET("/api/users/:line_id/routes", func(c *gin.Context) {
		if pool == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "db not configured"})
			return
		}
		lineID := c.Param("line_id")
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()
		rows, err := pool.Query(ctx,
			`SELECT r.id, r.start_lat, r.start_lng, r.end_lat, r.end_lng, r.created_at
            FROM routes r
            JOIN users u ON u.id = r.user_id
            WHERE u.line_id = $1
            ORDER BY r.created_at DESC`, lineID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		type routeDTO struct {
			ID        string    `json:"id"`
			StartLat  float64   `json:"start_lat"`
			StartLng  float64   `json:"start_lng"`
			EndLat    float64   `json:"end_lat"`
			EndLng    float64   `json:"end_lng"`
			CreatedAt time.Time `json:"created_at"`
		}
		var list []routeDTO
		for rows.Next() {
			var rDTO routeDTO
			if err := rows.Scan(&rDTO.ID, &rDTO.StartLat, &rDTO.StartLng, &rDTO.EndLat, &rDTO.EndLng, &rDTO.CreatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			list = append(list, rDTO)
		}
		c.JSON(http.StatusOK, gin.H{"routes": list})
	})

	// ----- Spots -----
	type createSpotRequest struct {
		RouteID  string   `json:"route_id" binding:"required"`
		Name     string   `json:"name" binding:"required"`
		Category string   `json:"category" binding:"required"`
		Lat      float64  `json:"lat" binding:"required"`
		Lng      float64  `json:"lng" binding:"required"`
		URL      *string  `json:"url"`
		Rating   *float64 `json:"rating"`
	}

	r.POST("/api/routes/:route_id/spots", func(c *gin.Context) {
		if pool == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "db not configured"})
			return
		}
		routeID := c.Param("route_id")
		var req createSpotRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()
		var id string
		sql := `INSERT INTO spots(route_id, name, category, lat, lng, url, rating)
                VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
		if err := pool.QueryRow(ctx, sql, routeID, req.Name, req.Category, req.Lat, req.Lng, req.URL, req.Rating).Scan(&id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id})
	})

	r.GET("/api/routes/:route_id/spots", func(c *gin.Context) {
		if pool == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "db not configured"})
			return
		}
		routeID := c.Param("route_id")
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()
		rows, err := pool.Query(ctx,
			`SELECT id, name, category, lat, lng, url, rating, created_at
            FROM spots WHERE route_id = $1 ORDER BY created_at DESC`, routeID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		type spotDTO struct {
			ID        string    `json:"id"`
			Name      string    `json:"name"`
			Category  string    `json:"category"`
			Lat       float64   `json:"lat"`
			Lng       float64   `json:"lng"`
			URL       *string   `json:"url"`
			Rating    *float64  `json:"rating"`
			CreatedAt time.Time `json:"created_at"`
		}
		var list []spotDTO
		for rows.Next() {
			var s spotDTO
			if err := rows.Scan(&s.ID, &s.Name, &s.Category, &s.Lat, &s.Lng, &s.URL, &s.Rating, &s.CreatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			list = append(list, s)
		}
		c.JSON(http.StatusOK, gin.H{"spots": list})
	})

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
