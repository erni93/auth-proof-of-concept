package router

import (
	response "authGo/router/response"
	"authGo/session"
	"authGo/token"
	"authGo/user"
	"authGo/validator"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func createUserRouter() *UserRouter {
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

	return &UserRouter{
		Services: services,
	}
}

func TestUserRouterGetUsersHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	userRouter := createUserRouter()
	adminUser, err := userRouter.Services.UserService.GetRepository().GetByName("admin")
	if err != nil {
		t.Fatal(err)
	}

	accessPayload := &token.AccessTokenPayload{UserId: adminUser.Id, IssuedAtTime: time.Now(), IsAdmin: adminUser.IsAdmin}
	addSession(t, *userRouter.Services, adminUser)
	accessToken, err := userRouter.Services.AccessTokenGenerator.CreateToken(accessPayload)
	if err != nil {
		t.Fatal(err)
	}
	addUserAndSession(t, *userRouter.Services, "user2", "user2", false)

	accessCookie := &http.Cookie{Name: "accessToken", Value: accessToken, HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("accessToken=%s", accessCookie.Value))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userRouter.GetUsersHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}
	var userResponse response.UserResponse
	json.Unmarshal(body, &userResponse)

	if len(userResponse.Users) != 2 {
		t.Errorf("expected userResponse.Users to have a len of 2: got %v", userResponse.Users)
	}

}

func TestUserRouterGetUsersHandlerInvalidAccessToken(t *testing.T) {
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	userRouter := createUserRouter()
	accessCookie := &http.Cookie{Name: "accessToken", Value: "123.123.123", HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("accessToken=%s", accessCookie.Value))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userRouter.GetUsersHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestUserRouterNewUserHandler(t *testing.T) {
	inputJson := validator.NewUserInput{Name: "user2", Password: "user2", IsAdmin: false}
	input, err := json.Marshal(&inputJson)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(input))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")

	userRouter := createUserRouter()
	adminUser, err := userRouter.Services.UserService.GetRepository().GetByName("admin")
	if err != nil {
		t.Fatal(err)
	}

	accessPayload := &token.AccessTokenPayload{UserId: adminUser.Id, IssuedAtTime: time.Now(), IsAdmin: adminUser.IsAdmin}
	addSession(t, *userRouter.Services, adminUser)
	accessToken, err := userRouter.Services.AccessTokenGenerator.CreateToken(accessPayload)
	if err != nil {
		t.Fatal(err)
	}

	accessCookie := &http.Cookie{Name: "accessToken", Value: accessToken, HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("accessToken=%s", accessCookie.Value))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userRouter.NewUserHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		body := strings.TrimSpace(rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v , body %s",
			status, http.StatusOK, body)
	}

	_, err = userRouter.Services.UserService.GetRepository().GetByName("user2")
	if err != nil {
		t.Error("new user not found")
	}
}

func TestUserRouterNewUserHandlerInvalidAccessToken(t *testing.T) {
	req, err := http.NewRequest("POST", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	userRouter := createUserRouter()
	accessCookie := &http.Cookie{Name: "accessToken", Value: "123.123.123", HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("accessToken=%s", accessCookie.Value))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userRouter.NewUserHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestUserRouterNewUserHandlerNotAdmin(t *testing.T) {
	req, err := http.NewRequest("POST", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	userRouter := createUserRouter()
	normalUser := addUserAndSession(t, *userRouter.Services, "user2", "user2", false)
	accessPayload := &token.AccessTokenPayload{UserId: normalUser.Id, IssuedAtTime: time.Now(), IsAdmin: normalUser.IsAdmin}
	accessToken, err := userRouter.Services.AccessTokenGenerator.CreateToken(accessPayload)
	if err != nil {
		t.Fatal(err)
	}

	accessCookie := &http.Cookie{Name: "accessToken", Value: accessToken, HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("accessToken=%s", accessCookie.Value))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userRouter.NewUserHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}
}

func TestUserRouterNewUserHandlerUserDataNotValid(t *testing.T) {
	inputJson := validator.NewUserInput{Name: "user2", Password: "user2", IsAdmin: false}
	input, err := json.Marshal(&inputJson)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(input))
	if err != nil {
		t.Fatal(err)
	}

	userRouter := createUserRouter()
	adminUser, err := userRouter.Services.UserService.GetRepository().GetByName("admin")
	if err != nil {
		t.Fatal(err)
	}

	accessPayload := &token.AccessTokenPayload{UserId: adminUser.Id, IssuedAtTime: time.Now(), IsAdmin: adminUser.IsAdmin}
	addSession(t, *userRouter.Services, adminUser)
	accessToken, err := userRouter.Services.AccessTokenGenerator.CreateToken(accessPayload)
	if err != nil {
		t.Fatal(err)
	}

	accessCookie := &http.Cookie{Name: "accessToken", Value: accessToken, HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("accessToken=%s", accessCookie.Value))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userRouter.NewUserHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := `{"error":"User data not valid"}`
	if want := strings.TrimSpace(rr.Body.String()); want != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", want, expected)
	}
}

func TestUserRouterNewUserHandlerErrorAddingUser(t *testing.T) {
	inputJson := validator.NewUserInput{Name: "admin", Password: "admin", IsAdmin: false}
	input, err := json.Marshal(&inputJson)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(input))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	userRouter := createUserRouter()
	adminUser, err := userRouter.Services.UserService.GetRepository().GetByName("admin")
	if err != nil {
		t.Fatal(err)
	}

	accessPayload := &token.AccessTokenPayload{UserId: adminUser.Id, IssuedAtTime: time.Now(), IsAdmin: adminUser.IsAdmin}
	addSession(t, *userRouter.Services, adminUser)
	accessToken, err := userRouter.Services.AccessTokenGenerator.CreateToken(accessPayload)
	if err != nil {
		t.Fatal(err)
	}

	accessCookie := &http.Cookie{Name: "accessToken", Value: accessToken, HttpOnly: true, Path: "/"}
	req.Header.Set("Cookie", fmt.Sprintf("accessToken=%s", accessCookie.Value))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userRouter.NewUserHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := `{"error":"Error creating new user"}`
	if want := strings.TrimSpace(rr.Body.String()); want != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", want, expected)
	}
}
