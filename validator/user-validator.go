package validator

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type NewUserInput struct {
	Name     string `json:"name"`
	IsAdmin  bool   `json:"isAdmin"`
	Password string `json:"password"`
}

type UserValidator struct {
	Validator Validator
}

var (
	ErrUserInvalidContentType = errors.New("user validator: invalid content-type")
)

func (v *UserValidator) GetNewUser() (*NewUserInput, error) {
	contentType := v.Validator.Request.Header.Get("Content-type")
	if contentType != "application/json" {
		return nil, ErrUserInvalidContentType
	}

	body, err := ioutil.ReadAll(v.Validator.Request.Body)
	if err != nil {
		return nil, err
	}

	var newUser NewUserInput
	err = json.Unmarshal([]byte(string(body)), &newUser)
	if err != nil {
		return nil, err
	}
	return &newUser, nil
}
