package session

import (
	"authGo/token"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSessionAlreadyExists    = errors.New("session handler: session already exists")
	ErrUserTokenNotFound       = errors.New("session handler: user token not found")
	ErrRenewUserTokenDifferent = errors.New("session handler: new user token has a different user id")
)

type SessionsHandler struct {
	sessions []*Session
}

func NewSessionHandler() *SessionsHandler {
	return &SessionsHandler{sessions: make([]*Session, 0)}
}

func (s *SessionsHandler) GetSession(userToken token.RefreshTokenPayload) (*Session, int, error) {
	for i, session := range s.sessions {
		if userToken.UserId == session.UserToken.UserId && userToken.IssuedAtTime.Equal(session.UserToken.IssuedAtTime) {
			return session, i, nil
		}
	}
	return nil, -1, ErrUserTokenNotFound
}

func (s *SessionsHandler) GetSessionById(id string) (*Session, error) {
	for _, session := range s.sessions {
		if session.Id == id {
			return session, nil
		}
	}
	return nil, ErrUserTokenNotFound
}

func (s *SessionsHandler) GetAllSessions() []*Session {
	return s.sessions
}

func (s *SessionsHandler) AddNewSession(userToken token.RefreshTokenPayload, deviceData DeviceData) error {
	if _, i, _ := s.GetSession(userToken); i != -1 {
		return ErrSessionAlreadyExists
	}
	id := strings.ReplaceAll(uuid.NewString(), "-", "")
	s.sessions = append(s.sessions, &Session{
		Id: id, UserToken: userToken, DeviceData: deviceData, LastUpdate: userToken.IssuedAtTime,
	})
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

func (s *SessionsHandler) DeleteSession(userToken token.RefreshTokenPayload) error {
	_, i, err := s.GetSession(userToken)
	if err != nil {
		return err
	}
	lastIndex := len(s.sessions) - 1
	s.sessions[i] = s.sessions[lastIndex]
	s.sessions[lastIndex] = nil
	s.sessions = s.sessions[:lastIndex]
	return nil
}

func (s *SessionsHandler) RefreshLastUpdate(session *Session) {
	session.LastUpdate = time.Now()
}
