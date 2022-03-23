package main

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/repositories/mock"
	"2022_1_OnlyGroup_back/app/usecases/impl"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const UrlUsers = "/users"
const ProfileUrl = "/profiles/{id:[0-9]+}"
const ProfileUrlShort = "/profiles/{id:[0-9]+}/shorts"
const ProfileUrlCandidates = "/profiles/candidates"

type APIServer struct {
	address        string
	authHandler    *handlers.AuthHandler
	profileHandler *handlers.ProfileHandler
	middlewares    handlers.Middlewares
}

func NewServer(addr string) APIServer {
	profileRepo := mock.NewProfileMock()
	authRepo := mock.NewAuthMock()
	return APIServer{address: addr, authHandler: handlers.CreateAuthHandler(impl.NewAuthUseCaseImpl(authRepo, profileRepo)),
		profileHandler: handlers.CreateProfileHandler(impl.NewProfileUseCaseImpl(profileRepo, authRepo)),
		middlewares:    handlers.MiddlewaresImpl{},
	}
}

func Cors(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
}

func (serv *APIServer) Run() error {
	multiplexor := mux.NewRouter()

	//Candidate for user
	multiplexor.HandleFunc(ProfileUrlCandidates, serv.profileHandler.GetCandidateHandler).Methods(http.MethodPost)
	//User own profile
	multiplexor.HandleFunc(ProfileUrl, serv.profileHandler.GetProfileHandler).Methods(http.MethodGet)
	multiplexor.HandleFunc(ProfileUrlShort, serv.profileHandler.GetProfileHandler).Methods(http.MethodGet) ///дописать
	multiplexor.HandleFunc(ProfileUrl, serv.profileHandler.ChangeProfileHandler).Methods(http.MethodPut)   //свой профиль

	multiplexor.Methods(http.MethodOptions).HandlerFunc(Cors)

	multiplexor.HandleFunc(UrlUsers, serv.authHandler.GET).Methods(http.MethodGet)
	multiplexor.HandleFunc(UrlUsers, serv.authHandler.PUT).Methods(http.MethodPut)
	multiplexor.HandleFunc(UrlUsers, serv.authHandler.POST).Methods(http.MethodPost)
	multiplexor.HandleFunc(UrlUsers, serv.authHandler.DELETE).Methods(http.MethodDelete)

	multiplexor.Use(serv.middlewares.AccessLogMiddleware)

	server := http.Server{Addr: serv.address, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, Handler: multiplexor}
	return server.ListenAndServe()
}
