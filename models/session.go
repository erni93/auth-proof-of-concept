package models

import "time"

type Session struct {
	RefreshToken AuthToken
	DeviceData   DeviceData
	Created      time.Time
}
