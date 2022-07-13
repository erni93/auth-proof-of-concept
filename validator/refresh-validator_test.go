package validator

import (
	"authGo/session"
	"authGo/token"
	"authGo/user"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestValidateRefreshToken(t *testing.T) {
	req, err := http.NewRequest("POST", "/auth/refresh", nil)
	if err != nil {
		t.Fatal(err)
	}

	refreshTokenGenerator := &token.TokenGenerator[token.RefreshTokenPayload]{Password: []byte("refreshKey"), Duration: time.Hour * 24 * 365}
	payload := &token.RefreshTokenPayload{UserId: "1", IssuedAtTime: time.Now()}
	refreshToken, err := refreshTokenGenerator.CreateToken(payload)
	if err != nil {
		t.Fatal(err)
	}

	refreshCookie := &http.Cookie{Name: "refreshToken", Value: refreshToken, HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("refreshToken=%s", refreshCookie.Value))

	sessionHandler := session.NewSessionHandler()
	sessionHandler.AddNewSession(*payload, session.DeviceData{IpAddress: "10.0.0.1", UserAgent: "vscode"})
	v := RefreshValidator{Validator: Validator{Request: req, Services: &Services{RefreshTokenGenerator: refreshTokenGenerator, SessionsHandler: sessionHandler}}}
	session, err := v.ValidateRefreshToken()
	if err != nil {
		t.Errorf("expected err to be nil, got %s", err)
	}
	if session == nil {
		t.Error("expected session to not be nil")
	}
}

func TestErrorsValidateRefreshToken(t *testing.T) {
	req, err := http.NewRequest("POST", "/auth/refresh", nil)
	if err != nil {
		t.Fatal(err)
	}

	refreshTokenGenerator1 := &token.TokenGenerator[token.RefreshTokenPayload]{Password: []byte("refreshKey1"), Duration: time.Hour * 24 * 365}
	refreshTokenGenerator2 := &token.TokenGenerator[token.RefreshTokenPayload]{Password: []byte("refreshKey2"), Duration: time.Hour * 24 * 365}
	payload := &token.RefreshTokenPayload{UserId: "1", IssuedAtTime: time.Now()}
	sessionHandler := session.NewSessionHandler()
	sessionHandler.AddNewSession(*payload, session.DeviceData{IpAddress: "10.0.0.1", UserAgent: "vscode"})
	v := RefreshValidator{Validator: Validator{Request: req, Services: &Services{RefreshTokenGenerator: refreshTokenGenerator1, SessionsHandler: sessionHandler}}}

	_, err = v.ValidateRefreshToken()
	if !errors.Is(err, ErrRefreshReadingRefreshCookie) {
		t.Errorf("expected err to be part of ErrRefreshReadingRefreshCookie, got %s", err)
	}

	invalidRefreshToken, err := refreshTokenGenerator2.CreateToken(payload)
	if err != nil {
		t.Fatal(err)
	}
	invalidRefreshCookie := &http.Cookie{Name: "refreshToken", Value: invalidRefreshToken, HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("refreshToken=%s", invalidRefreshCookie.Value))

	_, err = v.ValidateRefreshToken()
	if !errors.Is(err, token.ErrInvalidSignature) {
		t.Errorf("expected err to be part of ErrInvalidSignature, got %s", err)
	}

	invalidPayload := &token.RefreshTokenPayload{UserId: "12345", IssuedAtTime: time.Now()}
	refreshToken, err := refreshTokenGenerator1.CreateToken(invalidPayload)
	if err != nil {
		t.Fatal(err)
	}
	refreshCookie := &http.Cookie{Name: "refreshToken", Value: refreshToken, HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("refreshToken=%s", refreshCookie.Value))

	_, err = v.ValidateRefreshToken()
	if !errors.Is(err, session.ErrUserTokenNotFound) {
		t.Errorf("expected err to be part of ErrUserTokenNotFound, got %s", err)
	}
}

func TestCreateAccessToken(t *testing.T) {
	accessTokenGenerator := &token.TokenGenerator[token.AccessTokenPayload]{Password: []byte("accessKey"), Duration: time.Minute * 2}
	user := &user.User{Id: "1", Name: "user1", Password: "user1", IsAdmin: true}

	v := RefreshValidator{Validator: Validator{Services: &Services{AccessTokenGenerator: accessTokenGenerator}}}

	accessToken, err := v.CreateAccessToken(user)
	if err != nil {
		t.Errorf("expected err to be nil, got %s", err)
	}
	if accessToken == nil {
		t.Error("expected accessToken to not be nil")
	}
}
