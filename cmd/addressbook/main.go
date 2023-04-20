package main

import (
	"context"
	"fmt"
	"infoblox-golang/cmd/addressbook/config"
	"infoblox-golang/cmd/addressbook/transport/grpc"
	"infoblox-golang/internal/addressbook"
	"infoblox-golang/internal/platform/storage"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// Read config
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	// Setup persistence layer
	db := storage.NewInMemory()

	// Setup services
	repository := addressbook.NewRepository(db)
	service := addressbook.NewService(repository)

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
