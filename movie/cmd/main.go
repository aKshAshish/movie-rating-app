package main

import (
	"log"
	"movie-rating-app/movie/internal/controller/movie"
	metadatgateway "movie-rating-app/movie/internal/gateway/metadata/http"
	ratinggateway "movie-rating-app/movie/internal/gateway/rating/http"
	httphandler "movie-rating-app/movie/internal/handler/http"
	"net/http"
)

func main() {
	log.Println("Starting movie service...")
	ratingGateway := ratinggateway.New("http://localhost:8082")
	metadataGateway := metadatgateway.New("http://localhost:8081")
	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httphandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
