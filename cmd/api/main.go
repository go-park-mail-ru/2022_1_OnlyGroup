package main

//import res"2022_1_OnlyGroup_back/app/repositories/redis"

func main() {
	serv := NewServer(":8080")
	serv.Run()
}

//func main() {
//	client := redis.NewClient(&redis.Options{})
//	result, err := client.Ping(context.Background()).Result()
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(result)
//
//	repo := redisRepo.CreateRedisSessionRepository(client, "test")
//
//	session, err := repo.AddSession(5, "test_additional_data")
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(session)
//
//	id, info, err := repo.GetIdBySession(session)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(id)
//	fmt.Println(info)
//
//	err = repo.RemoveSession(session)
//
//	id, info, err = repo.GetIdBySession(session)
//	if err != nil {
//		fmt.Println(err)
//	}
//}

//func main() {
//	connConf, err := pgx.ParseConnectionString("postgresql://ilya:il28102001@localhost/ilya")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	conn, err := pgx.Connect(connConf)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	repo, err := postgres.CreatePostgresUsersRepo(conn, "testTable")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	user, err := repo.AddUser("test@mail.ru", "test_pass")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(user)
//
//	authorize, err := repo.Authorize("test@mail.ru", "test_pass")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(authorize)
//
//	a, err := repo.Authorize("test@mail.ru", "test_peass")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(a)
//
//}
