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
	ProfilesDbTableName        string
	InterestsDbTableName       string
	StaticInterestsDbTableName string
	AvatarDbTableName          string
	LikesDbTableName           string
	FiltersDbTableName         string
}

type ProfileServerConf struct {
	PostgresConf PostgresConnectConf
	ServerPort   string
	ServerAddr   string
}

var ApiServerDefaultConf = ProfileServerConf{
	ServerPort: "8080",
	ServerAddr: "0.0.0.0",
	PostgresConf: PostgresConnectConf{
		Username:                   "postgres",
		Password:                   "postgres",
		Address:                    "localhost",
		Port:                       "5432",
		DbName:                     "postgres",
		Params:                     "",
		ProfilesDbTableName:        "os_profiles",
		InterestsDbTableName:       "os_interests",
		StaticInterestsDbTableName: "os_staticInterests",
		AvatarDbTableName:          "os_avatars",
		LikesDbTableName:           "os_likes",
		FiltersDbTableName:         "os_filters",
	},
}

func (conf *ProfileServerConf) ReadFromFile(filename string) error {
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

func (conf ProfileServerConf) WriteBasicConfFile(filename string) error {
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

func (conf *ProfileServerConf) ProcessConfiguration(filename string) bool {
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
