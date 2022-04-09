package main

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/repositories/postgres"
	redis_repo "2022_1_OnlyGroup_back/app/repositories/redis"
	_ "2022_1_OnlyGroup_back/app/usecases"
	"2022_1_OnlyGroup_back/app/usecases/impl"
	"2022_1_OnlyGroup_back/pkg/dataValidator"
	"2022_1_OnlyGroup_back/pkg/sessionGenerator"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"
)

const UrlUsersPostfix = "/users"

const UrlProfileIdPostfix = "/profiles/{id:[0-9]+}"
const UrlProfileIdShortPostfix = "/profiles/{id:[0-9]+}/shorts"
const UrlProfileCandidatesPostfix = "/profiles/candidates"

const UrlLikesPostfix = "/likes"

type APIServer struct {
	conf           APIServerConf
	authHandler    *handlers.AuthHandler
	profileHandler *handlers.ProfileHandler
	likesHandler   *handlers.LikesHandler
	middlewares    handlers.Middlewares
}

func NewServer(conf APIServerConf) (APIServer, error) {
	//connections
	postgresSourceString := "postgresql://" + conf.PostgresConf.Username + ":" + conf.PostgresConf.Password +
		"@" + conf.PostgresConf.Address + ":" + conf.PostgresConf.Port + "/" + conf.PostgresConf.DbName + "?" + conf.PostgresConf.Params
	postgresConnect, err := sqlx.Open("pgx", postgresSourceString)
	if err != nil {
		return APIServer{}, err
	}
	redisConnOptions := redis.Options{
		Username: conf.RedisConf.Username,
		Password: conf.RedisConf.Password,
		Addr:     conf.RedisConf.Address + ":" + conf.RedisConf.Port,
		DB:       conf.RedisConf.DbNum,
	}
	redisConnect := redis.NewClient(&redisConnOptions)
	_, err = redisConnect.Ping(context.TODO()).Result()
	if err != nil {
		return APIServer{}, err
	}
	//repositories
	usersRepo, err := postgres.NewPostgresUsersRepo(postgresConnect, conf.PostgresConf.UsersDbTableName)
	if err != nil {
		return APIServer{}, err
	}
	profilesRepo, err := postgres.NewProfilePostgres(postgresConnect, conf.PostgresConf.ProfilesDbTableName, conf.PostgresConf.UsersDbTableName, conf.PostgresConf.InterestsDbTableName)
	if err != nil {
		return APIServer{}, err
	}
	likesRepo, err := postgres.NewLikesPostgres(postgresConnect, conf.PostgresConf.LikesDbTableName, conf.PostgresConf.UsersDbTableName)
	if err != nil {
		return APIServer{}, err
	}
	sessionsRepo := redis_repo.NewRedisSessionRepository(redisConnect, conf.RedisConf.SessionsPrefix, sessionGenerator.NewRandomGenerator())
	//set validators
	dataValidator.SetValidators()
	//useCases
	authUseCase := impl.NewAuthUseCaseImpl(usersRepo, sessionsRepo, profilesRepo)
	profileUseCase := impl.NewProfileUseCaseImpl(profilesRepo)
	likeUseCase := impl.NewLikesUseCaseImpl(likesRepo)

	return APIServer{
		conf:           conf,
		authHandler:    handlers.CreateAuthHandler(authUseCase),
		profileHandler: handlers.CreateProfileHandler(profileUseCase),
		likesHandler:   handlers.CreateLikesHandler(likeUseCase),
		middlewares:    handlers.MiddlewaresImpl{AuthUseCase: authUseCase},
	}, nil
}

func Cors(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
}

func (serv *APIServer) Run() error {
	UrlUsers := serv.conf.ApiPathPrefix + UrlUsersPostfix
	UrlProfileId := serv.conf.ApiPathPrefix + UrlProfileIdPostfix
	UrlProfileIdShort := serv.conf.ApiPathPrefix + UrlProfileIdShortPostfix
	UrlProfileCandidates := serv.conf.ApiPathPrefix + UrlProfileCandidatesPostfix
	UrlLikes := serv.conf.ApiPathPrefix + UrlLikesPostfix

	//main multiplexor
	multiplexor := mux.NewRouter()

	multiplexor.Use(serv.middlewares.AccessLogMiddleware)
	multiplexor.Use(serv.middlewares.PanicMiddleware)
	multiplexor.Use(serv.middlewares.CorsMiddleware)
	multiplexor.Methods(http.MethodOptions).HandlerFunc(Cors)

	multiplexor.HandleFunc(UrlUsers, serv.authHandler.GET).Methods(http.MethodGet)
	multiplexor.HandleFunc(UrlUsers, serv.authHandler.PUT).Methods(http.MethodPut)
	multiplexor.HandleFunc(UrlUsers, serv.authHandler.POST).Methods(http.MethodPost)
	multiplexor.HandleFunc(UrlUsers, serv.authHandler.DELETE).Methods(http.MethodDelete)

	//profile multiplexor
	multiplexorProfile := multiplexor.PathPrefix("").Subrouter()

	multiplexorProfile.Use(serv.middlewares.CheckAuthMiddleware)
	//сandidate methods
	multiplexorProfile.HandleFunc(UrlProfileCandidates, serv.profileHandler.GetCandidateHandler).Methods(http.MethodPost)
	//profile methods

	multiplexorProfile.HandleFunc(UrlLikes, serv.likesHandler.Set).Methods(http.MethodPost)
	multiplexorProfile.HandleFunc(UrlLikes, serv.likesHandler.Get).Methods(http.MethodGet)

	multiplexorProfile.HandleFunc(UrlProfileId, serv.profileHandler.GetProfileHandler).Methods(http.MethodGet)
	multiplexorProfile.HandleFunc(UrlProfileIdShort, serv.profileHandler.GetShortProfileHandler).Methods(http.MethodGet) ///дописать
	multiplexorProfile.HandleFunc(UrlProfileId, serv.profileHandler.ChangeProfileHandler).Methods(http.MethodPut)        //свой профиль

	serverAddr := serv.conf.ServerAddr + ":" + serv.conf.ServerPort
	server := http.Server{Addr: serverAddr, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, Handler: multiplexor}
	return server.ListenAndServe()
}
