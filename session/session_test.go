package session

import (
	"authGo/authentication"
	"testing"
	"time"
)

func createTestSession(userId string) *Session {
	return &Session{
		UserToken: authentication.TokenData{
			UserId:       userId,
			IssuedAtTime: time.Now(),
		},
		DeviceData: DeviceData{
			IpAddress: "ip-" + userId,
			UserAgent: "userAgent-" + userId,
		},
		Created: time.Now(),
	}
}

func createTestSessionHandler(t *testing.T) *SessionsHandler {
	sessionHandler := NewSessionHandler()
	err := sessionHandler.AddSession(createTestSession("user1"))
	if err != nil {
		t.Errorf("error adding session, %s", err)
	}
	err = sessionHandler.AddSession(createTestSession("user2"))
	if err != nil {
		t.Errorf("error adding session, %s", err)
	}
	return sessionHandler
}

func TestNewSessionHandler(t *testing.T) {
	if sessionHandler := NewSessionHandler(); sessionHandler == nil {
		test := *sessionHandler
		t.Errorf("expected NewSessionHandler to return SessionsHandler, got %v", test)
	}
}

func TestGetUserSessions(t *testing.T) {
	sessionHandler := createTestSessionHandler(t)

	user1Sessions := sessionHandler.GetUserSessions("user1")
	if len(user1Sessions) != 1 {
		t.Errorf("expected user1Sessions len to be 1, got %d", len(user1Sessions))
	}

	time.Sleep(1 * time.Millisecond)
	err := sessionHandler.AddSession(createTestSession("user1"))
	if err != nil {
		t.Errorf("error adding session, %s", err)
	}
	user1Sessions = sessionHandler.GetUserSessions("user1")
	if len(user1Sessions) != 2 {
		t.Errorf("expected user1Sessions len to be 2, got %d", len(user1Sessions))
	}
}

func TestAddSession(t *testing.T) {
	sessionHandler := createTestSessionHandler(t)
	err := sessionHandler.AddSession(createTestSession("user1"))
	if err != ErrSessionAlreadyExists {
		t.Errorf("expected error to be ErrSessionAlreadyExists, got: %s", err)
	}
}

func TestDeleteSession(t *testing.T) {
	sessionHandler := createTestSessionHandler(t)

	time.Sleep(1 * time.Millisecond)
	newSession := createTestSession("user1")

	err := sessionHandler.AddSession(newSession)
	if err != nil {
		t.Errorf("error adding session, %s", err)
	}

	user1Sessions := sessionHandler.GetUserSessions("user1")
	err = sessionHandler.DeleteSession(user1Sessions[0].UserToken)
	if err != nil {
		t.Errorf("error deleting session, %s", err)
	}

	err = sessionHandler.DeleteSession(authentication.TokenData{UserId: "user3", IssuedAtTime: time.Now()})
	if err != ErrUserTokenNotFound {
		t.Errorf("expected error to be ErrUserTokenNotFound, got: %s", err)
	}
}

func TestRenewUserToken(t *testing.T) {
	sessionHandler := createTestSessionHandler(t)
	time.Sleep(1 * time.Millisecond)
	oldData := sessionHandler.GetUserSessions("user1")[0].UserToken
	newData := authentication.TokenData{UserId: "user1", IssuedAtTime: time.Now()}

	err := sessionHandler.RenewUserToken(oldData, newData)
	if err != nil {
		t.Errorf("error renewing session, %s", err)
	}
	if sessionHandler.getSession(newData) == nil {
		t.Errorf("new user token was not found, %s", newData)
	}

	time.Sleep(1 * time.Millisecond)
	err = sessionHandler.RenewUserToken(newData, authentication.TokenData{UserId: "user3", IssuedAtTime: time.Now()})
	if err != ErrRenewUserTokenDifferent {
		t.Errorf("expected error to be ErrRenewUserTokenDifferent, got: %s", err)
	}

	time.Sleep(1 * time.Millisecond)
	err = sessionHandler.RenewUserToken(authentication.TokenData{UserId: "user3", IssuedAtTime: time.Now()}, newData)
	if err != ErrUserTokenNotFound {
		t.Errorf("expected error to be ErrUserTokenNotFound, got: %s", err)
	}
}
