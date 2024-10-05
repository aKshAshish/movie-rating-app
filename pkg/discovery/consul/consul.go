package consul

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	consul "github.com/hashicorp/consul/api"
)

type Registry struct {
	client *consul.Client
}

func New(addr string) (*Registry, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Registry{client: client}, nil
}

func (r *Registry) Register(ctx context.Context, instID, serviceName, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("hostPort must be in a form of <host>:<port>, example -> localhost:9000")
	}
	port, err := strconv.Atoi(parts[1])

	if err != nil {
		return err
	}

	return r.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		Port:    port,
		Address: parts[0],
		ID:      instID,
		Name:    serviceName,
		Check:   &consul.AgentServiceCheck{CheckID: instID, TTL: "5s"},
	})
}

func (r *Registry) DeRegister(ctx context.Context, instID, _ string) error {
	return r.client.Agent().ServiceDeregister(instID)
}

func (r *Registry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}
	var res []string

	for _, e := range entries {
		res = append(res, fmt.Sprintf("%s:%d", e.Service.Address, e.Service.Port))
	}

	return res, nil
}

func (r *Registry) ReportHealth(instID, _ string) error {
	return r.client.Agent().PassTTL(instID, "")
}
