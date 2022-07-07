package authentication

import "time"

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

var DefaultHeader = &Header{Alg: "HS256", Typ: "JWT"}

type TokenPayload interface {
	AccessTokenPayload | RefreshTokenPayload
}

type AccessTokenPayload struct {
	UserId       string    `json:"userId"`
	IssuedAtTime time.Time `json:"issuedAtTime"`
	IsAdmin      bool      `json:"isAdmin"`
}

type RefreshTokenPayload struct {
	UserId       string    `json:"userId"`
	IssuedAtTime time.Time `json:"issuedAtTime"`
}

type IssuedAtTime struct {
	IssuedAtTime time.Time `json:"issuedAtTime"`
}
