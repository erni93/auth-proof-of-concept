package router

import (
	"authGo/session"
	"authGo/token"
	"authGo/user"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type LoginRouter struct {
	Services       *LoginRouterServices
	SessionHandler *session.SessionsHandler
}

func (l *LoginRouter) Handler(w http.ResponseWriter, r *http.Request) {
	validator := LoginRouterValidator{writer: w, request: r, services: l.Services}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	loginDetails, err := validator.GetLoginDetails()
	if err != nil {
		log.Print(err)
		if errors.Is(err, ErrLoginRouterReadingFormData) {
			writeGeneralError(w)
		} else if errors.Is(err, ErrLoginRouterEmptyNamePassword) {
			writeError(w, "Empty name or password")
		}
		return
	}

	user, err := validator.GetUser(loginDetails)
	if err != nil {
		log.Print(err)
		if errors.Is(err, ErrLoginRouterUserNotFound) {
			writeError(w, "User doesn't exist")
		} else if errors.Is(err, ErrLoginRouterPasswordNotValid) {
			writeError(w, "Password not valid")
		}
		return
	}

	tokens, err := validator.GetTokens(user)
	if err != nil {
		log.Print(err)
		writeGeneralError(w)
		return
	}

	newSession := &session.Session{UserToken: tokens.refreshPayload, DeviceData: validator.GetDeviceData(), Created: time.Now()}
	l.SessionHandler.AddSession(newSession)

	writeSuccessLogin(w, tokens)
}

type LoginRouterValidator struct {
	writer   http.ResponseWriter
	request  *http.Request
	services *LoginRouterServices
}

type LoginDetails struct {
	Name     string
	Password string
}

type JwtTokens struct {
	accessToken    string
	refreshToken   string
	refreshPayload token.RefreshTokenPayload
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
	err := v.request.ParseForm()
	if err != nil {
		return nil, fmt.Errorf("%w, %s", ErrLoginRouterReadingFormData, err)
	}
	var name, password string
	for key, value := range v.request.Form {
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
	u, err := v.services.UserService.GetRepository().GetByName(loginDetails.Name)
	if err == user.ErrUserNotFound {
		return nil, ErrLoginRouterUserNotFound
	}

	if isPasswordValid := v.services.UserService.IsPasswordValid(loginDetails.Name, loginDetails.Password); !isPasswordValid {
		return nil, ErrLoginRouterPasswordNotValid
	}
	return u, nil
}

func (v *LoginRouterValidator) GetTokens(user *user.User) (*JwtTokens, error) {
	now := time.Now()
	accessPayload := &token.AccessTokenPayload{UserId: user.Id, IssuedAtTime: now, IsAdmin: user.IsAdmin}
	accessJWT, err := v.services.AccessTokenGenerator.CreateToken(accessPayload)
	if err != nil {
		return nil, ErrLoginRouterCreatingAccessToken
	}
	refreshPayload := &token.RefreshTokenPayload{UserId: user.Id, IssuedAtTime: now}
	refreshJWT, err := v.services.RefreshTokenGenerator.CreateToken(refreshPayload)
	if err != nil {
		return nil, ErrLoginRouterCreatingRefreshToken
	}
	return &JwtTokens{accessToken: accessJWT, refreshToken: refreshJWT, refreshPayload: *refreshPayload}, nil
}

func (v *LoginRouterValidator) GetDeviceData() session.DeviceData {
	return session.DeviceData{IpAddress: v.request.RemoteAddr, UserAgent: v.request.UserAgent()}
}

func writeSuccessLogin(w http.ResponseWriter, tokens *JwtTokens) {
	accessCookie := &http.Cookie{Name: "accessToken", Value: tokens.accessToken, HttpOnly: true, Path: "/"}
	refreshCookie := &http.Cookie{Name: "refreshToken", Value: tokens.refreshToken, HttpOnly: true, Path: "/"}
	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)
	w.WriteHeader(http.StatusOK)
}
