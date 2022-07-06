package authentication

import "time"

type Header struct {
	alg string
	typ string
}

var DefaultHeader = &Header{alg: "HS256", typ: "JWT"}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
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
