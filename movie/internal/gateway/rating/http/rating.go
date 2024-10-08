package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"movie-rating-app/movie/internal/gateway"
	"movie-rating-app/pkg/discovery"
	model "movie-rating-app/rating/pkg"
	"movie-rating-app/rating/pkg/constants"
	"net/http"

	"golang.org/x/exp/rand"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

func (g *Gateway) GetAggregatedRating(ctx context.Context, recordId model.RecordID, recordType model.RecordType) (float64, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, constants.SERVICE_NAME)
	if err != nil {
		return 0, err
	}

	url := "http://" + addrs[rand.Intn(len(addrs))] + "/rating"
	log.Printf("Calling metadata service. Request: GET " + url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordId))
	values.Add("type", string(recordType))
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return 0, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("non 2xx response: %v", resp)
	}

	var v float64
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (g *Gateway) PutRating(ctx context.Context, rating *model.Rating) error {
	addrs, err := g.registry.ServiceAddresses(ctx, constants.SERVICE_NAME)
	if err != nil {
		return err
	}

	url := "http://" + addrs[rand.Intn(len(addrs))] + "/rating"
	log.Printf("Calling metadata service. Request: PUT " + url)

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", rating.RecordID)
	values.Add("type", rating.RecordType)
	values.Add("userId", string(rating.UserID))
	values.Add("value", fmt.Sprintf("%v", rating.Value))

	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("non 2xx response: %v", resp)
	}

	return nil
}
