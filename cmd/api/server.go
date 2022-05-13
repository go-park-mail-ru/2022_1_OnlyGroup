package main

import (
	http2 "2022_1_OnlyGroup_back/app/handlers/http"
	"2022_1_OnlyGroup_back/app/repositories"
	profileService "2022_1_OnlyGroup_back/app/repositories/grpc"
	"2022_1_OnlyGroup_back/app/repositories/postgres"
	redis_repo "2022_1_OnlyGroup_back/app/repositories/redis"
	_ "2022_1_OnlyGroup_back/app/usecases"
	"2022_1_OnlyGroup_back/app/usecases/impl"
	"2022_1_OnlyGroup_back/microservices/profile/proto"
	csrf "2022_1_OnlyGroup_back/pkg/csrf/impl"
	"2022_1_OnlyGroup_back/pkg/dataValidator"
	fileService "2022_1_OnlyGroup_back/pkg/fileService/impl"
	randomGenerator "2022_1_OnlyGroup_back/pkg/randomGenerator/impl"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
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
const UrlProfilePhotosPostfix = "/profiles/{id:[0-9]+}/photos"
const UrlProfilePhotosAvatarPostfix = "/profiles/{id:[0-9]+}/photos/avatar"

const UrlLikesPostfix = "/likes"
const UrlInterestsPostfix = "/interests"
const UrlInterestsStrParamsPostfix = "/interests/{str}"

const UrlFiltersPostfix = "/filters"

type APIServer struct {
	conf             APIServerConf
	authHandler      *http2.AuthHandler
	profileHandler   *http2.ProfileHandler
	photosHandler    *http2.PhotosHandler
	likesHandler     *http2.LikesHandler
	interestsHandler *http2.InterestsHandler
	middlewares      http2.Middlewares
	jwtHandler       *http2.CSRFHandler
	filtersHandler   *http2.FiltersHandler
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
	//profiles repositories
	var profilesRepo repositories.ProfileRepository
	if conf.ProfileServiceConf.Enable {
		connProfileService, err := grpc.Dial(conf.ProfileServiceConf.Address+":"+conf.ProfileServiceConf.Port, grpc.WithInsecure())
		if err != nil {
			return APIServer{}, err
		}
		client := proto.NewProfileRepositoryClient(connProfileService)
		profilesRepo = profileService.NewProfileGrpc(client)
	} else {
		profilesRepo, err = postgres.NewProfilePostgres(postgresConnect, conf.PostgresConf.ProfilesDbTableName, conf.PostgresConf.InterestsDbTableName, conf.PostgresConf.StaticInterestsDbTableName, conf.PostgresConf.FiltersDbTableName, conf.PostgresConf.LikesDbTableName)
		if err != nil {
			return APIServer{}, err
		}
	}

	//jwtToken
	jwt := csrf.NewJwtTokenGenerator("поменяй здесь генерацию", conf.CSRFConf.TimeToLife)
	//useCases
	photosRepo, err := postgres.NewPostgresPhotoRepository(postgresConnect, conf.PostgresConf.PhotosDbTableName, conf.PostgresConf.UsersDbTableName, conf.PostgresConf.AvatarDbTableName)
	if err != nil {
		return APIServer{}, err
	}

	sessionsRepo := redis_repo.NewRedisSessionRepository(redisConnect, conf.RedisConf.SessionsPrefix, randomGenerator.NewMathRandomGenerator())

	//set validators
	dataValidator.SetValidators()
	//useCases
	authUseCase := impl.NewAuthUseCaseImpl(usersRepo, sessionsRepo, profilesRepo)
	profileUseCase := impl.NewProfileUseCaseImpl(profilesRepo) //profilesRepo
	photosUseCase := impl.NewPhotosUseCase(photosRepo)

	//fileService
	photosService, err := fileService.NewFileServicePhotos(conf.PhotosStorageDirectory)
	if err != nil {
		return APIServer{}, err
	}

	return APIServer{
		conf:             conf,
		authHandler:      http2.CreateAuthHandler(authUseCase),
		profileHandler:   http2.CreateProfileHandler(profileUseCase),
		jwtHandler:       http2.CreateCSRFHandler(jwt),
		photosHandler:    http2.CreatePhotosHandler(photosUseCase, photosService),
		likesHandler:     http2.CreateLikesHandler(profileUseCase),
		interestsHandler: http2.CreateInterestsHandler(profileUseCase),
		middlewares:      http2.MiddlewaresImpl{AuthUseCase: authUseCase},
		filtersHandler:   http2.CreateFiltersHandler(profileUseCase),
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
	UrlFilters := serv.conf.ApiPathPrefix + UrlFiltersPostfix

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
	multiplexorWithAuth.HandleFunc(UrlInterestsStrParamsPostfix, serv.interestsHandler.GetDynamic).Methods(http.MethodGet)
	//likes
	multiplexorWithAuth.HandleFunc(UrlProfilePhotosAvatar, serv.photosHandler.PUTAvatar).Methods(http.MethodPut)
	multiplexorWithAuth.HandleFunc(UrlLikes, serv.likesHandler.Get).Methods(http.MethodGet)
	//profile
	multiplexorWithAuth.HandleFunc(UrlProfileId, serv.profileHandler.GetProfileHandler).Methods(http.MethodGet)
	multiplexorWithAuth.HandleFunc(UrlProfileIdShort, serv.profileHandler.GetShortProfileHandler).Methods(http.MethodGet) ///дописать
	multiplexorWithAuth.HandleFunc(UrlCSRF, serv.jwtHandler.PostCSRF).Methods(http.MethodPost)
	multiplexorWithCsrf.HandleFunc(UrlFilters, serv.filtersHandler.Get).Methods(http.MethodGet)

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
	//filters with CSRF
	multiplexorWithCsrf.HandleFunc(UrlFilters, serv.filtersHandler.Put).Methods(http.MethodPut)

	serverAddr := serv.conf.ServerAddr + ":" + serv.conf.ServerPort
	server := http.Server{Addr: serverAddr, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, Handler: multiplexor}
	return server.ListenAndServe()
}
