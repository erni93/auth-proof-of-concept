package validator

import (
	"authGo/session"
	"authGo/token"
	"authGo/user"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type LoginRouterValidator struct {
	Writer   http.ResponseWriter
	Request  *http.Request
	Services *LoginRouterServices
}

type LoginDetails struct {
	Name     string
	Password string
}

type JwtTokens struct {
	AccessToken    string
	RefreshToken   string
	RefreshPayload token.RefreshTokenPayload
}

type LoginRouterServices struct {
	UserService           *user.UserService
	AccessTokenGenerator  *token.TokenGenerator[token.AccessTokenPayload]
	RefreshTokenGenerator *token.TokenGenerator[token.RefreshTokenPayload]
}

var (
	ErrLoginRouterReadingFormData      = errors.New("login router: error reading form data")
	ErrLoginRouterEmptyNamePassword    = errors.New("login router: empty name or password")
	ErrLoginRouterUserNotFound         = errors.New("login router: user not found")
	ErrLoginRouterPasswordNotValid     = errors.New("login router: password not valid")
	ErrLoginRouterCreatingAccessToken  = errors.New("login router: error creating accessToken")
	ErrLoginRouterCreatingRefreshToken = errors.New("login router: error creating accessToken")
)

func (v *LoginRouterValidator) GetLoginDetails() (*LoginDetails, error) {
	err := v.Request.ParseForm()
	if err != nil {
		return nil, fmt.Errorf("%w, %s", ErrLoginRouterReadingFormData, err)
	}
	var name, password string
	for key, value := range v.Request.Form {
		switch key {
		case "name":
			name = value[0]
		case "password":
			password = value[0]
		}
	}
	if name == "" || password == "" {
		return nil, ErrLoginRouterEmptyNamePassword
	}
	return &LoginDetails{Name: name, Password: password}, nil
}

func (v *LoginRouterValidator) GetUser(loginDetails *LoginDetails) (*user.User, error) {
	u, err := v.Services.UserService.GetRepository().GetByName(loginDetails.Name)
	if err == user.ErrUserNotFound {
		return nil, ErrLoginRouterUserNotFound
	}

	if isPasswordValid := v.Services.UserService.IsPasswordValid(loginDetails.Name, loginDetails.Password); !isPasswordValid {
		return nil, ErrLoginRouterPasswordNotValid
	}
	return u, nil
}

func (v *LoginRouterValidator) GetTokens(user *user.User) (*JwtTokens, error) {
	now := time.Now()
	accessPayload := &token.AccessTokenPayload{UserId: user.Id, IssuedAtTime: now, IsAdmin: user.IsAdmin}
	accessJWT, err := v.Services.AccessTokenGenerator.CreateToken(accessPayload)
	if err != nil {
		return nil, ErrLoginRouterCreatingAccessToken
	}
	refreshPayload := &token.RefreshTokenPayload{UserId: user.Id, IssuedAtTime: now}
	refreshJWT, err := v.Services.RefreshTokenGenerator.CreateToken(refreshPayload)
	if err != nil {
		return nil, ErrLoginRouterCreatingRefreshToken
	}
	return &JwtTokens{AccessToken: accessJWT, RefreshToken: refreshJWT, RefreshPayload: *refreshPayload}, nil
}

func (v *LoginRouterValidator) GetDeviceData() session.DeviceData {
	return session.DeviceData{IpAddress: v.Request.RemoteAddr, UserAgent: v.Request.UserAgent()}
}
