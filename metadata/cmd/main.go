package main

import (
	"log"
	"movie-rating-app/metadata/internal/controller/metadata"
	httphandler "movie-rating-app/metadata/internal/handler/http"
	"movie-rating-app/metadata/internal/repository/memory"
	"net/http"
)

func main() {
	log.Println("Starting the metadata service...")
	repo := memory.New()
	ctrl := metadata.New(repo)
	h := httphandler.New(ctrl)

	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
