package router

import (
	"authGo/validator"
	"log"
	"net/http"
)

type RefreshRouter struct {
	Services *validator.Services
}

func (ref *RefreshRouter) Handler(w http.ResponseWriter, r *http.Request) {
	v := validator.RefreshValidator{Validator: validator.Validator{Writer: w, Request: r, Services: ref.Services}}

	session, err := v.ValidateRefreshToken()
	if err != nil {
		log.Print(err)
		writeError(w, "User session not found")
		return
	}

	user, err := v.Validator.Services.UserService.GetRepository().GetById(session.UserToken.UserId)
	if err != nil {
		log.Print(err)
		writeError(w, "User not valid")
		return
	}

	token, err := v.CreateAccessToken(user)
	if err != nil {
		log.Print(err)
		writeGeneralError(w)
		return
	}
	v.Validator.Services.SessionsHandler.RefreshLastUpdate(session)

	writeSuccessfulRefresh(w, token)
}
