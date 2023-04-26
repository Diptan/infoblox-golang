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
	"time"

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
	gorm := connectToDB(cfg.Db.Host)
	gorm.WithContext(ctx)

	// Setup services
	repository := addressbook.NewRepository(gorm)
	service := addressbook.NewService(repository)

	// Start HTTP server (and proxy calls to gRPC server endpoint)

	grpcSrv := grpc.NewServer(service)
	go func() {
		log.Println("Starting grpc gateway on port", cfg.Server.GatewayAddr)
		if err := grpcSrv.RunGatewayServer(ctx, cfg.Server.GatewayAddr); err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
}

var counts int64

func openDb(dbURL string) (*gorm.DB, error) {
	log.Println("Connecting to db")
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&addressbook.User{})

	return db, nil
}

func connectToDB(dbURL string) *gorm.DB {
	for {
		connection, err := openDb(dbURL)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}
