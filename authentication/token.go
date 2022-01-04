package authentication

import "time"

type Token struct {
	UserId       string
	IssuedAtTime time.Time
}
