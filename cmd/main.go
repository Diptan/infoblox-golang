package main

import (
	"address-book/cmd/config"
	"address-book/cmd/transport/rest"
	"address-book/internal/addressbook"
	"address-book/pkg/storage"
	"fmt"
	"net/http"
	"time"
)

func main() {

	// Read config
	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Println(err)
	}

	// Setup persistence layer
	db := storage.NewInMemory()

	// Setup services
	repository := addressbook.NewRepository(db)
	service := addressbook.NewService(repository)

	controller := rest.NewController(service)

	// Register all endpoints
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      controller.RegisterHandlers(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
