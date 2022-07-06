package authentication

import (
	"authGo/user"
	"crypto/hmac"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type TokenGenerator struct {
	accessTokenPassword  []byte
	refreshTokenPassword []byte
}

func NewTokenGenerator(accessTokenPassword string, refreshTokenPassword string) *TokenGenerator {
	return &TokenGenerator{accessTokenPassword: []byte(accessTokenPassword), refreshTokenPassword: []byte(refreshTokenPassword)}
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
	return t.validateJWT(jwt, t.accessTokenPassword)
}

func (t *TokenGenerator) IsRefreshTokenValid(jwt string) (bool, error) {
	return t.validateJWT(jwt, t.refreshTokenPassword)
}

func (t *TokenGenerator) validateJWT(jwt string, password []byte) (bool, error) {
	jwtParts := strings.Split(jwt, ".")
	if len(jwtParts) != 3 {
		return false, fmt.Errorf("token generator: jwt length not valid, len %d", len(jwtParts))
	}
	signatureB64 := t.createSignature(jwtParts[0], jwtParts[1], password)
	if signatureB64 != jwtParts[2] {
		return false, fmt.Errorf("token generator: signature not valid, got %s want %s", jwtParts[2], signatureB64)
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
