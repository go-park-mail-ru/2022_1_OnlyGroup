package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

const defaultWriteFilePerm = 0644

type PostgresConnectConf struct {
	Username                   string
	Password                   string
	Address                    string
	Port                       string
	DbName                     string
	Params                     string
	UsersDbTableName           string
	ProfilesDbTableName        string
	InterestsDbTableName       string
	StaticInterestsDbTableName string
	PhotosDbTableName          string
	AvatarDbTableName          string
	LikesDbTableName           string
	FiltersDbTableName         string
}

type RedisConnectConf struct {
	Username       string
	Password       string
	Address        string
	Port           string
	DbNum          int
	SessionsPrefix string
}

type JwtConf struct {
	Enable     bool
	TimeToLife int64
}

type ProfileMicroserviceConf struct {
	Enable  bool
	Address string
	Port    string
}

type APIServerConf struct {
	CSRFConf               JwtConf
	RedisConf              RedisConnectConf
	PostgresConf           PostgresConnectConf
	ProfileServiceConf     ProfileMicroserviceConf
	ServerPort             string
	ServerAddr             string
	ApiPathPrefix          string
	PhotosStorageDirectory string
}

var ApiServerDefaultConf = APIServerConf{
	ServerPort:             "8080",
	ServerAddr:             "0.0.0.0",
	ApiPathPrefix:          "",
	PhotosStorageDirectory: "./photos",
	RedisConf: RedisConnectConf{
		Username:       "",
		Password:       "",
		Address:        "localhost",
		Port:           "6379",
		DbNum:          1,
		SessionsPrefix: "os_sessions",
	},
	PostgresConf: PostgresConnectConf{
		Username:                   "postgres",
		Password:                   "postgres",
		Address:                    "localhost",
		Port:                       "5432",
		DbName:                     "postgres",
		Params:                     "",
		UsersDbTableName:           "os_users",
		ProfilesDbTableName:        "os_profiles",
		InterestsDbTableName:       "os_interests",
		StaticInterestsDbTableName: "os_staticInterests",
		PhotosDbTableName:          "os_photos",
		AvatarDbTableName:          "os_avatars",
		LikesDbTableName:           "os_likes",
		FiltersDbTableName:         "os_filters",
	},
	CSRFConf: JwtConf{
		Enable:     false,
		TimeToLife: 86400,
	},
	ProfileServiceConf: ProfileMicroserviceConf{
		Enable:  false,
		Port:    "9111",
		Address: "127.0.0.1",
	},
}

func (conf *APIServerConf) ReadFromFile(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("ReadFile failed: %w", err)
	}
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		return fmt.Errorf("yaml read failed: %w", err)
	}
	return nil
}

func (conf APIServerConf) WriteBasicConfFile(filename string) error {
	bytes, err := yaml.Marshal(ApiServerDefaultConf)
	if err != nil {
		return fmt.Errorf("marshalling default conf failed: %w", err)
	}

	err = os.WriteFile(filename, bytes, defaultWriteFilePerm)
	if err != nil {
		return fmt.Errorf("WriteFile default conf failed: %w", err)
	}

	return nil
}

func (conf *APIServerConf) ProcessConfiguration(filename string) bool {
	err := conf.ReadFromFile(filename)
	if err == nil {
		return false
	}
	fmt.Println(err)
	fmt.Println("Creating new configuration: " + filename + ".new")
	err = conf.WriteBasicConfFile(filename + ".new")
	if err != nil {
		fmt.Println(err)
		return true
	}
	fmt.Println("Configuration file successfully wrote")
	return true
}
