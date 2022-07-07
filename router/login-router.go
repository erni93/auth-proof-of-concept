package router

import (
	"authGo/session"
	"authGo/token"
	"authGo/user"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type LoginRouter struct {
	UserService           *user.UserService
	SessionHandler        *session.SessionsHandler
	AccessTokenGenerator  *token.TokenGenerator[token.AccessTokenPayload]
	RefreshTokenGenerator *token.TokenGenerator[token.RefreshTokenPayload]
}

func (l *LoginRouter) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	vars := mux.Vars(r)
	name := vars["name"]
	password := vars["password"]
	if name == "" || password == "" {
		writeError(w, "Empty name or password")
	}

	u, err := l.UserService.GetRepository().GetByName(name)
	if err == user.ErrUserNotFound {
		writeError(w, "User not found")
	}

	if isPasswordValid := l.UserService.IsPasswordValid(name, password); !isPasswordValid {
		writeError(w, "Password not valid")
	}

	now := time.Now()

	accessPayload := &token.AccessTokenPayload{UserId: u.Id, IssuedAtTime: now, IsAdmin: u.IsAdmin}
	accessJWT, err := l.AccessTokenGenerator.CreateToken(accessPayload)
	if err != nil {
		log.Printf("error creating accessToken jwt for user %s: %s", name, err)
		writeGeneralError(w)
	}

	refreshPayload := &token.RefreshTokenPayload{UserId: u.Id, IssuedAtTime: now}
	refreshJWT, err := l.RefreshTokenGenerator.CreateToken(refreshPayload)
	if err != nil {
		log.Printf("error creating refreshToken jwt for user %s: %s", name, err)
		writeGeneralError(w)
	}

	newSession := &session.Session{UserToken: *refreshPayload, DeviceData: getDeviceData(r), Created: time.Now()}
	l.SessionHandler.AddSession(newSession)

	writeSuccessLogin(w, accessJWT, refreshJWT)
}

func getDeviceData(r *http.Request) session.DeviceData {
	return session.DeviceData{IpAddress: r.RemoteAddr, UserAgent: r.UserAgent()}
}

func writeSuccessLogin(w http.ResponseWriter, accessJWT string, refreshJWT string) {
	w.WriteHeader(http.StatusOK)
	accessCookie := &http.Cookie{Name: "accessToken", Value: accessJWT, HttpOnly: true}
	refreshCookie := &http.Cookie{Name: "refreshToken", Value: refreshJWT, HttpOnly: true}
	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)
}
