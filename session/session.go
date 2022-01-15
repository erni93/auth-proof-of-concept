package session

import (
	"authGo/authentication"
	"time"
)

type Session struct {
	UserToken  authentication.TokenData
	DeviceData DeviceData
	Created    time.Time
}
