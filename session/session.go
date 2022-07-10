package session

import (
	"authGo/token"
	"time"
)

type Session struct {
	Id         string
	UserToken  token.RefreshTokenPayload
	DeviceData DeviceData
	LastUpdate time.Time
}
