package router

import (
	response "authGo/router/response"
	"authGo/session"
	"authGo/token"
	"authGo/user"
	"authGo/validator"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func createSessionRouter() *SessionRouter {
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

	return &SessionRouter{
		Services: services,
	}
}

func addUserAndSession(t *testing.T, services validator.Services, name string, password string, isAdmin bool) *user.User {
	err := services.UserService.CreateUser(name, password, isAdmin)
	if err != nil {
		t.Fatal(err)
	}
	user, err := services.UserService.GetRepository().GetByName(name)
	if err != nil {
		t.Fatal(err)
	}
	addSession(t, services, user)
	return user
}

func addSession(t *testing.T, services validator.Services, user *user.User) {
	refreshPayload := &token.RefreshTokenPayload{UserId: user.Id, IssuedAtTime: time.Now()}
	err := services.SessionsHandler.AddNewSession(*refreshPayload, session.DeviceData{IpAddress: "10.0.0.1", UserAgent: "vscode"})
	if err != nil {
		t.Fatal(err)
	}
}
func TestSessionRouterHandlerAdmin(t *testing.T) {
	req, err := http.NewRequest("GET", "/sessions", nil)
	if err != nil {
		t.Fatal(err)
	}

	sessionRouter := createSessionRouter()
	user, err := sessionRouter.Services.UserService.GetRepository().GetByName("admin")
	if err != nil {
		t.Fatal(err)
	}

	accessPayload := &token.AccessTokenPayload{UserId: user.Id, IssuedAtTime: time.Now(), IsAdmin: user.IsAdmin}
	addSession(t, *sessionRouter.Services, user)
	accessToken, err := sessionRouter.Services.AccessTokenGenerator.CreateToken(accessPayload)
	if err != nil {
		t.Fatal(err)
	}
	addUserAndSession(t, *sessionRouter.Services, "user2", "user2", false)

	accessCookie := &http.Cookie{Name: "accessToken", Value: accessToken, HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("accessToken=%s", accessCookie.Value))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sessionRouter.GetSessionsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}
	var sessionResponse response.SessionResponse
	json.Unmarshal(body, &sessionResponse)

	if len(sessionResponse.Sessions) != 2 {
		t.Errorf("expected sessionResponse.Sessions to have a len of 2: got %v", sessionResponse.Sessions)
	}

}

func TestSessionRouterHandlerNotAdmin(t *testing.T) {
	req, err := http.NewRequest("GET", "/sessions", nil)
	if err != nil {
		t.Fatal(err)
	}

	sessionRouter := createSessionRouter()
	adminUser, err := sessionRouter.Services.UserService.GetRepository().GetByName("admin")
	if err != nil {
		t.Fatal(err)
	}
	normalUser := addUserAndSession(t, *sessionRouter.Services, "user2", "user2", false)
	accessPayload := &token.AccessTokenPayload{UserId: normalUser.Id, IssuedAtTime: time.Now(), IsAdmin: normalUser.IsAdmin}
	addSession(t, *sessionRouter.Services, adminUser)
	accessToken, err := sessionRouter.Services.AccessTokenGenerator.CreateToken(accessPayload)
	if err != nil {
		t.Fatal(err)
	}

	accessCookie := &http.Cookie{Name: "accessToken", Value: accessToken, HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("accessToken=%s", accessCookie.Value))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sessionRouter.GetSessionsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}
	var sessionResponse response.SessionResponse
	json.Unmarshal(body, &sessionResponse)

	if len(sessionResponse.Sessions) != 1 {
		t.Errorf("expected sessionResponse.Sessions to have a len of 2: got %v", sessionResponse.Sessions)
	}

}

