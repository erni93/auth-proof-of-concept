package user

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	IsAdmin  bool   `json:"isAdmin"`
	Password string `json:"-"`
}
