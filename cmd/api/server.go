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
	multiplexor.HandleFunc("/users", serv.authHandler.GET).Methods(http.MethodGet)
	multiplexor.HandleFunc("/users", serv.authHandler.PUT).Methods(http.MethodPut)
	multiplexor.HandleFunc("/users", serv.authHandler.POST).Methods(http.MethodPost)
	multiplexor.HandleFunc("/users", serv.authHandler.DELETE).Methods(http.MethodDelete)

	server := http.Server{Addr: serv.address, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, Handler: multiplexor}
	return server.ListenAndServe()
}
