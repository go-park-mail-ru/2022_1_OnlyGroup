package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

const defaultWriteFilePerm = 0644

type PostgresConnectConf struct {
	Username string
	Password string
	Address  string
	Port     string
	DbName   string
	Params   string
}

type RedisConnectConf struct {
	Username string
	Password string
	Address  string
	Port     string
	DbNum    int
}

type APIServerConf struct {
	RedisConf        RedisConnectConf
	PostgresConf     PostgresConnectConf
	ServerPort       string
	ServerAddr       string
	ApiPathPrefix    string
	UsersDbTable     string
	ProfilesDbTable  string
	InterestsDbTable string
}

var ApiServerDefaultConf = APIServerConf{
	ServerPort:       "8080",
	ServerAddr:       "0.0.0.0",
	ApiPathPrefix:    "",
	UsersDbTable:     "os_users",
	ProfilesDbTable:  "os_profiles",
	InterestsDbTable: "os_interests",
	RedisConf: RedisConnectConf{
		Username: "",
		Password: "",
		Address:  "localhost",
		Port:     "6379",
		DbNum:    1,
	},
	PostgresConf: PostgresConnectConf{
		Username: "postgres",
		Password: "postgres",
		Address:  "localhost",
		Port:     "5432",
		DbName:   "postgres",
		Params:   "",
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
