package authentication

import (
	"authGo/authentication"
	"errors"
	"time"
)

var ErrSessionAlreadyExists = errors.New("session handler: session already exists")
var ErrRefreshTokenNotFound = errors.New("session handler: refreshToken not found")

type Session struct {
	RefreshToken authentication.Token
	DeviceData   DeviceData
	Created      time.Time
}

type SessionsHandler struct {
	sessions []*Session
}

func NewSessionHandler() *SessionsHandler {
	return &SessionsHandler{sessions: make([]*Session, 0)}
}

func (s *SessionsHandler) Add(session *Session) error {
	if s.ExistToken(session.RefreshToken) {
		return ErrSessionAlreadyExists
	}
	s.sessions = append(s.sessions, session)
	return nil
}

func (s *SessionsHandler) ExistToken(refreshToken authentication.Token) bool {
	for _, session := range s.sessions {
		if refreshToken == session.RefreshToken {
			return true
		}
	}
	return false
}

func (s *SessionsHandler) DeleteSession(refreshToken authentication.Token) error {
	for i, v := range s.sessions {
		if v.RefreshToken == refreshToken {
			lastIndex := len(s.sessions) - 1
			s.sessions[i] = s.sessions[lastIndex]
			s.sessions[lastIndex] = nil
			s.sessions = s.sessions[:lastIndex]
			return nil
		}
	}
	return ErrRefreshTokenNotFound
}
