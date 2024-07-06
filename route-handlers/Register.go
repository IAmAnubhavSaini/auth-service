package route_handlers

import (
	t "auth-service/types"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user t.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		panic(err)
	}

	t.ServiceUsers[user.Username] = string(hashedPassword)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
	if err != nil {
		panic(err)
	}
}
