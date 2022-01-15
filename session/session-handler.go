package session

import (
	"authGo/authentication"
	"errors"
)

var ErrSessionAlreadyExists = errors.New("session handler: session already exists")
var ErrUserTokenNotFound = errors.New("session handler: user token not found")
var ErrRenewUserTokenDifferent = errors.New("session handler: new user token has a different user id")

type SessionsHandler struct {
	sessions []*Session
}

func NewSessionHandler() *SessionsHandler {
	return &SessionsHandler{sessions: make([]*Session, 0)}
}

func (s *SessionsHandler) getSession(userToken authentication.TokenData) *Session {
	for _, session := range s.sessions {
		if userToken.UserId == session.UserToken.UserId && userToken.IssuedAtTime.Equal(session.UserToken.IssuedAtTime) {
			return session
		}
	}
	return nil
}

func (s *SessionsHandler) GetUserSessions(userId string) []*Session {
	sessions := make([]*Session, 0)
	for _, session := range s.sessions {
		if session.UserToken.UserId == userId {
			sessions = append(sessions, session)
		}
	}
	return sessions
}

func (s *SessionsHandler) AddSession(session *Session) error {
	if s.getSession(session.UserToken) != nil {
		return ErrSessionAlreadyExists
	}
	s.sessions = append(s.sessions, session)
	return nil
}

func (s *SessionsHandler) DeleteSession(userToken authentication.TokenData) error {
	for i, v := range s.sessions {
		if v.UserToken == userToken {
			lastIndex := len(s.sessions) - 1
			s.sessions[i] = s.sessions[lastIndex]
			s.sessions[lastIndex] = nil
			s.sessions = s.sessions[:lastIndex]
			return nil
		}
	}
	return ErrUserTokenNotFound
}

func (s *SessionsHandler) RenewUserToken(oldData authentication.TokenData, newData authentication.TokenData) error {
	session := s.getSession(oldData)
	if session == nil {
		return ErrUserTokenNotFound
	}
	if session.UserToken.UserId != newData.UserId {
		return ErrRenewUserTokenDifferent
	}
	session.UserToken = newData
	return nil
}
