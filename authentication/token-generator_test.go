package authentication

import (
	"authGo/user"
	"errors"
	"testing"
	"time"
)

type CreateTokenTest struct {
	name string
	jwt  string
	want string
	err  error
}

type ValidateTokenTest struct {
	jwt     string
	isValid bool
	err     error
}

var (
	AccessTokenJWT  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxIiwiaXNzdWVkQXRUaW1lIjoiMjAyMi0wOC0wNlQwMDowMDowMFoiLCJpc0FkbWluIjp0cnVlfQ.aqrs8ystc9s5KUpXeAWQaQCG8YffKsp-o-2cXKy80DE"
	RefreshTokenJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxIiwiaXNzdWVkQXRUaW1lIjoiMjAyMi0wOC0wNlQwMDowMDowMFoifQ.s-L_k5JIq-CqBs44avYIu09CPMDczFbrRavktZvX8bU"
)

func TestCreateToken(t *testing.T) {
	tg := &TokenGenerator{accessTokenPassword: []byte("accessKey"), refreshTokenPassword: []byte("refreshKey")}
	issuedAtTime := time.Date(2022, 8, 6, 0, 0, 0, 0, time.UTC)
	user := &user.User{Id: "1", Name: "admin", IsAdmin: true, Password: ""}

	// JWT created from here https://jwt.io/ note, secret is not in base64 format
	testTokens := make([]*CreateTokenTest, 0)
	testAccessToken := &CreateTokenTest{name: "accessToken", want: AccessTokenJWT}
	testAccessToken.jwt, testAccessToken.err = tg.CreateAccessToken(user, issuedAtTime)
	testRefreshToken := &CreateTokenTest{name: "refreshToken", want: RefreshTokenJWT}
	testRefreshToken.jwt, testRefreshToken.err = tg.CreateRefreshToken(user, issuedAtTime)
	testTokens = append(testTokens, testAccessToken, testRefreshToken)

	for _, testToken := range testTokens {
		t.Run(testToken.name, func(t *testing.T) {
			if testToken.err != nil {
				t.Errorf("expected err to be nil, got %s", testToken.err)
			}
			if testToken.jwt != testToken.want {
				t.Errorf("want %s, got %s", testToken.want, testToken.jwt)
			}
		})
	}
}

func TestIsAccessTokenValid(t *testing.T) {
	tg := &TokenGenerator{
		accessTokenPassword: []byte("accessKey"), refreshTokenPassword: []byte("refreshKey"),
		accessTokenDuration: time.Minute * 2, refreshTokenDuration: time.Hour * 24 * 365,
	}
	user := &user.User{Id: "1", Name: "admin", IsAdmin: true, Password: ""}

	validAccessIssuedAtTime := time.Now().Add(-time.Minute * 1)
	expiredAccessIssuedAtTime := time.Now().Add(-time.Minute * 3)

	accessTokens := make([]*ValidateTokenTest, 0)
	validAccessJWT, err := tg.CreateAccessToken(user, validAccessIssuedAtTime)
	if err != nil {
		t.Errorf("expected err to be nil, got %s", err)
	}
	expiredAccessJWT, err := tg.CreateAccessToken(user, expiredAccessIssuedAtTime)
	if err != nil {
		t.Errorf("expected err to be nil, got %s", err)
	}

	tg2 := &TokenGenerator{
		accessTokenPassword: []byte("accessKey2"), refreshTokenPassword: []byte("refreshKey2"),
	}
	invalidSignatureJWT, err := tg2.CreateAccessToken(user, validAccessIssuedAtTime)
	if err != nil {
		t.Errorf("expected err to be nil, got %s", err)
	}

	accessTokens = append(accessTokens,
		&ValidateTokenTest{jwt: validAccessJWT, isValid: true, err: nil},
		&ValidateTokenTest{jwt: "aa.bb", isValid: false, err: ErrInvalidJWTLength},
		&ValidateTokenTest{jwt: expiredAccessJWT, isValid: false, err: ErrTokenExpired},
		&ValidateTokenTest{jwt: invalidSignatureJWT, isValid: false, err: ErrInvalidSignature},
	)

	for _, token := range accessTokens {
		isValid, err := tg.IsAccessTokenValid(token.jwt)
		if isValid != token.isValid {
			t.Errorf("expected isValid to be %t, err %s", token.isValid, err)
		}
		if !(err == nil && token.err == nil) && errors.Is(token.err, err) {
			t.Errorf("expected error to be %s, got %s", token.err, err)
		}
	}

}

func TestIsRefreshTokenValid(t *testing.T) {
	tg := &TokenGenerator{
		accessTokenPassword: []byte("accessKey"), refreshTokenPassword: []byte("refreshKey"),
		accessTokenDuration: time.Minute * 2, refreshTokenDuration: time.Hour * 24 * 365,
	}
	user := &user.User{Id: "1", Name: "admin", IsAdmin: true, Password: ""}

	validRefreshIssuedAtTime := time.Now().Add(-time.Hour * 24 * 200)
	expiredRefreshIssuedAtTime := time.Now().Add(-time.Hour * 24 * 500)

	refreshTokens := make([]*ValidateTokenTest, 0)
	validRefreshJWT, err := tg.CreateRefreshToken(user, validRefreshIssuedAtTime)
	if err != nil {
		t.Errorf("expected err to be nil, got %s", err)
	}
	expiredRefreshJWT, err := tg.CreateRefreshToken(user, expiredRefreshIssuedAtTime)
	if err != nil {
		t.Errorf("expected err to be nil, got %s", err)
	}

	tg2 := &TokenGenerator{
		accessTokenPassword: []byte("accessKey2"), refreshTokenPassword: []byte("refreshKey2"),
	}
	invalidSignatureJWT, err := tg2.CreateRefreshToken(user, validRefreshIssuedAtTime)
	if err != nil {
		t.Errorf("expected err to be nil, got %s", err)
	}

	refreshTokens = append(refreshTokens,
		&ValidateTokenTest{jwt: validRefreshJWT, isValid: true, err: nil},
		&ValidateTokenTest{jwt: "aa.bb", isValid: false, err: ErrInvalidJWTLength},
		&ValidateTokenTest{jwt: expiredRefreshJWT, isValid: false, err: ErrTokenExpired},
		&ValidateTokenTest{jwt: invalidSignatureJWT, isValid: false, err: ErrInvalidSignature},
	)

	for _, token := range refreshTokens {
		isValid, err := tg.IsRefreshTokenValid(token.jwt)
		if isValid != token.isValid {
			t.Errorf("expected isValid to be %t, err %s", token.isValid, err)
		}
		if !(err == nil && token.err == nil) && errors.Is(token.err, err) {
			t.Errorf("expected error to be %s, got %s", token.err, err)
		}
	}

}
