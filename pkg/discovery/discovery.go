package discovery

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Registry interface {
	Register(ctx context.Context, instID string, serviceName string, hostPort string) error

	DeRegister(ctx context.Context, instID string, serviceName string) error

	ServiceAddresses(ctx context.Context, serviceId string) ([]string, error)

	ReportHealth(instID string, serviceName string) error
}

var ErrorNotFound = errors.New("no services found")

func GenerateInstanceId(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName,
		rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
