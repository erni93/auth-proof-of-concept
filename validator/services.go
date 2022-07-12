package validator

import (
	"authGo/session"
	"authGo/token"
	"authGo/user"
)

type Services struct {
	UserService           *user.UserService
	AccessTokenGenerator  *token.TokenGenerator[token.AccessTokenPayload]
	RefreshTokenGenerator *token.TokenGenerator[token.RefreshTokenPayload]
	SessionsHandler       *session.SessionsHandler
}
