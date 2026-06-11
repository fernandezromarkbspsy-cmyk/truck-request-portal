package middleware

import (
	"net/http"
	"truck-request-portal/pkg/database"
)

// RequireRole checks if the authenticated user has the required role
func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Extract Clerk User ID from the request context
			// (Assuming clerk-go middleware has already verified the JWT and placed the UserID in context)
			clerkUserID, ok := r.Context().Value("clerkUserID").(string)
			if !ok || clerkUserID == "" {
				http.Error(w, "Unauthorized: Missing Clerk User ID", http.StatusUnauthorized)
				return
			}

			// 2. Fetch user role from database (Prevents spoofing)
			var role string
			query := `SELECT role FROM users WHERE clerk_id = $1 AND is_active = true`
			err := database.DB.QueryRow(r.Context(), query, clerkUserID).Scan(&role)
			if err != nil {
				http.Error(w, "Unauthorized: User not found or inactive", http.StatusForbidden)
				return
			}

			// 3. Check if the user's role matches the required role
			// Note: In a real app, you might want a hierarchy (e.g., fte_ops can do ops_pic things)
			if role != requiredRole {
				http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
				return
			}

			// 4. If authorized, proceed to the next handler (Controller)
			next.ServeHTTP(w, r)
		})
	}
}
