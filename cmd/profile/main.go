package main

import (
	"2022_1_OnlyGroup_back/app/handlers/grpcDelivery"
	"2022_1_OnlyGroup_back/app/repositories/postgres"
	"2022_1_OnlyGroup_back/microservices/profile/implGrpc"
	"2022_1_OnlyGroup_back/microservices/profile/proto"
	_ "github.com/jackc/pgx"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"log"
	"net"
)

//microserviceAddress
//microservicePort
//dbUserName
//dbPassword
//dbAddress
//dbPort
//dbName
//dbParams
//ProfilesDbTableName
//InterestsDbTableName
//StaticInterestsDbTableName
//FiltersDbTableName
//LikesDbTableName

func main() {
	//var conf ProfileServerConf
	//exit := conf.ProcessConfiguration("os_profile_server.conf")
	//if exit {
	//	return
	//}
	//postgresSourceString := "postgresql://" + os.Getenv("dbUserName") + ":" + os.Getenv("dbPassword") +
	//	"@" + os.Getenv("dbAddress") + ":" + os.Getenv("dbPort") + "/" + os.Getenv("dbName") + "?" + os.Getenv("dbParams")
	//postgresConnect, err := sqlx.Open("pgx", postgresSourceString)
	//if err != nil {
	//	log.Fatalf("failed connect dataBase, %v", err)
	//}
	//listener, err := net.Listen("tcp", os.Getenv("microserviceAddress")+":"+os.Getenv("microservicePort"))
	//if err != nil {
	//	log.Fatalf("failed to listen: %v", err)
	//}
	//profilesRepo, err := postgres.NewProfilePostgres(postgresConnect, os.Getenv("ProfilesDbTableName"), os.Getenv("InterestsDbTableName"),
	//	os.Getenv("StaticInterestsDbTableName"), os.Getenv("FiltersDbTableName"),
	//	os.Getenv("LikesDbTableName"))
	//if err != nil {
	//	log.Fatalf("failed connect dataBase, %v", err)
	//}
	//profileUseCase := implGrpc.NewProfileUseCaseGRPCDelivery(profilesRepo)
	//grpcServer := grpc.NewServer()
	//proto.RegisterProfileRepositoryServer(grpcServer, grpcDelivery.NewProfileHandler(profileUseCase))
	//if err := grpcServer.Serve(listener); err != nil {
	//	log.Fatal(err)
	//}

	var conf ProfileServerConf
	exit := conf.ProcessConfiguration("os_profile_server.conf")
	if exit {
		return
	}
	postgresSourceString := "postgresql://" + conf.PostgresConf.Username + ":" + conf.PostgresConf.Password +
		"@" + conf.PostgresConf.Address + ":" + conf.PostgresConf.Port + "/" + conf.PostgresConf.DbName + "?" + conf.PostgresConf.Params
	postgresConnect, err := sqlx.Open("pgx", postgresSourceString)
	if err != nil {
		log.Fatalf("failed connect dataBase, %v", err)
	}
	listener, err := net.Listen("tcp", conf.ServerAddr+":"+conf.ServerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	profilesRepo, err := postgres.NewProfilePostgres(postgresConnect, conf.PostgresConf.ProfilesDbTableName, conf.PostgresConf.InterestsDbTableName,
		conf.PostgresConf.StaticInterestsDbTableName, conf.PostgresConf.FiltersDbTableName, conf.PostgresConf.LikesDbTableName)
	if err != nil {
		log.Fatalf("failed connect dataBase, %v", err)
	}
	profileUseCase := implGrpc.NewProfileUseCaseGRPCDelivery(profilesRepo)
	grpcServer := grpc.NewServer()
	proto.RegisterProfileRepositoryServer(grpcServer, grpcDelivery.NewProfileHandler(profileUseCase))
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

