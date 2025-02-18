package services

import (
	"context"
	"fmt"
)

type ServiceConfiguration struct{}

// Service defines the contract for implementing services
// that hold business logic.
type Service[R, V any] interface {
	Name() string
	Config() ServiceConfiguration
	Execute(context.Context, R) (V, error)
}

// A base service for any services to build upon. Helps in keeping things DRY
// and extensibility of the interface with backward compatibility.
type unimplementedService[R, V any] struct{}

func (us *unimplementedService[R, V]) Name() string {
	return "UnimplementedService"
}

func (us *unimplementedService[R, V]) Config() ServiceConfiguration {
	return ServiceConfiguration{}
}

func (us *unimplementedService[R, V]) Execute(ctx context.Context, request R) (V, error) {
	var nilValue V
	return nilValue, fmt.Errorf("unimplemented service")
}
