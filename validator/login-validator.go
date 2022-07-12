package validator

import (
	"authGo/session"
	"authGo/token"
	"authGo/user"
	"errors"
	"fmt"
	"time"
)

type LoginValidator struct {
	Validator Validator
}

type LoginDetails struct {
	Name     string
	Password string
}

type JwtTokens struct {
	AccessToken    string
	RefreshToken   string
	AccessPayload  token.AccessTokenPayload
	RefreshPayload token.RefreshTokenPayload
}

var (
	ErrLoginRouterReadingFormData      = errors.New("login validator: error reading form data")
	ErrLoginRouterEmptyNamePassword    = errors.New("login validator: empty name or password")
	ErrLoginRouterUserNotFound         = errors.New("login validator: user not found")
	ErrLoginRouterPasswordNotValid     = errors.New("login validator: password not valid")
	ErrLoginRouterCreatingAccessToken  = errors.New("login validator: error creating accessToken")
	ErrLoginRouterCreatingRefreshToken = errors.New("login validator: error creating accessToken")
)

func (v *LoginValidator) GetLoginDetails() (*LoginDetails, error) {
	err := v.Validator.Request.ParseForm()
	if err != nil {
		return nil, fmt.Errorf("%w, %s", ErrLoginRouterReadingFormData, err)
	}
	var name, password string
	for key, value := range v.Validator.Request.Form {
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

func (v *LoginValidator) GetUser(loginDetails *LoginDetails) (*user.User, error) {
	u, err := v.Validator.Services.UserService.GetRepository().GetByName(loginDetails.Name)
	if err == user.ErrUserNotFound {
		return nil, ErrLoginRouterUserNotFound
	}

	if isPasswordValid := v.Validator.Services.UserService.IsPasswordValid(loginDetails.Name, loginDetails.Password); !isPasswordValid {
		return nil, ErrLoginRouterPasswordNotValid
	}
	return u, nil
}

func (v *LoginValidator) CreateTokens(user *user.User) (*JwtTokens, error) {
	now := time.Now()
	accessPayload := &token.AccessTokenPayload{UserId: user.Id, IssuedAtTime: now, IsAdmin: user.IsAdmin}
	accessJWT, err := v.Validator.Services.AccessTokenGenerator.CreateToken(accessPayload)
	if err != nil {
		return nil, ErrLoginRouterCreatingAccessToken
	}
	refreshPayload := &token.RefreshTokenPayload{UserId: user.Id, IssuedAtTime: now}
	refreshJWT, err := v.Validator.Services.RefreshTokenGenerator.CreateToken(refreshPayload)
	if err != nil {
		return nil, ErrLoginRouterCreatingRefreshToken
	}
	return &JwtTokens{AccessToken: accessJWT, RefreshToken: refreshJWT, AccessPayload: *accessPayload, RefreshPayload: *refreshPayload}, nil
}

func (v *LoginValidator) GetDeviceData() session.DeviceData {
	return session.DeviceData{IpAddress: v.Validator.Request.RemoteAddr, UserAgent: v.Validator.Request.UserAgent()}
}
