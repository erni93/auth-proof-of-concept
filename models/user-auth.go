package models

type UserAuth struct {
	AccessToken  AuthToken
	RefreshToken AuthToken
}
