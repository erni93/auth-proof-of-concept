package router

import (
	"authGo/session"
	"authGo/validator"
	"errors"
	"log"
	"net/http"
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

	tokens, err := v.CreateTokens(user)
	if err != nil {
		log.Print(err)
		writeGeneralError(w)
		return
	}

	l.SessionHandler.AddNewSession(tokens.RefreshPayload, v.GetDeviceData())

	writeSuccessfulLogin(w, tokens)
}
