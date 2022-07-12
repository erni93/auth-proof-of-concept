package validator

import (
	"authGo/token"
	"errors"
	"fmt"
)

type AccessTokenValidator struct {
	Validator Validator
}

var (
	ErrAccessReadingAccessCookie = errors.New("access token validator: error reading accessToken cookie")
)

func (v *AccessTokenValidator) ValidateAccessToken() (*token.AccessTokenPayload, error) {
	accessToken, err := v.Validator.Request.Cookie("accessToken")
	if err != nil {
		return nil, fmt.Errorf("%w, %s", ErrAccessReadingAccessCookie, err)
	}
	err = v.Validator.Services.AccessTokenGenerator.IsTokenValid(accessToken.Value)
	if err != nil {
		return nil, err
	}
	payload := &token.AccessTokenPayload{}
	err = v.Validator.Services.AccessTokenGenerator.LoadPayload(accessToken.Value, payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
