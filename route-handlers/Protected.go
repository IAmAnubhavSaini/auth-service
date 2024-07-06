package route_handlers

import (
	"encoding/json"
	"net/http"
)

func ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(map[string]string{"message": "Welcome to the protected endpoint, " + username + "!"})
	if err != nil {
		panic(err)
	}
}
