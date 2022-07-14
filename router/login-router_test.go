package router

import (
	"authGo/session"
	"authGo/token"
	"authGo/user"
	"authGo/validator"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func createValidator() *LoginRouter {
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

	return &LoginRouter{
		Services: services,
	}
}

func TestHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/auth/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("admin", "admin")

	loginRouter := createValidator()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(loginRouter.Handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	cookies := rr.Header().Values("Set-Cookie")
	var hasAccessCookie, hasRefreshCookie bool
	for _, cookie := range cookies {
		if strings.Contains(cookie, "accessToken") {
			hasAccessCookie = true
		}
		if strings.Contains(cookie, "refreshToken") {
			hasRefreshCookie = true
		}
	}
	if !hasAccessCookie {
		t.Errorf("accessToken cookie not found, got %s", cookies)
	}
	if !hasRefreshCookie {
		t.Errorf("refreshToken cookie not found, got %s", cookies)
	}

	sessions := loginRouter.Services.SessionsHandler.GetAllSessions()
	if len(sessions) != 1 {
		t.Errorf("expected sessions len to be 1, got %v", sessions)
	}

	expected := `userData`
	if want := strings.TrimSpace(rr.Body.String()); !strings.Contains(want, expected) {
		t.Errorf("handler returned unexpected body: got %v should contain %v", want, expected)
	}
}

func TestHandlerEmptyLogin(t *testing.T) {
	req, err := http.NewRequest("POST", "/auth/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	loginRouter := createValidator()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(loginRouter.Handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := `{"error":"Empty name or password"}`
	if want := strings.TrimSpace(rr.Body.String()); want != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", want, expected)
	}
}

func TestHandlerUserNotFound(t *testing.T) {
	req, err := http.NewRequest("POST", "/auth/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("user123", "user123")

	loginRouter := createValidator()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(loginRouter.Handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := `{"error":"User doesn't exist"}`
	if want := strings.TrimSpace(rr.Body.String()); want != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", want, expected)
	}
}

func TestHandlerPasswordNotValid(t *testing.T) {
	req, err := http.NewRequest("POST", "/auth/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("admin", "12345")

	loginRouter := createValidator()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(loginRouter.Handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := `{"error":"Password not valid"}`
	if want := strings.TrimSpace(rr.Body.String()); want != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", want, expected)
	}
}
