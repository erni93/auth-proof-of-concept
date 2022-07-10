package router

import (
	"authGo/token"
	"authGo/validator"
	"encoding/json"
	"net/http"
)

type AuthResponse struct {
	UserData token.AccessTokenPayload `json:"userData"`
}

func writeSuccessfulLogin(w http.ResponseWriter, tokens *validator.JwtTokens) {
	accessCookie := &http.Cookie{Name: "accessToken", Value: tokens.AccessToken, HttpOnly: true, Path: "/"}
	refreshCookie := &http.Cookie{Name: "refreshToken", Value: tokens.RefreshToken, HttpOnly: true, Path: "/"}
	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)
	json.NewEncoder(w).Encode(AuthResponse{UserData: tokens.AccessPayload})
}

func writeSuccessfulRefresh(w http.ResponseWriter, token *validator.AccessJwtToken) {
	accessCookie := &http.Cookie{Name: "accessToken", Value: token.AccessToken, HttpOnly: true, Path: "/"}
	http.SetCookie(w, accessCookie)
	json.NewEncoder(w).Encode(AuthResponse{UserData: token.AccessPayload})
}
