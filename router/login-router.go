package router

import (
	"authGo/session"
	"authGo/validator"
	"errors"
	"log"
	"net/http"
	"time"
)

type LoginRouter struct {
	Services       *validator.LoginRouterServices
	SessionHandler *session.SessionsHandler
}

func (l *LoginRouter) Handler(w http.ResponseWriter, r *http.Request) {
	v := validator.LoginRouterValidator{Writer: w, Request: r, Services: l.Services}

	loginDetails, err := v.GetLoginDetails()
	if err != nil {
		log.Print(err)
		if errors.Is(err, validator.ErrLoginRouterReadingFormData) {
			writeGeneralError(w)
		} else if errors.Is(err, validator.ErrLoginRouterEmptyNamePassword) {
			writeError(w, "Empty name or password")
		}
		return
	}

	user, err := v.GetUser(loginDetails)
	if err != nil {
		log.Print(err)
		if errors.Is(err, validator.ErrLoginRouterUserNotFound) {
			writeError(w, "User doesn't exist")
		} else if errors.Is(err, validator.ErrLoginRouterPasswordNotValid) {
			writeError(w, "Password not valid")
		}
		return
	}

	tokens, err := v.GetTokens(user)
	if err != nil {
		log.Print(err)
		writeGeneralError(w)
		return
	}

	newSession := &session.Session{UserToken: tokens.RefreshPayload, DeviceData: v.GetDeviceData(), Created: time.Now()}
	l.SessionHandler.AddSession(newSession)

	writeSuccessLogin(w, tokens)
}

func writeSuccessLogin(w http.ResponseWriter, tokens *validator.JwtTokens) {
	accessCookie := &http.Cookie{Name: "accessToken", Value: tokens.AccessToken, HttpOnly: true, Path: "/"}
	refreshCookie := &http.Cookie{Name: "refreshToken", Value: tokens.RefreshToken, HttpOnly: true, Path: "/"}
	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)
	w.WriteHeader(http.StatusOK)
}
