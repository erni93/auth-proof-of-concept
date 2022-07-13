package router

import (
	response "authGo/router/response"
	"authGo/session"
	"authGo/validator"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

func (s *SessionRouter) DeleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	v := validator.AccessTokenValidator{Validator: validator.Validator{Writer: w, Request: r, Services: s.Services}}

	payload, err := v.ValidateAccessToken()
	if err != nil {
		log.Print(err)
		response.WriteTokenError(w)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	session, err := s.Services.SessionsHandler.GetSessionById(id)
	if err != nil {
		log.Print(err)
		response.WriteError(w, "Session not found")
		return
	}
	if !payload.IsAdmin && session.UserToken.UserId != payload.UserId {
		response.WriteForbidden(w)
	}

	err = s.Services.SessionsHandler.DeleteSession(session.UserToken)
	if err != nil {
		log.Print(err)
		response.WriteError(w, "There was an error deleting the session")
		return
	}

	w.WriteHeader(http.StatusOK)
}
