package validator

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestGetNewUser(t *testing.T) {
	inputJson := NewUserInput{Name: "user1", Password: "user1", IsAdmin: false}
	input, err := json.Marshal(&inputJson)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(input))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	v := UserValidator{Validator: Validator{Request: req}}
	user, err := v.GetNewUser()
	if err != nil {
		t.Error("expected err to be nil")
	}
	if user == nil {
		t.Error("expected user to not be nil")
	}
}

func TestErrorsGetNewUser(t *testing.T) {
	input := "123456"
	req, err := http.NewRequest("POST", "/users", strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	v := UserValidator{Validator: Validator{Request: req}}
	_, err = v.GetNewUser()
	if err != ErrUserInvalidContentType {
		t.Errorf("expected err to be ErrUserInvalidContentType, got %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = v.GetNewUser()
	if err == nil {
		t.Error("expected err to not be nil")
	}
}
