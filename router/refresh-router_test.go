package router

import (
	"authGo/session"
	"authGo/token"
	"authGo/user"
	"authGo/validator"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func createRefreshRouter() *RefreshRouter {
	userService := user.NewUserService()
	userService.CreateUser("admin", "admin", true)
	accessTokenGenerator := &token.TokenGenerator[token.AccessTokenPayload]{Password: []byte("accessKey"), Duration: time.Minute * 2}
	refreshTokenGenerator := &token.TokenGenerator[token.RefreshTokenPayload]{Password: []byte("refreshKey"), Duration: time.Hour * 24 * 365}
	sessionHandler := session.NewSessionHandler()

	services := &validator.Services{
		UserService:           userService,
		AccessTokenGenerator:  accessTokenGenerator,
		RefreshTokenGenerator: refreshTokenGenerator,
		SessionsHandler:       sessionHandler,
	}

	return &RefreshRouter{
		Services: services,
	}
}

func TestRefreshRouterHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/auth/refresh", nil)
	if err != nil {
		t.Fatal(err)
	}

	refreshRouter := createRefreshRouter()
	payload := &token.RefreshTokenPayload{UserId: refreshRouter.Services.UserService.GetRepository().GetAll()[0].Id, IssuedAtTime: time.Now()}
	err = refreshRouter.Services.SessionsHandler.AddNewSession(*payload, session.DeviceData{IpAddress: "10.0.0.1", UserAgent: "vscode"})
	if err != nil {
		t.Fatal(err)
	}
	session := refreshRouter.Services.SessionsHandler.GetAllSessions()[0]
	lastSessionUpdate := session.LastUpdate
	time.Sleep(1 * time.Millisecond)
	refreshToken, err := refreshRouter.Services.RefreshTokenGenerator.CreateToken(&session.UserToken)
	if err != nil {
		t.Fatal(err)
	}

	refreshCookie := &http.Cookie{Name: "refreshToken", Value: refreshToken, HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("refreshToken=%s", refreshCookie.Value))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(refreshRouter.Handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	cookies := rr.Header().Values("Set-Cookie")
	var hasAccessCookie bool
	for _, cookie := range cookies {
		if strings.Contains(cookie, "accessToken") {
			hasAccessCookie = true
		}
	}
	if !hasAccessCookie {
		t.Errorf("accessToken cookie not found, got %s", cookies)
	}

	expected := `userData`
	if want := strings.TrimSpace(rr.Body.String()); !strings.Contains(want, expected) {
		t.Errorf("handler returned unexpected body: got %v should contain %v", want, expected)
	}

	if session.LastUpdate == lastSessionUpdate {
		t.Error("session LastUpdate was not updated")
	}
}

func TestRefreshRouterSessionNotFound(t *testing.T) {
	req, err := http.NewRequest("POST", "/auth/refresh", nil)
	if err != nil {
		t.Fatal(err)
	}

	refreshRouter := createRefreshRouter()
	refreshToken, err := refreshRouter.Services.RefreshTokenGenerator.CreateToken(&token.RefreshTokenPayload{UserId: "12345", IssuedAtTime: time.Now()})
	if err != nil {
		t.Fatal(err)
	}

	refreshCookie := &http.Cookie{Name: "refreshToken", Value: refreshToken, HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("refreshToken=%s", refreshCookie.Value))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(refreshRouter.Handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := `{"error":"User session not found"}`
	if want := strings.TrimSpace(rr.Body.String()); want != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", want, expected)
	}
}

func TestRefreshRouterUserNotValid(t *testing.T) {
	req, err := http.NewRequest("POST", "/auth/refresh", nil)
	if err != nil {
		t.Fatal(err)
	}

	refreshRouter := createRefreshRouter()
	payload := &token.RefreshTokenPayload{UserId: "12345", IssuedAtTime: time.Now()}
	err = refreshRouter.Services.SessionsHandler.AddNewSession(*payload, session.DeviceData{IpAddress: "10.0.0.1", UserAgent: "vscode"})
	if err != nil {
		t.Fatal(err)
	}
	session := refreshRouter.Services.SessionsHandler.GetAllSessions()[0]
	refreshToken, err := refreshRouter.Services.RefreshTokenGenerator.CreateToken(&session.UserToken)
	if err != nil {
		t.Fatal(err)
	}

	refreshCookie := &http.Cookie{Name: "refreshToken", Value: refreshToken, HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("refreshToken=%s", refreshCookie.Value))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(refreshRouter.Handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := `{"error":"User not valid"}`
	if want := strings.TrimSpace(rr.Body.String()); want != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", want, expected)
	}
}
