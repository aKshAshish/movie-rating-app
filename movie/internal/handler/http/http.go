package http

import (
	"encoding/json"
	"errors"
	"log"
	"movie-rating-app/movie/internal/controller/movie"
	"net/http"
)

type Handler struct {
	ctrl *movie.Controller
}

func New(ctrl *movie.Controller) *Handler {
	return &Handler{ctrl}
}

func (h *Handler) GetMovieDetails(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	resp, err := h.ctrl.Get(req.Context(), id)
	if err != nil && errors.Is(err, movie.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Repository get error %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Response encode error %v", err)
	}
}
