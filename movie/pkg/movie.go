package model

import model "movie-rating-app/metadata/pkg"

type MovieDetails struct {
	Rating   float64        `json:"rating,omitempty"`
	Metadata model.Metadata `json:"metadata"`
}
