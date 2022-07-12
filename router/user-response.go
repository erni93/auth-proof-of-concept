package router

import (
	"authGo/user"
	"encoding/json"
	"net/http"
)

type UserResponse struct {
	Users []*user.User `json:"users"`
}

func writeUserList(w http.ResponseWriter, users []*user.User) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UserResponse{Users: users})
}
