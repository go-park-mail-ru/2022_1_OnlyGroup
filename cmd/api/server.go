package main

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/repositories/mock"
	"2022_1_OnlyGroup_back/app/repositories/postgres"
	"2022_1_OnlyGroup_back/app/usecases/impl"
	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"

	//_ "database/sql"
	"github.com/gorilla/mux"
	//_ "github.com/lib/pq"
	_ "github.com/jackc/pgx/stdlib"
	"log"
	"net/http"
	"time"
)

const UrlUsers = "/users"
const ProfileUrl = "/profiles/{id:[0-9]+}"
const ProfileUrlShort = "/profiles/{id:[0-9]+}/shorts"
const ProfileUrlCandidates = "/profiles/candidates"

const (
	host     = "localhost"
	port     = 5432
	user     = "kdv"
	password = "5051"
	dbname   = "kdv"
)

type APIServer struct {
	address        string
	authHandler    *handlers.AuthHandler
	profileHandler *handlers.ProfileHandler
}

func NewServer(addr string) APIServer {

	db, err := sqlx.Connect("pgx", "user=kdv password=5051 dbname=kdv sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	profRepo, err := postgres.NewProfilePostgres(db, "kdv", "users")
	if err != nil {
		log.Fatalln(err)
	}
	//m := models.Profile{FirstName: faker.FirstName(), LastName: faker.LastName(), Birthday: faker.Date(), City: "Moscow", Interests: []string{faker.Word(), faker.Word()}, AboutUser: faker.Sentence(), UserId: 1, Gender: faker.Gender()}
	log.Println(profRepo.CheckProfileFiled(2))

	profileRepo := mock.NewProfileMock()
	authRepo := mock.NewAuthMock()
	return APIServer{address: addr, authHandler: handlers.CreateAuthHandler(impl.NewAuthUseCaseImpl(authRepo, profileRepo)),
		profileHandler: handlers.CreateProfileHandler(impl.NewProfileUseCaseImpl(profileRepo, authRepo)),
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

	server := http.Server{Addr: serv.address, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, Handler: multiplexor}
	return server.ListenAndServe()
}
