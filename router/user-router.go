package router

import (
	response "authGo/router/response"
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
		response.WriteTokenError(w)
		return
	}

	users := v.Validator.Services.UserService.GetRepository().GetAll()
	response.WriteUserList(w, users)
}

func (u *UserRouter) NewUserHandler(w http.ResponseWriter, r *http.Request) {
	tokenV := validator.AccessTokenValidator{Validator: validator.Validator{Writer: w, Request: r, Services: u.Services}}
	userV := validator.UserValidator{Validator: validator.Validator{Writer: w, Request: r, Services: u.Services}}

	payload, err := tokenV.ValidateAccessToken()
	if err != nil {
		log.Print(err)
		response.WriteTokenError(w)
		return
	}

	if !payload.IsAdmin {
		response.WriteForbidden(w)
		return
	}

	newUser, err := userV.GetNewUser()
	if err != nil {
		log.Print(err)
		response.WriteError(w, "User data not valid")
		return
	}

	err = u.Services.UserService.CreateUser(newUser.Name, newUser.Password, newUser.IsAdmin)
	if err != nil {
		log.Print(err)
		response.WriteError(w, "Error creating new user")
		return
	}
	w.WriteHeader(http.StatusOK)
}
