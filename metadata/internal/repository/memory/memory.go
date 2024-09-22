package memory

import (
	"context"
	"movie-rating-app/metadata/internal/repository"
	model "movie-rating-app/metadata/pkg"
	"sync"
)

type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

func New() *Repository {
	return &Repository{data: map[string]*model.Metadata{}}
}

// Get retrieves metadata from repository for a given movie id.
func (r *Repository) Get(_ context.Context, id string) (*model.Metadata, error) {
	r.Lock()
	defer r.Unlock()

	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}

	return m, nil
}

// Put adds metadata to repository for a given movie id.
func (r *Repository) Put(_ context.Context, id string, m *model.Metadata) error {
	r.Lock()
	defer r.Unlock()

	r.data[id] = m
	return nil
}
