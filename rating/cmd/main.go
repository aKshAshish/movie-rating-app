package main

import (
	"log"
	"movie-rating-app/rating/internal/controller/rating"
	httphandler "movie-rating-app/rating/internal/handler/http"
	"movie-rating-app/rating/internal/repository/memory"
	"net/http"
)

func main() {
	log.Println("Starting rating service...")
	repo := memory.New()
	ctrl := rating.New(repo)
	h := httphandler.New(ctrl)

	http.Handle("/rating", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
