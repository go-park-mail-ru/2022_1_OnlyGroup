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

func CorsMock(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
}

func (serv *APIServer) Run() error {
	multiplexor := mux.NewRouter()
	multiplexor.HandleFunc("/user", serv.authHandler.AuthUserHandler).Methods(http.MethodGet)
	multiplexor.HandleFunc("/user", serv.authHandler.LoginUserHandler).Methods(http.MethodPut)
	multiplexor.HandleFunc("/user", serv.authHandler.LogupUserHandler).Methods(http.MethodPost)
	multiplexor.HandleFunc("/user", serv.authHandler.LogoutUserHandler).Methods(http.MethodDelete)
	multiplexor.HandleFunc("/user", CorsMock).Methods(http.MethodOptions)

	server := http.Server{Addr: serv.address, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, Handler: multiplexor}
	return server.ListenAndServe()
}
