package authentication

import "time"

type TokenData struct {
	UserId       string
	IssuedAtTime time.Time
}
