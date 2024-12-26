package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Define a struct to parse the JWT claims

type contextKey string

const UserKey contextKey = "user"

// Define a secret key for JWT signing (keep this secure in production)
var jwtSecret = []byte("KA11EL4943")

// Middleware to validate the JWT cookie
func CookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/storePIDetails" {
			// Pass the request directly to the next handler
			next.ServeHTTP(w, r)
			return
		}

		// Retrieve the session cookie
		cookie, err := r.Cookie("session")
		if err != nil {
			// No cookie: send JSON response to frontend
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "unauthorized",
				"msg":   "Session is missing or expired.",
			})
			return
		}

		// Validate the JWT in the cookie
		userEmail, isValid := validateJWT(cookie.Value)
		if !isValid {
			// Invalid or expired JWT
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "unauthorized",
				"msg":   "Session is invalid or expired.",
			})
			return
		}

		// Store the user email in the request context
		ctx := context.WithValue(r.Context(), UserKey, userEmail)

		// Pass the updated request to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Validate the JWT token in the cookie
func validateJWT(tokenString string) (string, bool) {
	// Parse and validate the JWT
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is as expected
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", false // Invalid token
	}

	// Extract claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || claims.ExpiresAt.Time.Before(time.Now()) {
		return "", false // Token expired or claims invalid
	}

	// Return the email
	return claims.Email, true
}
