package models

import "time"

type AuthToken struct {
	UserId       string
	IssuedAtTime time.Time
}
