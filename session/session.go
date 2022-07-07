package session

import (
	"authGo/token"
	"time"
)

type Session struct {
	UserToken  token.RefreshTokenPayload
	DeviceData DeviceData
	Created    time.Time
}
