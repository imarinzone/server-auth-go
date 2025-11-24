package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"server-auth-go/internal/auth"
	"server-auth-go/internal/token"
	"server-auth-go/pkg/middleware"
)

func main() {
	// Configuration
	secretKey := getEnv("SECRET_KEY", "super-secret-key-change-me")
	issuer := getEnv("ISSUER", "auth-server")
	port := getEnv("PORT", "8080")
	pgConn := getEnv("POSTGRES_CONN", "host=localhost user=postgres password=postgres dbname=auth sslmode=disable")
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")

	// Initialize services
	tokenService := token.NewService(secretKey, issuer)

	// Initialize Store (Postgres + Redis Cache)
	var store auth.Store
	pgStore, err := auth.NewPostgresStore(pgConn)
	if err != nil {
		log.Printf("Failed to connect to Postgres, falling back to in-memory (for dev only): %v", err)
		store = auth.NewInMemoryStore()
	} else {
		log.Println("Connected to Postgres")
		// Seed for demo
		pgStore.AddClient("service-a", "secret-a")
		pgStore.AddClient("service-b", "secret-b")

		// Add Redis Cache
		redisStore, err := auth.NewRedisStore(redisAddr, "", 0, pgStore)
		if err != nil {
			log.Printf("Failed to connect to Redis, using Postgres only: %v", err)
			store = pgStore
		} else {
			log.Println("Connected to Redis")
			store = redisStore
		}
	}

	authHandler := auth.NewHandler(store, tokenService)

	// Setup routes
	mux := http.NewServeMux()

	// Public endpoint: Get Token
	mux.HandleFunc("/token", authHandler.HandleToken)

	// Protected endpoint: Example resource
	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(middleware.ClaimsContextKey)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Access granted to protected resource",
			"claims":  claims,
		})
	})

	// Apply middleware to protected route
	mux.Handle("/protected", middleware.AuthMiddleware(tokenService)(protectedHandler))

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
