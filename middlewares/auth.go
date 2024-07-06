package middlewares

import (
	"auth-service/config"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
)

// Auth is a middlewares that checks if the request has a valid JWT token.
// If the token is valid, it sets the username in the request header.
// If the token is invalid, it returns a 401 Unauthorized error.
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, keyFunc())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			r.Header.Set("username", claims["username"].(string))
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// https://pkg.go.dev/github.com/golang-jwt/jwt#Parse
// Parse parses, validates, and returns a token.
// keyFunc will receive the parsed token and should return the key for validating.
// If everything is correct, a token will be returned, and err will be nil.
// If there is an error, the token will be nil, and err will contain information about what went wrong.
func keyFunc() func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.JwtSecret, nil
	}
}
