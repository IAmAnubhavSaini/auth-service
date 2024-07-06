package route_handlers

import (
	"auth-service/config"
	t "auth-service/types"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var user t.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	storedPassword, exists := t.ServiceUsers[user.Username]
	if !exists || bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password)) != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * time.Duration(config.JwtExpiration)).Unix(),
	})

	tokenString, err := token.SignedString(config.JwtSecret)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	if err != nil {
		panic(err)
	}
}
