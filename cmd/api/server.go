package main

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/repositories/mock"
	"2022_1_OnlyGroup_back/app/usecases/impl"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type APIServer struct {
	address     string
	authHandler *handlers.AuthHandler
}

func NewServer(addr string) APIServer {
	return APIServer{address: addr, authHandler: handlers.CreateAuthHandler(impl.NewAuthUseCaseImpl(mock.NewAuthMock()))}
}

func (serv *APIServer) Run() error {
	multiplexor := mux.NewRouter()
	multiplexor.HandleFunc("/auth", serv.authHandler.MainAuthHandler)
	multiplexor.HandleFunc("/register", serv.authHandler.RegisterAuthHandler)
	server := http.Server{Addr: serv.address, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, Handler: multiplexor}
	return server.ListenAndServe()
}
