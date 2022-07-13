package router

import (
	response "authGo/router/response"
	"authGo/session"
	"authGo/validator"
	"log"
	"net/http"
)

type SessionRouter struct {
	Services *validator.Services
}

func (s *SessionRouter) GetSessionsHandler(w http.ResponseWriter, r *http.Request) {
	v := validator.AccessTokenValidator{Validator: validator.Validator{Writer: w, Request: r, Services: s.Services}}

	payload, err := v.ValidateAccessToken()
	if err != nil {
		log.Print(err)
		response.WriteTokenError(w)
		return
	}

	var sessions []*session.Session
	if payload.IsAdmin {
		sessions = s.Services.SessionsHandler.GetAllSessions()
	} else {
		sessions = s.Services.SessionsHandler.GetUserSessions(payload.UserId)
	}

	response.WriteSessionList(w, sessions)
}
