package memory

import (
	"context"
	"errors"
	"movie-rating-app/pkg/discovery"
	"sync"
	"time"
)

type Registry struct {
	sync.RWMutex
	serviceAddrs map[string]map[string]*serviceInst
}

type serviceInst struct {
	hostPort   string
	lastActive time.Time
}

func New() *Registry {
	return &Registry{serviceAddrs: map[string]map[string]*serviceInst{}}
}

func (r *Registry) Register(ctx context.Context, instID, serviceName, hostPort string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName]; !ok {
		r.serviceAddrs[serviceName] = map[string]*serviceInst{}
	}
	r.serviceAddrs[serviceName][instID] = &serviceInst{
		hostPort:   hostPort,
		lastActive: time.Now(),
	}
	return nil
}

func (r *Registry) DeRegister(ctx context.Context, instID, serviceName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName]; !ok {
		return nil
	}

	delete(r.serviceAddrs[serviceName], instID)

	return nil
}

func (r *Registry) ReportHealth(instID, serviceName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName]; !ok {
		return errors.New("service is not registered yet")
	}

	if _, ok := r.serviceAddrs[serviceName][instID]; !ok {
		return errors.New("service instance is not registered yet")
	}

	r.serviceAddrs[serviceName][instID].lastActive = time.Now()

	return nil
}

func (r *Registry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	r.Lock()
	defer r.Unlock()

	var res []string

	if len(r.serviceAddrs[serviceName]) == 0 {
		return nil, discovery.ErrorNotFound
	}

	for _, inst := range r.serviceAddrs[serviceName] {
		if inst.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		res = append(res, inst.hostPort)
	}
	return res, nil
}
