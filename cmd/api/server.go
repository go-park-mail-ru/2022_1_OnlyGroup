package main

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/repositories/postgres"
	redis_repo "2022_1_OnlyGroup_back/app/repositories/redis"
	_ "2022_1_OnlyGroup_back/app/usecases"
	"2022_1_OnlyGroup_back/app/usecases/impl"
	"2022_1_OnlyGroup_back/pkg/dataValidator"
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
const UrlCSRFPostfix = "/csrf"
const UrlPhotosPostfix = "/photos"
const UrlPhotosIdPostfix = "/photos/{id:[0-9]+}"
const UrlPhotosIdParamsPostfix = "/photos/{id:[0-9]+}/params"
const UrlProfilePhotosPostfix = "/profile/{id:[0-9]+}/photos"
const UrlProfilePhotosAvatarPostfix = "/profile/{id:[0-9]+}/photos/avatar"

const UrlLikesPostfix = "/likes"

type APIServer struct {
	conf           APIServerConf
	authHandler    *handlers.AuthHandler
	profileHandler *handlers.ProfileHandler
	photosHandler  *handlers.PhotosHandler
	likesHandler   *handlers.LikesHandler
	middlewares    handlers.Middlewares
	jwtHandler     *handlers.CSRFHandler
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
	//jwtToken
	randGenerator := sessionGenerator.NewRandomGenerator()
	jwt, _ := handlers.NewJwtToken(randGenerator.String(16))
	//useCases
	photosRepo, err := postgres.NewPostgresPhotoRepository(postgresConnect, conf.PostgresConf.PhotosDbTableName, conf.PostgresConf.UsersDbTableName, conf.PostgresConf.AvatarDbTableName)
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
	photosUseCase := impl.NewPhotosUseCase(photosRepo)

	//fileService
	photosService, err := impl2.NewFileServicePhotos(conf.PhotosStorageDirectory)
	if err != nil {
		return APIServer{}, err
	}
	likeUseCase := impl.NewLikesUseCaseImpl(likesRepo)

	return APIServer{
		conf:           conf,
		authHandler:    handlers.CreateAuthHandler(authUseCase),
		profileHandler: handlers.CreateProfileHandler(profileUseCase),
		jwtHandler:     handlers.CreateCSRFHandler(*jwt, conf.JwtConf.TimeToLife),
		photosHandler:  handlers.CreatePhotosHandler(photosUseCase, photosService),
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
	UrlCSRF := serv.conf.ApiPathPrefix + UrlCSRFPostfix
	UrlPhotos := serv.conf.ApiPathPrefix + UrlPhotosPostfix
	UrlPhotosId := serv.conf.ApiPathPrefix + UrlPhotosIdPostfix
	UrlPhotosIdParams := serv.conf.ApiPathPrefix + UrlPhotosIdParamsPostfix
	UrlProfilePhotos := serv.conf.ApiPathPrefix + UrlProfilePhotosPostfix
	UrlProfilePhotosAvatar := serv.conf.ApiPathPrefix + UrlProfilePhotosAvatarPostfix
	UrlLikes := serv.conf.ApiPathPrefix + UrlLikesPostfix
	//main multiplexor
	multiplexor := mux.NewRouter()
	//multiplexor.Use()
	multiplexor.Use(serv.middlewares.AccessLogMiddleware)
	multiplexor.Use(serv.middlewares.PanicMiddleware)
	multiplexor.Use(serv.middlewares.CorsMiddleware)
	multiplexor.Methods(http.MethodOptions).HandlerFunc(Cors)

	multiplexor.HandleFunc(UrlUsers, serv.authHandler.GET).Methods(http.MethodGet)
	multiplexor.HandleFunc(UrlUsers, serv.authHandler.PUT).Methods(http.MethodPut)
	multiplexor.HandleFunc(UrlUsers, serv.authHandler.POST).Methods(http.MethodPost)

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
	multiplexorWithAuth.HandleFunc(UrlPhotos, serv.photosHandler.POST).Methods(http.MethodPost)
	multiplexorWithAuth.HandleFunc(UrlPhotosId, serv.photosHandler.GETPhoto).Methods(http.MethodGet)
	multiplexorWithAuth.HandleFunc(UrlPhotosId, serv.photosHandler.POSTPhoto).Methods(http.MethodPost)
	multiplexorWithAuth.HandleFunc(UrlPhotosId, serv.photosHandler.DELETE).Methods(http.MethodDelete)
	multiplexorWithAuth.HandleFunc(UrlPhotosIdParams, serv.photosHandler.PUTParams).Methods(http.MethodPut)
	multiplexorWithAuth.HandleFunc(UrlPhotosIdParams, serv.photosHandler.GETParams).Methods(http.MethodGet)
	multiplexorWithAuth.HandleFunc(UrlProfilePhotos, serv.photosHandler.GETAll).Methods(http.MethodGet)
	multiplexorWithAuth.HandleFunc(UrlProfilePhotosAvatar, serv.photosHandler.GETAvatar).Methods(http.MethodGet)
	multiplexorWithAuth.HandleFunc(UrlProfilePhotosAvatar, serv.photosHandler.PUTAvatar).Methods(http.MethodPut)

	multiplexorWithAuth.HandleFunc(UrlLikes, serv.likesHandler.Set).Methods(http.MethodPost)
	multiplexorWithAuth.HandleFunc(UrlLikes, serv.likesHandler.Get).Methods(http.MethodGet)

	multiplexorWithAuth.HandleFunc(UrlProfileId, serv.profileHandler.GetProfileHandler).Methods(http.MethodGet)
	multiplexorWithAuth.HandleFunc(UrlProfileIdShort, serv.profileHandler.GetShortProfileHandler).Methods(http.MethodGet) ///дописать
	multiplexorWithAuth.HandleFunc(UrlProfileId, serv.profileHandler.ChangeProfileHandler).Methods(http.MethodPut)        //свой профиль

	multiplexorWithAuth.HandleFunc(UrlCSRF, serv.jwtHandler.GetCSRF).Methods(http.MethodGet)

	//csrf multiplexor
	multiplexorCsrf := multiplexorWithAuth.PathPrefix("").Subrouter()
	multiplexorCsrf.Use(serv.middlewares.CSRFMiddleware)
	multiplexorCsrf.HandleFunc(UrlLikes, serv.likesHandler.Set).Methods(http.MethodPost)
	multiplexorCsrf.HandleFunc(UrlProfileId, serv.profileHandler.ChangeProfileHandler).Methods(http.MethodPut)
	multiplexorCsrf.HandleFunc(UrlUsers, serv.authHandler.DELETE).Methods(http.MethodDelete)

	serverAddr := serv.conf.ServerAddr + ":" + serv.conf.ServerPort
	server := http.Server{Addr: serverAddr, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, Handler: multiplexor}
	return server.ListenAndServe()
}
