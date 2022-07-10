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

type RefreshRouterValidator struct {
	Writer   http.ResponseWriter
	Request  *http.Request
	Services *RefreshRouterServices
}

type AccessJwtToken struct {
	AccessToken   string
	AccessPayload token.AccessTokenPayload
}

type RefreshRouterServices struct {
	UserService           *user.UserService
	AccessTokenGenerator  *token.TokenGenerator[token.AccessTokenPayload]
	RefreshTokenGenerator *token.TokenGenerator[token.RefreshTokenPayload]
	SessionsHandler       *session.SessionsHandler
}

var (
	ErrRefreshRouterCreatingAccessToken  = errors.New("refresh router: error creating accessToken")
	ErrRefreshRouterReadingRefreshCookie = errors.New("refresh router: error reading refreshToken cookie")
)

func (v *RefreshRouterValidator) ValidateRefreshToken() (*session.Session, error) {
	refreshCookie, err := v.Request.Cookie("refreshToken")
	if err != nil {
		return nil, fmt.Errorf("%w, %s", ErrRefreshRouterReadingRefreshCookie, err)
	}
	err = v.Services.RefreshTokenGenerator.IsTokenValid(refreshCookie.Value)
	if err != nil {
		return nil, err
	}
	payload := token.RefreshTokenPayload{}
	err = v.Services.RefreshTokenGenerator.LoadPayload(refreshCookie.Value, &payload)
	if err != nil {
		return nil, err
	}
	session, _, err := v.Services.SessionsHandler.GetSession(payload)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (v *RefreshRouterValidator) CreateAccessToken(user *user.User) (*AccessJwtToken, error) {
	now := time.Now()
	accessPayload := &token.AccessTokenPayload{UserId: user.Id, IssuedAtTime: now, IsAdmin: user.IsAdmin}
	accessJWT, err := v.Services.AccessTokenGenerator.CreateToken(accessPayload)
	if err != nil {
		return nil, ErrRefreshRouterCreatingAccessToken
	}
	return &AccessJwtToken{AccessToken: accessJWT, AccessPayload: *accessPayload}, nil
}
