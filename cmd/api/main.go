package main

import (
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	var conf APIServerConf
	exit := conf.ProcessConfiguration("os_api_server.conf")
	if exit {
		return
	}
	serv, err := NewServer(conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	serv.Run()
}
