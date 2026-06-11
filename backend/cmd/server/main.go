package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	// Internal packages
	"truck-request-portal/features/clusters"
	"truck-request-portal/features/requests"
	"truck-request-portal/features/users"
	"truck-request-portal/pkg/cache"
	"truck-request-portal/pkg/database"
	appmiddleware "truck-request-portal/pkg/middleware"
)

func main() {
	// 1. Load Environment Variables (Rule: Never hardcode secrets)
	if !loadEnvFile(".env", "backend/.env", "../../.env") {
		log.Println("Warning: .env file not found, relying on system environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback default
	}

	// 2. Initialize Database (Supabase/PostgreSQL)
	if err := database.InitDB(); err != nil {
		log.Fatalf("Fatal: Could not connect to database: %v", err)
	}
	defer database.DB.Close()
	log.Println("✅ Database connected successfully")

	// 3. Initialize Cache (Upstash Redis)
	if err := cache.InitRedis(); err != nil {
		log.Fatalf("Fatal: Could not connect to Redis: %v", err)
	}
	defer cache.RedisClient.Close()
	log.Println("✅ Redis cache connected successfully")

	// 4. Initialize Chi Router
	r := chi.NewRouter()

	// 5. Global Middleware (Performance & Security)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)                    // Logs all HTTP requests
	r.Use(chimiddleware.Recoverer)                 // Prevents server crashes on panics
	r.Use(chimiddleware.Timeout(60 * time.Second)) // Prevents hanging requests

	// CORS Configuration (Allow Frontend to communicate with Backend)
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		log.Println("Warning: FRONTEND_URL is not set; only *.vercel.app origins will be allowed")
	}
	allowedOrigins := []string{"https://*.vercel.app"}
	if frontendURL != "" {
		allowedOrigins = append(allowedOrigins, frontendURL)
	}
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// 6. Public Routes (No Auth Required)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is running and healthy"))
	})

	// Clerk Webhook (Requires Clerk webhook signature verification in production)
	r.Post("/api/v1/users/webhook", users.HandleClerkWebhook)

	// 7. Protected API Routes (Requires Valid Clerk JWT)
	r.Route("/api/v1", func(r chi.Router) {
		// Apply Auth Middleware to all routes in this block
		r.Use(appmiddleware.RequireAuth)

		// Cluster Routes (Cached, Read-only for authenticated users)
		r.Get("/clusters", clusters.HandleGetClusters)

		// Request Routes
		r.Route("/requests", func(r chi.Router) {
			r.Post("/", requests.HandleCreateRequest)
			r.Get("/", requests.HandleGetPendingRequests)

			// RBAC: Only 'fte_ops' can approve
			r.With(appmiddleware.RequireRole("fte_ops")).Put("/{id}/approve", requests.HandleApproveRequest)

			// RBAC: Only 'fte_mm' can assign or reject
			r.With(appmiddleware.RequireRole("fte_mm")).Get("/approved", requests.HandleGetApprovedRequests)
			r.With(appmiddleware.RequireRole("fte_mm")).Put("/{id}/assign", requests.HandleAssignTruck)
			r.With(appmiddleware.RequireRole("fte_mm")).Put("/{id}/reject", requests.HandleRejectRequest)
		})
	})

	// 8. Start Server with Graceful Shutdown (Reliability Principle)
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("🚀 Server is listening on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Wait for interrupt signal (Ctrl+C) to gracefully shut down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited properly")
}

func loadEnvFile(paths ...string) bool {
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Load(path); err != nil {
				log.Printf("Warning: failed to load %s: %v", path, err)
				return false
			}
			return true
		}
	}
	return false
}
