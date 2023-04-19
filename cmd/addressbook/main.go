package main

import (
	"context"
	"fmt"
	"infoblox-golang/cmd/addressbook/config"
	"infoblox-golang/cmd/addressbook/transport/grpc"
	"infoblox-golang/internal/addressbook"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	// Read config
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup persistence layer
	gorm := connectToDB(cfg.Db.Address)
	gorm.WithContext(ctx)

	// Setup services
	repository := addressbook.NewRepository(gorm)
	service := addressbook.NewService(repository)

	// Start HTTP server (and proxy calls to gRPC server endpoint)

	grpcSrv := grpc.NewServer(service)
	go func() {
		log.Println("Starting grpc gateway on port", cfg.Server.GrpcGatewayServerPort)
		if err := grpcSrv.RunGatewayServer(ctx, cfg.Server.GrpcGatewayServerPort); err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
}

func connectToDB(dbURL string) *gorm.DB {
	log.Println("Connecting to db")
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&addressbook.User{})

	return db
}
