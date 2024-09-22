package movie

import (
	"context"
	"errors"
	metadatamodel "movie-rating-app/metadata/pkg"
	"movie-rating-app/movie/internal/gateway"
	model "movie-rating-app/movie/pkg"
	ratingmodel "movie-rating-app/rating/pkg"
)

var ErrNotFound = errors.New("not found")

type ratingGateway interface {
	GetAggregatedRating(ctx context.Context, recordId ratingmodel.RecordID, reordType ratingmodel.RecordType) (float64, error)
	PutRating(ctx context.Context, rating *ratingmodel.Rating) error
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadatamodel.Metadata, error)
}

type Controller struct {
	rating   ratingGateway
	metadata metadataGateway
}

func New(r ratingGateway, m metadataGateway) *Controller {
	return &Controller{rating: r, metadata: m}
}

func (c *Controller) Get(ctx context.Context, id string) (*model.MovieDetails, error) {
	metadata, err := c.metadata.Get(ctx, id)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	details := &model.MovieDetails{Metadata: *metadata}
	rating, err := c.rating.GetAggregatedRating(ctx, ratingmodel.RecordID(id), ratingmodel.RecordTypeMovie)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
	} else if err != nil {
		return nil, err
	} else {
		details.Rating = rating
	}
	return details, nil
}
