package validator

import (
	"authGo/token"
	"authGo/user"
	"net/http"
	"testing"
	"time"
)

func TestGetLoginDetails(t *testing.T) {
	req, err := http.NewRequest("POST", "/auth/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("user1", "user1")

	v := LoginValidator{Validator: Validator{Request: req}}

	loginDetails, err := v.GetLoginDetails()
	if err != nil {
		t.Errorf("expected err to be nil, got %s", err)
	}
	if loginDetails == nil || loginDetails.Name != "user1" || loginDetails.Password != "user1" {
		t.Errorf("expected loginDetails to has user1:user1, got %s", loginDetails)
	}

	req.SetBasicAuth("", "")
	loginDetails, err = v.GetLoginDetails()
	if err != ErrLoginRouterEmptyNamePassword {
		t.Error("expected err to be ErrLoginRouterEmptyNamePassword")
	}
	if loginDetails != nil {
		t.Errorf("expected loginDetails to be nil, got %s", loginDetails)
	}
}

func TestGetUser(t *testing.T) {
	req, err := http.NewRequest("POST", "/auth/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	userService := user.NewUserService()
	userService.CreateUser("user1", "user1", true)
	v := LoginValidator{Validator: Validator{Request: req, Services: &Services{UserService: userService}}}

	// Should retrieve the user
	req.SetBasicAuth("user1", "user1")
	loginDetails, err := v.GetLoginDetails()
	if err != nil {
		t.Fatal(err)
	}
	user, err := v.GetUser(loginDetails)
	if err != nil {
		t.Errorf("expected err to be nil, got %s", err)
	}
	if user == nil || user.Name != "user1" {
		t.Errorf("expected user to have user1 data, got %v", user)
	}

	// Should return user not found
	req.SetBasicAuth("user2", "user2")
	loginDetails, err = v.GetLoginDetails()
	if err != nil {
		t.Fatal(err)
	}
	user, err = v.GetUser(loginDetails)
	if err != ErrLoginRouterUserNotFound {
		t.Errorf("expected err to be ErrLoginRouterUserNotFound, got %s", err)
	}
	if user != nil {
		t.Errorf("expected user to be nil, got %v", user)
	}

	// Should return invalid password
	req.SetBasicAuth("user1", "user11111")
	loginDetails, err = v.GetLoginDetails()
	if err != nil {
		t.Fatal(err)
	}
	user, err = v.GetUser(loginDetails)
	if err != ErrLoginRouterPasswordNotValid {
		t.Errorf("expected err to be ErrLoginRouterPasswordNotValid, got %s", err)
	}
	if user != nil {
		t.Errorf("expected user to be nil, got %v", user)
	}
}

func TestCreateTokens(t *testing.T) {
	req, err := http.NewRequest("POST", "/auth/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	accessTokenGenerator := &token.TokenGenerator[token.AccessTokenPayload]{Password: []byte("accessKey"), Duration: time.Minute * 2}
	refreshTokenGenerator := &token.TokenGenerator[token.RefreshTokenPayload]{Password: []byte("refreshKey"), Duration: time.Hour * 24 * 365}
	v := LoginValidator{Validator: Validator{Request: req, Services: &Services{
		AccessTokenGenerator: accessTokenGenerator, RefreshTokenGenerator: refreshTokenGenerator}},
	}

	jwtTokens, err := v.CreateTokens(&user.User{Id: "1", Name: "user1", Password: "user1", IsAdmin: true})
	if err != nil {
		t.Errorf("expected err to be nil, got %s", err)
	}
	if jwtTokens == nil {
		t.Error("expected jwtTokens to not be nil")
	}
	if jwtTokens.AccessToken == "" || jwtTokens.RefreshToken == "" {
		t.Errorf("expected AccessToken and RefreshToken to not be empty, got %v", jwtTokens)
	}
	if jwtTokens.AccessPayload.UserId != "1" || jwtTokens.RefreshPayload.UserId != "1" {
		t.Errorf("expected AccessPayload.UserId and RefreshPayload.UserId to be 1, got %s and %s", jwtTokens.AccessPayload.UserId, jwtTokens.RefreshPayload.UserId)
	}
}

func TestGetDeviceData(t *testing.T) {
	req, err := http.NewRequest("POST", "/auth/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("User-Agent", "golang")
	req.RemoteAddr = "127.0.0.1"

	v := LoginValidator{Validator: Validator{Request: req}}

	deviceData := v.GetDeviceData()
	if deviceData.IpAddress != "127.0.0.1" {
		t.Errorf("expected deviceData.IpAddress to be 127.0.0.1, got %s", deviceData.IpAddress)
	}
	if deviceData.UserAgent != "golang" {
		t.Errorf("expected deviceData.UserAgent to be golang, got %s", deviceData.UserAgent)
	}
}
