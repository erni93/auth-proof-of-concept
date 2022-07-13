package router

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteError(w http.ResponseWriter, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorResponse{Error: err})
}

func WriteGeneralError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorResponse{Error: "Something went wrong. Please try again later"})
}

func WriteTokenError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}

func WriteForbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
}
