package authentication

import (
	"authGo/user"
	"testing"
	"time"
)

type TestToken struct {
	name string
	jwt  string
	want string
	err  error
}

func TestCreateToken(t *testing.T) {
	tg := NewTokenGenerator("accessKey", "refreshKey")
	issuedAtTime := time.Date(2022, 8, 6, 0, 0, 0, 0, time.UTC)
	user := &user.User{Id: "1", Name: "admin", IsAdmin: true, Password: ""}

	// JWT created from here https://jwt.io/ note, secret is not in base64 format
	testTokens := make([]*TestToken, 0)
	testAccessToken := &TestToken{name: "accessToken", want: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxIiwiaXNzdWVkQXRUaW1lIjoiMjAyMi0wOC0wNlQwMDowMDowMFoiLCJpc0FkbWluIjp0cnVlfQ.aqrs8ystc9s5KUpXeAWQaQCG8YffKsp-o-2cXKy80DE"}
	testAccessToken.jwt, testAccessToken.err = tg.CreateAccessToken(user, issuedAtTime)
	testRefreshToken := &TestToken{name: "refreshToken", want: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxIiwiaXNzdWVkQXRUaW1lIjoiMjAyMi0wOC0wNlQwMDowMDowMFoifQ.s-L_k5JIq-CqBs44avYIu09CPMDczFbrRavktZvX8bU"}
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
