package router

import (
	response "authGo/router/response"
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
		response.WriteError(w, "User session not found")
		return
	}

	user, err := v.Validator.Services.UserService.GetRepository().GetById(session.UserToken.UserId)
	if err != nil {
		log.Print(err)
		response.WriteError(w, "User not valid")
		return
	}

	token, err := v.CreateAccessToken(user)
	if err != nil {
		log.Print(err)
		response.WriteGeneralError(w)
		return
	}
	v.Validator.Services.SessionsHandler.RefreshLastUpdate(session)

	response.WriteSuccessfulRefresh(w, token)
}
