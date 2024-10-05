package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"movie-rating-app/movie/internal/controller/movie"
	metadatgateway "movie-rating-app/movie/internal/gateway/metadata/http"
	ratinggateway "movie-rating-app/movie/internal/gateway/rating/http"
	httphandler "movie-rating-app/movie/internal/handler/http"
	"movie-rating-app/movie/pkg/constants"
	"movie-rating-app/pkg/discovery"
	"movie-rating-app/pkg/discovery/consul"
	"net/http"
	"time"
)

func getPort() int {
	var port int
	flag.IntVar(&port, "port", 9001, "Server Port")
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
	log.Printf("Starting the movie service at port %d", port)

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

	ratingGateway := ratinggateway.New(registry)
	metadataGateway := metadatgateway.New(registry)

	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httphandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
