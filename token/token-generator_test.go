package token

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

type ValidateTokenTest struct {
	jwt string
	err error
}

// JWT created from here https://jwt.io/ password "accessKey", secret is not in base64 format
var AccessTokenJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxIiwiaXNzdWVkQXRUaW1lIjoiMjAyMi0wOC0wNlQwMDowMDowMFoiLCJpc0FkbWluIjp0cnVlfQ.aqrs8ystc9s5KUpXeAWQaQCG8YffKsp-o-2cXKy80DE"

func TestCreateToken(t *testing.T) {
	tg := &TokenGenerator[AccessTokenPayload]{Password: []byte("accessKey")}
	payload := &AccessTokenPayload{UserId: "1", IssuedAtTime: time.Date(2022, 8, 6, 0, 0, 0, 0, time.UTC), IsAdmin: true}
	jwt, err := tg.CreateToken(payload)
	if err != nil {
		t.Errorf("expected err to be nil, got %s", err)
	}
	if jwt != AccessTokenJWT {
		t.Errorf("want %s, got %s", AccessTokenJWT, jwt)
	}
}

func TestIsTokenValid(t *testing.T) {
	tg := &TokenGenerator[RefreshTokenPayload]{Password: []byte("refreshKey"), Duration: time.Hour * 24 * 365}
	validDuration := time.Now().Add(-time.Hour * 24 * 200)
	expiredDuration := time.Now().Add(-time.Hour * 24 * 500)

	validToken, err := tg.CreateToken(&RefreshTokenPayload{UserId: "1", IssuedAtTime: validDuration})
	if err != nil {
		t.Errorf("expected err to be nil, got %s", err)
	}
	expiredToken, err := tg.CreateToken(&RefreshTokenPayload{UserId: "1", IssuedAtTime: expiredDuration})
	if err != nil {
		t.Errorf("expected err to be nil, got %s", err)
	}
	tg2 := &TokenGenerator[RefreshTokenPayload]{Password: []byte("refreshKey2"), Duration: time.Minute * 2}
	invalidSignatureJWT, err := tg2.CreateToken(&RefreshTokenPayload{UserId: "1", IssuedAtTime: validDuration})
	if err != nil {
		t.Errorf("expected err to be nil, got %s", err)
	}

	testTokens := make([]*ValidateTokenTest, 0)
	testTokens = append(testTokens,
		&ValidateTokenTest{jwt: validToken, err: nil},
		&ValidateTokenTest{jwt: "aa.bb", err: ErrInvalidJWTLength},
		&ValidateTokenTest{jwt: expiredToken, err: ErrTokenExpired},
		&ValidateTokenTest{jwt: invalidSignatureJWT, err: ErrInvalidSignature},
	)

	for _, token := range testTokens {
		err := tg.IsTokenValid(token.jwt)
		if !(err == nil && token.err == nil) && errors.Is(token.err, err) {
			t.Errorf("expected error to be %s, got %s", token.err, err)
		}
	}

}

func TestLoadPayload(t *testing.T) {
	tg := &TokenGenerator[AccessTokenPayload]{Password: []byte("accessKey")}
	payload := &AccessTokenPayload{UserId: "1", IssuedAtTime: time.Date(2022, 8, 6, 0, 0, 0, 0, time.UTC), IsAdmin: true}
	jwt, err := tg.CreateToken(payload)
	if err != nil {
		t.Errorf("expected err to be nil, got %s", err)
	}
	if jwt != AccessTokenJWT {
		t.Errorf("want %s, got %s", AccessTokenJWT, jwt)
	}
	samePayload := &AccessTokenPayload{}
	err = tg.LoadPayload(jwt, samePayload)
	if err != nil {
		t.Errorf("expected err to be nil, got %s", err)
	}
	if !reflect.DeepEqual(samePayload, payload) {
		t.Errorf("wanted %v to be %v", samePayload, payload)
	}
}
