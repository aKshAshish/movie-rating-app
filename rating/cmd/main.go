package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"movie-rating-app/pkg/discovery"
	"movie-rating-app/pkg/discovery/consul"
	"movie-rating-app/rating/internal/controller/rating"
	httphandler "movie-rating-app/rating/internal/handler/http"
	"movie-rating-app/rating/internal/repository/memory"
	"movie-rating-app/rating/pkg/constants"
	"net/http"
	"time"
)

func getPort() int {
	var port int
	flag.IntVar(&port, "port", 8801, "Server Port")
	flag.Parse()
	return port
}

func reportHealth(registry *consul.Registry, instID string) {
	for {
		if err := registry.ReportHealth(instID, constants.SERVICE_NAME); err != nil {
			log.Println("Failed to report health state, error: ", err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	port := getPort()
	log.Printf("Starting the rating service at port %d", port)

	registry, err := consul.New("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instID := discovery.GenerateInstanceId(constants.SERVICE_NAME)
	err = registry.Register(ctx, instID, constants.SERVICE_NAME, fmt.Sprintf("localhost:%d", port))
	if err != nil {
		panic(err)
	}
	defer registry.DeRegister(ctx, instID, constants.SERVICE_NAME)

	go reportHealth(registry, instID)

	repo := memory.New()
	ctrl := rating.New(repo)
	h := httphandler.New(ctrl)

	http.Handle("/rating", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
