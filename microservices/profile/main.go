package main

import (
	"2022_1_OnlyGroup_back/app/handlers/grpcDelivery"
	"2022_1_OnlyGroup_back/app/repositories/postgres"
	"2022_1_OnlyGroup_back/app/usecases/implGrpc"
	"2022_1_OnlyGroup_back/microservices/profile/proto"
	_ "github.com/jackc/pgx"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	var conf APIServerConf
	exit := conf.ProcessConfiguration("os_api_server.conf")
	if exit {
		return
	}
	postgresSourceString := "postgresql://" + conf.PostgresConf.Username + ":" + conf.PostgresConf.Password +
		"@" + conf.PostgresConf.Address + ":" + conf.PostgresConf.Port + "/" + conf.PostgresConf.DbName + "?" + conf.PostgresConf.Params
	postgresConnect, err := sqlx.Open("pgx", postgresSourceString)
	if err != nil {
		log.Fatalf("failed connect dataBase, %v", err)
	}
	listener, err := net.Listen("tcp", "127.0.0.1:9111")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	profilesRepo, err := postgres.NewProfilePostgres(postgresConnect, conf.PostgresConf.ProfilesDbTableName, conf.PostgresConf.UsersDbTableName, conf.PostgresConf.InterestsDbTableName, conf.PostgresConf.StaticInterestsDbTableName, conf.PostgresConf.FiltersDbTableName, conf.PostgresConf.LikesDbTableName)
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
