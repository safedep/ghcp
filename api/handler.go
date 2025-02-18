package api

import (
	"net/http"

	"connectrpc.com/connect"
)

// Handler defines the contract for implementing API (gRPC)
// service handlers that can be registered with our H2 server
type Handler interface {
	// Build returns the path, handler and error for service registration with
	// H2 server. Handlers must include the options along with their own
	Build(...connect.HandlerOption) (string, http.Handler, error)

	// Name returns the name of the service.
	Name() string
}
