package main

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/repositories/postgres"
	redis_repo "2022_1_OnlyGroup_back/app/repositories/redis"
	_ "2022_1_OnlyGroup_back/app/usecases"
	"2022_1_OnlyGroup_back/app/usecases/impl"
	csrf "2022_1_OnlyGroup_back/pkg/csrf/impl"
	"2022_1_OnlyGroup_back/pkg/dataValidator"
	fileService "2022_1_OnlyGroup_back/pkg/fileService/impl"
	randomGenerator "2022_1_OnlyGroup_back/pkg/randomGenerator/impl"
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
const UrlInterestsPostfix = "/interests"

type APIServer struct {
	conf             APIServerConf
	authHandler      *handlers.AuthHandler
	profileHandler   *handlers.ProfileHandler
	photosHandler    *handlers.PhotosHandler
	likesHandler     *handlers.LikesHandler
	interestsHandler *handlers.InterestsHandler
	middlewares      handlers.Middlewares
	jwtHandler       *handlers.CSRFHandler
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
	usersRepo, err := postgres.NewPostgresUsersRepo(postgresConnect, conf.PostgresConf.UsersDbTableName, randomGenerator.NewCryptoRandomGenerator())
	if err != nil {
		return APIServer{}, err
	}
	profilesRepo, err := postgres.NewProfilePostgres(postgresConnect, conf.PostgresConf.ProfilesDbTableName, conf.PostgresConf.UsersDbTableName, conf.PostgresConf.InterestsDbTableName, conf.PostgresConf.StaticInterestsDbTableName)
	if err != nil {
		return APIServer{}, err
	}
	interestsRepo, err := postgres.NewInterestsPostgres(postgresConnect, conf.PostgresConf.StaticInterestsDbTableName)
	if err != nil {
		return APIServer{}, err
	}
	//jwtToken
	jwt := csrf.NewJwtTokenGenerator("поменяй здесь генерацию", conf.CSRFConf.TimeToLife)
	//useCases
	photosRepo, err := postgres.NewPostgresPhotoRepository(postgresConnect, conf.PostgresConf.PhotosDbTableName, conf.PostgresConf.UsersDbTableName, conf.PostgresConf.AvatarDbTableName)
	if err != nil {
		return APIServer{}, err
	}

	likesRepo, err := postgres.NewLikesPostgres(postgresConnect, conf.PostgresConf.LikesDbTableName, conf.PostgresConf.UsersDbTableName)
	if err != nil {
		return APIServer{}, err
	}
	sessionsRepo := redis_repo.NewRedisSessionRepository(redisConnect, conf.RedisConf.SessionsPrefix, randomGenerator.NewMathRandomGenerator())
	//set validators
	dataValidator.SetValidators()
	//useCases
	authUseCase := impl.NewAuthUseCaseImpl(usersRepo, sessionsRepo, profilesRepo)
	profileUseCase := impl.NewProfileUseCaseImpl(profilesRepo)
	photosUseCase := impl.NewPhotosUseCase(photosRepo)
	interestsUseCase := impl.NewInterestsUseCaseImpl(interestsRepo)

	//fileService
	photosService, err := fileService.NewFileServicePhotos(conf.PhotosStorageDirectory)
	if err != nil {
		return APIServer{}, err
	}
	likeUseCase := impl.NewLikesUseCaseImpl(likesRepo)

	return APIServer{
		conf:             conf,
		authHandler:      handlers.CreateAuthHandler(authUseCase),
		profileHandler:   handlers.CreateProfileHandler(profileUseCase),
		jwtHandler:       handlers.CreateCSRFHandler(jwt),
		photosHandler:    handlers.CreatePhotosHandler(photosUseCase, photosService),
		likesHandler:     handlers.CreateLikesHandler(likeUseCase),
		interestsHandler: handlers.CreateInterestsHandler(interestsUseCase),
		middlewares:      handlers.MiddlewaresImpl{AuthUseCase: authUseCase},
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
	UrlInterests := serv.conf.ApiPathPrefix + UrlInterestsPostfix

	//main multiplexor
	multiplexor := mux.NewRouter()
	//log middleware
	multiplexor.Use(serv.middlewares.AccessLogMiddleware)
	//panic middleware
	multiplexor.Use(serv.middlewares.PanicMiddleware)
	//cors middlewares
	multiplexor.Use(serv.middlewares.CorsMiddleware)
	//multiplexor with auth
	multiplexorWithAuth := multiplexor.PathPrefix("").Subrouter()
	//auth middleware
	multiplexorWithAuth.Use(serv.middlewares.CheckAuthMiddleware)
	//csrf multiplexor
	multiplexorWithCsrf := multiplexorWithAuth.PathPrefix("").Subrouter()
	//CSRF middleware
	if serv.conf.CSRFConf.Enable {
		multiplexorWithCsrf.Use(serv.middlewares.CSRFMiddleware)
	}
	//cors
	multiplexor.Methods(http.MethodOptions).HandlerFunc(Cors)
	//auth
	multiplexor.HandleFunc(UrlUsers, serv.authHandler.GET).Methods(http.MethodGet)
	multiplexor.HandleFunc(UrlUsers, serv.authHandler.PUT).Methods(http.MethodPut)
	multiplexor.HandleFunc(UrlUsers, serv.authHandler.POST).Methods(http.MethodPost)
	//profile methods
	multiplexorWithAuth.HandleFunc(UrlProfileId, serv.profileHandler.GetProfileHandler).Methods(http.MethodGet)
	multiplexorWithAuth.HandleFunc(UrlProfileIdShort, serv.profileHandler.GetShortProfileHandler).Methods(http.MethodGet) ///дописать
	//photos
	multiplexorWithAuth.HandleFunc(UrlPhotosId, serv.photosHandler.GETPhoto).Methods(http.MethodGet)
	multiplexorWithAuth.HandleFunc(UrlPhotosIdParams, serv.photosHandler.GETParams).Methods(http.MethodGet)
	multiplexorWithAuth.HandleFunc(UrlProfilePhotos, serv.photosHandler.GETAll).Methods(http.MethodGet)
	multiplexorWithAuth.HandleFunc(UrlProfilePhotosAvatar, serv.photosHandler.GETAvatar).Methods(http.MethodGet)
	//interests
	multiplexorWithAuth.HandleFunc(UrlInterests, serv.interestsHandler.Get).Methods(http.MethodGet)

	//likes
	multiplexorWithAuth.HandleFunc(UrlLikes, serv.likesHandler.Get).Methods(http.MethodGet)
	//profile
	multiplexorWithAuth.HandleFunc(UrlProfileId, serv.profileHandler.GetProfileHandler).Methods(http.MethodGet)
	multiplexorWithAuth.HandleFunc(UrlProfileIdShort, serv.profileHandler.GetShortProfileHandler).Methods(http.MethodGet) ///дописать
	multiplexorWithAuth.HandleFunc(UrlCSRF, serv.jwtHandler.PostCSRF).Methods(http.MethodPost)

	//photos with CSRF
	multiplexorWithCsrf.HandleFunc(UrlProfilePhotosAvatar, serv.photosHandler.PUTAvatar).Methods(http.MethodPut)
	//likes with CSRF
	multiplexorWithCsrf.HandleFunc(UrlLikes, serv.likesHandler.Set).Methods(http.MethodPost)
	//profile with CSRF
	multiplexorWithCsrf.HandleFunc(UrlProfileId, serv.profileHandler.ChangeProfileHandler).Methods(http.MethodPut)
	multiplexorWithCsrf.HandleFunc(UrlProfileCandidates, serv.profileHandler.GetCandidateHandler).Methods(http.MethodPost)
	//users with CSRF
	multiplexorWithCsrf.HandleFunc(UrlUsers, serv.authHandler.DELETE).Methods(http.MethodDelete)
	//photos with CSRF
	multiplexorWithCsrf.HandleFunc(UrlPhotos, serv.photosHandler.POST).Methods(http.MethodPost)
	multiplexorWithCsrf.HandleFunc(UrlPhotosId, serv.photosHandler.POSTPhoto).Methods(http.MethodPost)
	multiplexorWithCsrf.HandleFunc(UrlPhotosId, serv.photosHandler.DELETE).Methods(http.MethodDelete)
	multiplexorWithCsrf.HandleFunc(UrlPhotosIdParams, serv.photosHandler.PUTParams).Methods(http.MethodPut)
	multiplexorWithCsrf.HandleFunc(UrlProfilePhotosAvatar, serv.photosHandler.PUTAvatar).Methods(http.MethodPut)

	serverAddr := serv.conf.ServerAddr + ":" + serv.conf.ServerPort
	server := http.Server{Addr: serverAddr, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, Handler: multiplexor}
	return server.ListenAndServe()
}
