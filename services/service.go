package services

import (
	"context"
)

type ServiceConfiguration struct{}

// Service defines the contract for implementing services
// that hold business logic.
type Service[R, V any] interface {
	Name() string
	Config() ServiceConfiguration
	Execute(context.Context, R) (V, error)
}
