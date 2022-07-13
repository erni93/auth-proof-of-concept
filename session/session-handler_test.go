package session

import (
	"authGo/token"
	"reflect"
	"testing"
	"time"
)

func addTestSession(sessionHandler *SessionsHandler, userId string, issuedTime time.Time) error {
	err := sessionHandler.AddNewSession(
		token.RefreshTokenPayload{
			UserId:       userId,
			IssuedAtTime: issuedTime,
		},
		DeviceData{
			IpAddress: "ip-" + userId,
			UserAgent: "userAgent-" + userId,
		},
	)
	return err
}

func createTestSessionHandler(t *testing.T, issuedTime time.Time) *SessionsHandler {
	sessionHandler := NewSessionHandler()
	err := addTestSession(sessionHandler, "user1", issuedTime)
	if err != nil {
		t.Errorf("error adding session, %s", err)
	}
	err = addTestSession(sessionHandler, "user2", issuedTime)
	if err != nil {
		t.Errorf("error adding session, %s", err)
	}
	return sessionHandler
}

func TestNewSessionHandler(t *testing.T) {
	if sessionHandler := NewSessionHandler(); sessionHandler == nil {
		t.Error("expected NewSessionHandler to return SessionsHandler, got nil")
	}
}

func TestGetSession(t *testing.T) {
	now := time.Date(2022, 8, 6, 0, 0, 0, 0, time.UTC)
	sessionHandler := createTestSessionHandler(t, now)
	payload := token.RefreshTokenPayload{
		UserId:       "user1",
		IssuedAtTime: now,
	}
	session, id, err := sessionHandler.GetSession(payload)
	if err != nil {
		t.Errorf("expected err to be nil, got %s", err)
	}
	if session == nil {
		t.Error("expected session to not be nil")
	}
	if id == -1 {
		t.Error("expected session to not be -1")
	}
	if !reflect.DeepEqual(session.UserToken, payload) {
		t.Errorf("wanted %v to be %v", session.UserToken, payload)
	}

	payload2 := token.RefreshTokenPayload{
		UserId:       "user3",
		IssuedAtTime: now,
	}
	session, id, err = sessionHandler.GetSession(payload2)
	if err != ErrUserTokenNotFound {
		t.Errorf("expected err to be ErrUserTokenNotFound, got %s", err)
	}
	if session != nil {
		t.Error("expected session to be nil")
	}
	if id != -1 {
		t.Error("expected session to be -1")
	}
}

func TestGetSessionById(t *testing.T) {
	now := time.Date(2022, 8, 6, 0, 0, 0, 0, time.UTC)
	sessionHandler := createTestSessionHandler(t, now)
	sessionId := sessionHandler.sessions[0].Id
	session, err := sessionHandler.GetSessionById(sessionId)
	if err != nil {
		t.Errorf("expected err to be nil, got %s", err)
	}
	if session == nil {
		t.Error("expected session to not be nil")
	}
	session, err = sessionHandler.GetSessionById("12345678")
	if err != ErrUserTokenNotFound {
		t.Errorf("expected err to be ErrUserTokenNotFound, got %s", err)
	}
	if session != nil {
		t.Error("expected session to be nil")
	}
}

func TestAddNewSession(t *testing.T) {
	now := time.Date(2022, 8, 6, 0, 0, 0, 0, time.UTC)
	sessionHandler := createTestSessionHandler(t, now)
	err := addTestSession(sessionHandler, "user1", now)
	if err != ErrSessionAlreadyExists {
		t.Errorf("expected error to be ErrSessionAlreadyExists, got: %s", err)
	}
}

func TestGetUserSessions(t *testing.T) {
	sessionHandler := createTestSessionHandler(t, time.Now())

	user1Sessions := sessionHandler.GetUserSessions("user1")
	if len(user1Sessions) != 1 {
		t.Errorf("expected user1Sessions len to be 1, got %d", len(user1Sessions))
	}

	time.Sleep(1 * time.Millisecond)

	err := addTestSession(sessionHandler, "user1", time.Now())
	if err != nil {
		t.Errorf("error adding session, %s", err)
	}

	user1Sessions = sessionHandler.GetUserSessions("user1")
	if len(user1Sessions) != 2 {
		t.Errorf("expected user1Sessions len to be 2, got %d", len(user1Sessions))
	}
}

func TestDeleteSession(t *testing.T) {
	sessionHandler := createTestSessionHandler(t, time.Now())

	time.Sleep(1 * time.Millisecond)
	err := addTestSession(sessionHandler, "user1", time.Now())
	if err != nil {
		t.Errorf("error adding session, %s", err)
	}

	user1Sessions := sessionHandler.GetUserSessions("user1")
	if len(user1Sessions) != 2 {
		t.Errorf("expected user1Sessions len to be 2, got %d", len(user1Sessions))
	}

	err = sessionHandler.DeleteSession(user1Sessions[0].UserToken)
	if err != nil {
		t.Errorf("error deleting session, %s", err)
	}
	user1Sessions = sessionHandler.GetUserSessions("user1")
	if len(user1Sessions) != 1 {
		t.Errorf("expected user1Sessions len to be 1, got %d", len(user1Sessions))
	}

	err = sessionHandler.DeleteSession(token.RefreshTokenPayload{UserId: "user3", IssuedAtTime: time.Now()})
	if err != ErrUserTokenNotFound {
		t.Errorf("expected error to be ErrUserTokenNotFound, got: %s", err)
	}
}

func TestRefreshLastUpdate(t *testing.T) {
	now := time.Date(2022, 8, 6, 0, 0, 0, 0, time.UTC)
	sessionHandler := createTestSessionHandler(t, now)
	payload := token.RefreshTokenPayload{
		UserId:       "user1",
		IssuedAtTime: now,
	}
	session, _, _ := sessionHandler.GetSession(payload)
	sessionHandler.RefreshLastUpdate(session)
	if session.LastUpdate == now {
		t.Error("LastUpdate was not updated correctly")
	}
}

func TestAllSessions(t *testing.T) {
	sessionHandler := createTestSessionHandler(t, time.Now())

	sessions := sessionHandler.GetAllSessions()
	if !reflect.DeepEqual(sessions, sessionHandler.sessions) {
		t.Errorf("expected %v to be %v", sessions, sessionHandler.sessions)
	}
}
