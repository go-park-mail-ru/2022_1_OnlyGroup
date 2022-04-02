package main

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/repositories/postgres"
	redis_repo "2022_1_OnlyGroup_back/app/repositories/redis"
	"2022_1_OnlyGroup_back/app/usecases/impl"
	impl2 "2022_1_OnlyGroup_back/pkg/fileService/impl"
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
const UrlPhotosIdPostfix = "/photos/{id:[0-9]+}"
const UrlPhotosIdParamsPostfix = "/photos/{id:[0-9]+}/params"
const UrlProfilePhotosPostfix = "/profile/{id:[0-9]+}/photos"
const UrlProfilePhotosAvatarPostfix = "/profile/{id:[0-9]+}/photos/avatar"

type APIServer struct {
	conf           APIServerConf
	authHandler    *handlers.AuthHandler
	profileHandler *handlers.ProfileHandler
	photosHandler  *handlers.PhotosHandler
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
	photosRepo, err := postgres.NewPostgresPhotoRepository(postgresConnect, conf.PostgresConf.PhotosDbTableName, conf.PostgresConf.UsersDbTableName, conf.PostgresConf.AvatarDbTableName)
	if err != nil {
		return APIServer{}, err
	}

	sessionsRepo := redis_repo.NewRedisSessionRepository(redisConnect, conf.RedisConf.SessionsPrefix, sessionGenerator.NewRandomGenerator())

	//useCases
	authUseCase := impl.NewAuthUseCaseImpl(usersRepo, sessionsRepo, profilesRepo)
	profileUseCase := impl.NewProfileUseCaseImpl(profilesRepo)
	photosUseCase := impl.NewPhotosUseCase(photosRepo)

	//fileService
	photosService, err := impl2.NewFileServicePhotos(conf.PhotosStorageDirectory)
	if err != nil {
		return APIServer{}, err
	}

	return APIServer{
		conf:           conf,
		authHandler:    handlers.CreateAuthHandler(authUseCase),
		profileHandler: handlers.CreateProfileHandler(profileUseCase),
		photosHandler:  handlers.CreatePhotosHandler(photosUseCase, photosService),
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
	UrlPhotosId := serv.conf.ApiPathPrefix + UrlPhotosIdPostfix
	UrlPhotosIdParams := serv.conf.ApiPathPrefix + UrlPhotosIdParamsPostfix
	UrlProfilePhotos := serv.conf.ApiPathPrefix + UrlProfilePhotosPostfix
	UrlProfilePhotosAvatar := serv.conf.ApiPathPrefix + UrlProfilePhotosAvatarPostfix

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
	multiplexorWithAuth := multiplexor.PathPrefix("").Subrouter()

	multiplexorWithAuth.Use(serv.middlewares.CheckAuthMiddleware)
	//сandidate methods
	multiplexorWithAuth.HandleFunc(UrlProfileCandidates, serv.profileHandler.GetCandidateHandler).Methods(http.MethodPost)
	//profile methods
	multiplexorWithAuth.HandleFunc(UrlProfileId, serv.profileHandler.GetProfileHandler).Methods(http.MethodGet)
	multiplexorWithAuth.HandleFunc(UrlProfileIdShort, serv.profileHandler.GetShortProfileHandler).Methods(http.MethodGet) ///дописать
	multiplexorWithAuth.HandleFunc(UrlProfileId, serv.profileHandler.ChangeProfileHandler).Methods(http.MethodPut)        //свой профиль

	//photos
	multiplexorWithAuth.HandleFunc(UrlPhotosId, serv.photosHandler.GETPhoto).Methods(http.MethodGet)
	multiplexorWithAuth.HandleFunc(UrlPhotosId, serv.photosHandler.POSTPhoto).Methods(http.MethodPost)
	multiplexorWithAuth.HandleFunc(UrlPhotosId, serv.photosHandler.DELETE).Methods(http.MethodDelete)
	multiplexorWithAuth.HandleFunc(UrlPhotosIdParams, serv.photosHandler.PUTParams).Methods(http.MethodPut)
	multiplexorWithAuth.HandleFunc(UrlPhotosIdParams, serv.photosHandler.GETParams).Methods(http.MethodGet)
	multiplexorWithAuth.HandleFunc(UrlProfilePhotos, serv.photosHandler.GETAll).Methods(http.MethodGet)
	multiplexorWithAuth.HandleFunc(UrlProfilePhotosAvatar, serv.photosHandler.GETAvatar).Methods(http.MethodGet)
	multiplexorWithAuth.HandleFunc(UrlProfilePhotosAvatar, serv.photosHandler.PUTAvatar).Methods(http.MethodPut)

	serverAddr := serv.conf.ServerAddr + ":" + serv.conf.ServerPort
	server := http.Server{Addr: serverAddr, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, Handler: multiplexor}
	return server.ListenAndServe()
}
