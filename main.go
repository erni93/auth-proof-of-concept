package main

import (
	"authGo/router"
	"authGo/session"
	"authGo/token"
	"authGo/user"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	port := ":8181"
	loginRouter := &router.LoginRouter{
		Services: &router.LoginRouterServices{
			UserService:           user.NewUserService(),
			AccessTokenGenerator:  &token.TokenGenerator[token.AccessTokenPayload]{Password: []byte("accessKey"), Duration: time.Minute * 2},
			RefreshTokenGenerator: &token.TokenGenerator[token.RefreshTokenPayload]{Password: []byte("refreshKey"), Duration: time.Hour * 24 * 365},
		},
		SessionHandler: session.NewSessionHandler(),
	}

	router := mux.NewRouter()
	router.HandleFunc("/auth/login", loginRouter.Handler).Methods("POST")
	http.Handle("/", router)
	log.Printf("Application listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}
