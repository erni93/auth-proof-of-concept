package session

import (
	"authGo/authentication"
	"time"
)

type Session struct {
	UserToken  authentication.RefreshTokenPayload
	DeviceData DeviceData
	Created    time.Time
}
