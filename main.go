package main

import (
	"authGo/router"
	"authGo/session"
	"authGo/token"
	"authGo/user"
	"authGo/validator"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	port := ":8181"

	userService := user.NewUserService()
	userService.CreateUser("admin", "admin", true)
	accessTokenGenerator := &token.TokenGenerator[token.AccessTokenPayload]{Password: []byte("accessKey"), Duration: time.Minute * 2}
	refreshTokenGenerator := &token.TokenGenerator[token.RefreshTokenPayload]{Password: []byte("refreshKey"), Duration: time.Hour * 24 * 365}
	sessionHandler := session.NewSessionHandler()

	services := &validator.Services{
		UserService:           userService,
		AccessTokenGenerator:  accessTokenGenerator,
		RefreshTokenGenerator: refreshTokenGenerator,
		SessionsHandler:       sessionHandler,
	}
	loginRouter := &router.LoginRouter{
		Services: services,
	}
	refreshRouter := &router.RefreshRouter{
		Services: services,
	}
	userRouter := &router.UserRouter{
		Services: services,
	}

	router := mux.NewRouter()
	router.HandleFunc("/auth/login", loginRouter.Handler).Methods("POST")
	router.HandleFunc("/auth/refresh", refreshRouter.Handler).Methods("POST")
	router.HandleFunc("/users", userRouter.GetUsersHandler).Methods("GET")
	http.Handle("/", router)
	log.Printf("Application listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}
