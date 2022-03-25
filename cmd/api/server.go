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

//const ProfileUrl = "/profiles"
const ProfileIdUrl = "/profiles/{id:[0-9]+}"
const ProfileUrlShort = "/profiles/{id:[0-9]+}/shorts"
const ProfileUrlCandidates = "/candidates"

type APIServer struct {
	address        string
	authHandler    *handlers.AuthHandler
	profileHandler *handlers.ProfileHandler
	middlewares    handlers.Middlewares
}

func NewServer(addr string) APIServer {
	profileRepo := mock.NewProfileMock()
	authRepo := mock.NewAuthMock()
	authUseCase := impl.NewAuthUseCaseImpl(authRepo, profileRepo)
	return APIServer{address: addr, authHandler: handlers.CreateAuthHandler(authUseCase),
		profileHandler: handlers.CreateProfileHandler(impl.NewProfileUseCaseImpl(profileRepo, authRepo)),
		middlewares:    handlers.MiddlewaresImpl{AuthUseCase: authUseCase},
	}
}

func Cors(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
}

func (serv *APIServer) Run() error {
	//main multiplexor
	multiplexor := mux.NewRouter()

	multiplexor.HandleFunc(UrlUsers, serv.authHandler.GET).Methods(http.MethodGet)
	multiplexor.HandleFunc(UrlUsers, serv.authHandler.PUT).Methods(http.MethodPut)
	multiplexor.HandleFunc(UrlUsers, serv.authHandler.POST).Methods(http.MethodPost)
	multiplexor.HandleFunc(UrlUsers, serv.authHandler.DELETE).Methods(http.MethodDelete)

	multiplexor.Use(serv.middlewares.AccessLogMiddleware)
	multiplexor.Use(serv.middlewares.PanicMiddleware)

	multiplexor.Methods(http.MethodOptions).HandlerFunc(Cors)

	//profile multiplexor
	multiplexorProfile := multiplexor.PathPrefix("").Subrouter()

	//Candidate methods
	multiplexorProfile.HandleFunc(ProfileUrlCandidates, serv.profileHandler.GetCandidateHandler).Methods(http.MethodPost)
	//profile methods
	multiplexorProfile.HandleFunc(ProfileIdUrl, serv.profileHandler.GetProfileHandler).Methods(http.MethodGet)
	multiplexorProfile.HandleFunc(ProfileUrlShort, serv.profileHandler.GetShortProfileHandler).Methods(http.MethodGet) ///дописать
	multiplexorProfile.HandleFunc(ProfileIdUrl, serv.profileHandler.ChangeProfileHandler).Methods(http.MethodPut)      //свой профиль
	//profile middlewares
	multiplexorProfile.Use(serv.middlewares.CheckAuthMiddleware)

	server := http.Server{Addr: serv.address, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, Handler: multiplexor}
	return server.ListenAndServe()
}
