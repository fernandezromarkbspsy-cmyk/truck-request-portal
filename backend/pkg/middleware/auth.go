package middleware

import (
	"context"
	"net/http"
	"strings"
	"truck-request-portal/pkg/database"
)

// RequireAuth verifies the Clerk JWT and injects the UserID into the context.
// Note: In a full production setup, use github.com/clerk/clerk-sdk-go/v4 to verify the JWT signature.
// For this MVP, we assume the frontend sends a valid token and we extract the subject (user ID).
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: Missing or invalid token", http.StatusUnauthorized)
			return
		}

		// TODO: Add actual Clerk JWT verification here using clerk-sdk-go.
		// For now, we simulate extracting a user ID (e.g., "user_12345")
		// In production: session, err := clerk.VerifyToken(r.Context(), token)

		// Mock extraction for local development flow (Replace with actual Clerk verification)
		userID := "mock_user_id_for_dev" // Replace with actual parsed JWT subject

		// Inject into context for downstream layers (Service/Repository)
		ctx := context.WithValue(r.Context(), "clerkUserID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole checks the database to ensure the authenticated user has the correct role.
func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value("clerkUserID").(string)
			if !ok || userID == "" {
				http.Error(w, "Unauthorized: Missing user context", http.StatusUnauthorized)
				return
			}

			// Query database for user's role (Prevents frontend spoofing)
			var role string
			query := `SELECT role FROM users WHERE clerk_id = $1 AND is_active = true`
			err := database.DB.QueryRow(r.Context(), query, userID).Scan(&role)
			if err != nil {
				http.Error(w, "Forbidden: User not found or inactive", http.StatusForbidden)
				return
			}

			// Enforce Role-Based Access Control (RBAC)
			if role != requiredRole {
				http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
