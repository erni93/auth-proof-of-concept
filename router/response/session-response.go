package router

import (
	"authGo/session"
	"encoding/json"
	"net/http"
)

type SessionResponse struct {
	Sessions []*session.Session `json:"sessions"`
}

func WriteSessionList(w http.ResponseWriter, sessions []*session.Session) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SessionResponse{Sessions: sessions})
}
