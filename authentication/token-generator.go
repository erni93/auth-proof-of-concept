package authentication

import (
	"authGo/user"
	"crypto/hmac"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidJWTLength = errors.New("token generator: jwt length not valid")
	ErrTokenExpired     = errors.New("token generator: expired issuedAtTime")
	ErrInvalidSignature = errors.New("token generator: signature not valid")
)

type TokenGenerator struct {
	accessTokenPassword  []byte
	refreshTokenPassword []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func (t *TokenGenerator) CreateAccessToken(user *user.User, issuedAtTime time.Time) (string, error) {
	header, err := json.Marshal(DefaultHeader)
	if err != nil {
		return "", err
	}
	payload, err := json.Marshal(&AccessTokenPayload{UserId: user.Id, IssuedAtTime: issuedAtTime, IsAdmin: user.IsAdmin})
	if err != nil {
		return "", err
	}
	return t.createJWT(header, payload, t.accessTokenPassword), nil
}

func (t *TokenGenerator) CreateRefreshToken(user *user.User, issuedAtTime time.Time) (string, error) {
	header, err := json.Marshal(DefaultHeader)
	if err != nil {
		return "", err
	}
	payload, err := json.Marshal(&RefreshTokenPayload{UserId: user.Id, IssuedAtTime: issuedAtTime})
	if err != nil {
		return "", err
	}
	return t.createJWT(header, payload, t.refreshTokenPassword), nil
}

func (t *TokenGenerator) IsAccessTokenValid(jwt string) (bool, error) {
	expirationDate := time.Now().Add(-t.accessTokenDuration)
	return t.validateJWT(jwt, t.accessTokenPassword, expirationDate)
}

func (t *TokenGenerator) IsRefreshTokenValid(jwt string) (bool, error) {
	expirationDate := time.Now().Add(-t.refreshTokenDuration)
	return t.validateJWT(jwt, t.refreshTokenPassword, expirationDate)
}

func (t *TokenGenerator) validateJWT(jwt string, password []byte, expirationDate time.Time) (bool, error) {
	jwtParts := strings.Split(jwt, ".")
	if len(jwtParts) != 3 {
		return false, fmt.Errorf("%w, len %d", ErrInvalidJWTLength, len(jwtParts))
	}

	var issuedAtTime IssuedAtTime
	payloadJson, err := b64.RawURLEncoding.DecodeString(strings.Split(jwt, ".")[1])
	if err != nil {
		return false, err
	}
	err = json.Unmarshal([]byte(payloadJson), &issuedAtTime)
	if err != nil {
		return false, err
	}
	if issuedAtTime.IssuedAtTime.Before(expirationDate) {
		return false, fmt.Errorf("%w, got %s expirationDate %s", ErrTokenExpired, issuedAtTime.IssuedAtTime.String(), expirationDate.String())
	}

	signatureB64 := t.createSignature(jwtParts[0], jwtParts[1], password)
	if signatureB64 != jwtParts[2] {
		return false, fmt.Errorf("%w, got %s want %s", ErrInvalidSignature, jwtParts[2], signatureB64)
	}
	return true, nil
}

func (t *TokenGenerator) createJWT(header []byte, payload []byte, password []byte) string {
	headerB64 := b64.RawURLEncoding.EncodeToString(header)
	payloadB64 := b64.RawURLEncoding.EncodeToString(payload)
	signatureB64 := t.createSignature(headerB64, payloadB64, password)
	return fmt.Sprintf("%s.%s.%s", headerB64, payloadB64, signatureB64)
}

func (t *TokenGenerator) createSignature(headerB64 string, payloadB64 string, password []byte) string {
	signature := fmt.Sprintf("%s.%s", headerB64, payloadB64)
	h := hmac.New(sha256.New, password)
	h.Write([]byte(signature))
	return b64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
