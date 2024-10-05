package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	model "movie-rating-app/metadata/pkg"
	"movie-rating-app/metadata/pkg/constants"
	"movie-rating-app/movie/internal/gateway"
	"movie-rating-app/pkg/discovery"
	"net/http"

	"golang.org/x/exp/rand"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, constants.SERVICE_NAME)
	if err != nil {
		return nil, err
	}

	url := "http://" + addrs[rand.Intn(len(addrs))] + "/metadata"
	log.Printf("Calling metadata service. Request: GET " + url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	values := req.URL.Query()
	values.Add("id", id)

	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non 2xx response: %v", resp)
	}

	var m model.Metadata

	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	return &m, nil
}
