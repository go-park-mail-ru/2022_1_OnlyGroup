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
	address        string
	authHandler    *handlers.AuthHandler
	profileHandler *handlers.ProfileHandler
}

func NewServer(addr string) APIServer {
	profileRepo := mock.NewProfileMock()
	authRepo := mock.NewAuthMock()
	return APIServer{address: addr, authHandler: handlers.CreateAuthHandler(impl.NewAuthUseCaseImpl(authRepo)),
		profileHandler: handlers.CreateProfileHandler(impl.NewProfileUseCaseImpl(profileRepo, authRepo)),
	}
}

func (serv *APIServer) Run() error {
	multiplexor := mux.NewRouter()

	multiplexor.HandleFunc("/user", serv.authHandler.AuthUserHandler).Methods(http.MethodGet)
	multiplexor.HandleFunc("/user", serv.authHandler.LoginUserHandler).Methods(http.MethodPut)
	multiplexor.HandleFunc("/user", serv.authHandler.LogupUserHandler).Methods(http.MethodPost)
	multiplexor.HandleFunc("/user", serv.authHandler.LogoutUserHandler).Methods(http.MethodDelete)

	multiplexor.HandleFunc("/auth", serv.authHandler.MainAuthHandler)
	multiplexor.HandleFunc("/register", serv.authHandler.RegisterAuthHandler)
	//Candidate for user
	multiplexor.HandleFunc("/profile", serv.profileHandler.GetCandidateHandler).Methods(http.MethodGet)
	//User own profile
	multiplexor.HandleFunc("/profile/my", serv.profileHandler.GetProfileHandler).Methods("GET")
	//multiplexor.HandleFunc("/profile/my", serv.profileHandler.GetProfileHandler).Methods("POST")
	multiplexor.HandleFunc("/profile/my", serv.profileHandler.ChangeProfileHandler).Methods("PUT")
	//multiplexor.HandleFunc("/profile/my", serv.profileHandler.DeleteProfileHandler).Methods("DELETE")

	server := http.Server{Addr: serv.address, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, Handler: multiplexor}
	return server.ListenAndServe()
}
