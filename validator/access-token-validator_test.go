package validator

import (
	"authGo/token"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestValidateAccessToken(t *testing.T) {
	req, err := http.NewRequest("POST", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	accessTokenGenerator := &token.TokenGenerator[token.AccessTokenPayload]{Password: []byte("accessKey"), Duration: time.Minute * 2}
	payload := &token.AccessTokenPayload{UserId: "1", IssuedAtTime: time.Now(), IsAdmin: true}
	accessToken, err := accessTokenGenerator.CreateToken(payload)
	if err != nil {
		t.Fatal(err)
	}

	accessCookie := &http.Cookie{Name: "accessToken", Value: accessToken, HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("accessToken=%s", accessCookie.Value))

	v := AccessTokenValidator{Validator: Validator{Request: req, Services: &Services{AccessTokenGenerator: accessTokenGenerator}}}

	requestPayload, err := v.ValidateAccessToken()
	if err != nil {
		t.Errorf("expected err to be nil, got %s", err)
	}
	if requestPayload.UserId != payload.UserId || requestPayload.IsAdmin != payload.IsAdmin || !requestPayload.IssuedAtTime.Equal(payload.IssuedAtTime) {
		t.Errorf("expected %v to be %v", requestPayload, payload)
	}
}

func TestErrorsValidateAccessToken(t *testing.T) {
	req, err := http.NewRequest("POST", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	accessTokenGenerator1 := &token.TokenGenerator[token.AccessTokenPayload]{Password: []byte("accessKey1"), Duration: time.Minute * 2}
	accessTokenGenerator2 := &token.TokenGenerator[token.AccessTokenPayload]{Password: []byte("accessKey2"), Duration: time.Minute * 2}
	payload := &token.AccessTokenPayload{UserId: "1", IssuedAtTime: time.Now(), IsAdmin: true}

	v := AccessTokenValidator{Validator: Validator{Request: req, Services: &Services{AccessTokenGenerator: accessTokenGenerator1}}}

	_, err = v.ValidateAccessToken()
	if !errors.Is(err, ErrAccessReadingAccessCookie) {
		t.Errorf("expected err to be part of ErrAccessReadingAccessCookie, got %s", err)
	}

	invalidAccessToken, err := accessTokenGenerator2.CreateToken(payload)
	if err != nil {
		t.Fatal(err)
	}
	invalidAccessCookie := &http.Cookie{Name: "accessToken", Value: invalidAccessToken, HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("accessToken=%s", invalidAccessCookie.Value))

	_, err = v.ValidateAccessToken()
	if !errors.Is(err, token.ErrInvalidSignature) {
		t.Errorf("expected err to be part of ErrInvalidSignature, got %s", err)
	}
}
