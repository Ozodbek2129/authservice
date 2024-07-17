package main

import (
	"auth/config"
	pb "auth/genproto/AuthService"
	"auth/service"
	"auth/storage/postgres"
	"fmt"
	"log"
	"net"
	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", config.Load().USER_SERVICE)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer listener.Close()

	db, err := postgres.ConnectDB()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// go func() {
	// 	redisClient := redis.NewClient(&redis.Options{
	// 		Addr: "localhost:6379",
	// 	})

	// 	fmt.Printf("Server is listening on port %s\n",config.Load().USER_ROUTER)
	// 	router := api.NewRouter(db, redisClient)
	// 	router.Run(config.Load().USER_ROUTER)
	// }()

	userservice := service.NewUserService(db)
	server := grpc.NewServer()

	pb.RegisterAuthUserServiceServer(server, userservice)

	fmt.Printf("Server is listening on port %s\n", config.Load().USER_SERVICE)
	if err = server.Serve(listener); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}
