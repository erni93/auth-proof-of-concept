package validator

import (
	"authGo/session"
	"authGo/token"
	"authGo/user"
	"errors"
	"fmt"
	"time"
)

type RefreshValidator struct {
	Validator Validator
}

type AccessJwtToken struct {
	AccessToken   string
	AccessPayload token.AccessTokenPayload
}

var (
	ErrRefreshCreatingAccessToken  = errors.New("refresh validator: error creating accessToken")
	ErrRefreshReadingRefreshCookie = errors.New("refresh validator: error reading refreshToken cookie")
)

func (v *RefreshValidator) ValidateRefreshToken() (*session.Session, error) {
	refreshCookie, err := v.Validator.Request.Cookie("refreshToken")
	if err != nil {
		return nil, fmt.Errorf("%w, %s", ErrRefreshReadingRefreshCookie, err)
	}
	err = v.Validator.Services.RefreshTokenGenerator.IsTokenValid(refreshCookie.Value)
	if err != nil {
		return nil, err
	}
	payload := token.RefreshTokenPayload{}
	err = v.Validator.Services.RefreshTokenGenerator.LoadPayload(refreshCookie.Value, &payload)
	if err != nil {
		return nil, err
	}
	session, _, err := v.Validator.Services.SessionsHandler.GetSession(payload)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (v *RefreshValidator) CreateAccessToken(user *user.User) (*AccessJwtToken, error) {
	now := time.Now()
	accessPayload := &token.AccessTokenPayload{UserId: user.Id, IssuedAtTime: now, IsAdmin: user.IsAdmin}
	accessJWT, err := v.Validator.Services.AccessTokenGenerator.CreateToken(accessPayload)
	if err != nil {
		return nil, ErrRefreshCreatingAccessToken
	}
	return &AccessJwtToken{AccessToken: accessJWT, AccessPayload: *accessPayload}, nil
}