func TestSessionRouterHandlerInvalidAccessToken(t *testing.T) {
	req, err := http.NewRequest("GET", "/sessions", nil)
	if err != nil {
		t.Fatal(err)
	}

	sessionRouter := createSessionRouter()
	accessCookie := &http.Cookie{Name: "accessToken", Value: "123.123.123", HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("accessToken=%s", accessCookie.Value))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sessionRouter.GetSessionsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestDeleteSessionHandler(t *testing.T) {
	sessionRouter := createSessionRouter()
	user, err := sessionRouter.Services.UserService.GetRepository().GetByName("admin")
	if err != nil {
		t.Fatal(err)
	}

	accessPayload := &token.AccessTokenPayload{UserId: user.Id, IssuedAtTime: time.Now(), IsAdmin: user.IsAdmin}
	addSession(t, *sessionRouter.Services, user)
	accessToken, err := sessionRouter.Services.AccessTokenGenerator.CreateToken(accessPayload)
	if err != nil {
		t.Fatal(err)
	}
	addUserAndSession(t, *sessionRouter.Services, "user2", "user2", false)

	adminSession := sessionRouter.Services.SessionsHandler.GetUserSessions(user.Id)[0]

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/sessions/%s", adminSession.Id), nil)
	if err != nil {
		t.Fatal(err)
	}

	accessCookie := &http.Cookie{Name: "accessToken", Value: accessToken, HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("accessToken=%s", accessCookie.Value))

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/sessions/{id}", sessionRouter.DeleteSessionHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		body := strings.TrimSpace(rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v , body %s",
			status, http.StatusOK, body)
	}
}

func TestDeleteSessionHandlerInvalidAccessToken(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/sessions/123", nil)
	if err != nil {
		t.Fatal(err)
	}

	sessionRouter := createSessionRouter()
	accessCookie := &http.Cookie{Name: "accessToken", Value: "123.123.123", HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("accessToken=%s", accessCookie.Value))

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/sessions/{id}", sessionRouter.DeleteSessionHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestDeleteSessionHandlerNotFound(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/sessions/123", nil)
	if err != nil {
		t.Fatal(err)
	}

	sessionRouter := createSessionRouter()

	user, err := sessionRouter.Services.UserService.GetRepository().GetByName("admin")
	if err != nil {
		t.Fatal(err)
	}

	accessPayload := &token.AccessTokenPayload{UserId: user.Id, IssuedAtTime: time.Now(), IsAdmin: user.IsAdmin}
	addSession(t, *sessionRouter.Services, user)
	accessToken, err := sessionRouter.Services.AccessTokenGenerator.CreateToken(accessPayload)
	if err != nil {
		t.Fatal(err)
	}
	accessCookie := &http.Cookie{Name: "accessToken", Value: accessToken, HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("accessToken=%s", accessCookie.Value))

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/sessions/{id}", sessionRouter.DeleteSessionHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := `{"error":"Session not found"}`
	if want := strings.TrimSpace(rr.Body.String()); want != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", want, expected)
	}
}

func TestDeleteSessionHandlerForbidden(t *testing.T) {
	sessionRouter := createSessionRouter()
	adminUser, err := sessionRouter.Services.UserService.GetRepository().GetByName("admin")
	if err != nil {
		t.Fatal(err)
	}
	addSession(t, *sessionRouter.Services, adminUser)
	adminSession := sessionRouter.Services.SessionsHandler.GetUserSessions(adminUser.Id)[0]

	normalUser := addUserAndSession(t, *sessionRouter.Services, "user2", "user2", false)
	accessPayload := &token.AccessTokenPayload{UserId: normalUser.Id, IssuedAtTime: time.Now(), IsAdmin: normalUser.IsAdmin}
	accessToken, err := sessionRouter.Services.AccessTokenGenerator.CreateToken(accessPayload)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/sessions/%s", adminSession.Id), nil)
	if err != nil {
		t.Fatal(err)
	}

	accessCookie := &http.Cookie{Name: "accessToken", Value: accessToken, HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("accessToken=%s", accessCookie.Value))

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/sessions/{id}", sessionRouter.DeleteSessionHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		body := strings.TrimSpace(rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v , body %s",
			status, http.StatusForbidden, body)
	}
}
