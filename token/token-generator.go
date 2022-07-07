package token

import (
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

type TokenGenerator[T TokenPayload] struct {
	Password []byte
	Duration time.Duration
}

func (t *TokenGenerator[T]) CreateToken(payload *T) (string, error) {
	header, err := json.Marshal(DefaultHeader)
	if err != nil {
		return "", err
	}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return t.createJWT(header, payloadJson), nil
}

func (t *TokenGenerator[T]) IsTokenValid(jwt string) (bool, error) {
	expirationDate := time.Now().Add(-t.Duration)
	return t.validateJWT(jwt, expirationDate)
}

func (t *TokenGenerator[T]) createJWT(header []byte, payload []byte) string {
	headerB64 := b64.RawURLEncoding.EncodeToString(header)
	payloadB64 := b64.RawURLEncoding.EncodeToString(payload)
	signatureB64 := t.createSignature(headerB64, payloadB64)
	return fmt.Sprintf("%s.%s.%s", headerB64, payloadB64, signatureB64)
}

func (t *TokenGenerator[T]) createSignature(headerB64 string, payloadB64 string) string {
	signature := fmt.Sprintf("%s.%s", headerB64, payloadB64)
	h := hmac.New(sha256.New, t.Password)
	h.Write([]byte(signature))
	return b64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func (t *TokenGenerator[T]) validateJWT(jwt string, expirationDate time.Time) (bool, error) {
	jwtParts := strings.Split(jwt, ".")
	if len(jwtParts) != 3 {
		return false, fmt.Errorf("%w, len %d", ErrInvalidJWTLength, len(jwtParts))
	}

	if err := t.validateIssuedAtTime(jwt, expirationDate); err != nil {
		return false, err
	}

	signatureB64 := t.createSignature(jwtParts[0], jwtParts[1])
	if signatureB64 != jwtParts[2] {
		return false, fmt.Errorf("%w, got %s want %s", ErrInvalidSignature, jwtParts[2], signatureB64)
	}
	return true, nil
}

func (t *TokenGenerator[T]) validateIssuedAtTime(jwt string, expirationDate time.Time) error {
	var issuedAtTime IssuedAtTime
	payloadJson, err := b64.RawURLEncoding.DecodeString(strings.Split(jwt, ".")[1])
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(payloadJson), &issuedAtTime)
	if err != nil {
		return err
	}
	if issuedAtTime.IssuedAtTime.Before(expirationDate) {
		return fmt.Errorf("%w, got %s expirationDate %s", ErrTokenExpired, issuedAtTime.IssuedAtTime.String(), expirationDate.String())
	}
	return nil
}
