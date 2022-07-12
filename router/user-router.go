package router

import (
	"authGo/validator"
	"log"
	"net/http"
)

type UserRouter struct {
	Services *validator.Services
}

func (u *UserRouter) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	v := validator.AccessTokenValidator{Validator: validator.Validator{Writer: w, Request: r, Services: u.Services}}

	_, err := v.ValidateAccessToken()
	if err != nil {
		log.Print(err)
		writeTokenError(w)
		return
	}

	users := v.Validator.Services.UserService.GetRepository().GetAll()
	writeUserList(w, users)
}
