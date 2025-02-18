package api

import (
	"net/http"

	"buf.build/gen/go/safedep/api/connectrpc/go/safedep/services/ghcp/v1/ghcpv1connect"
	"connectrpc.com/connect"
)

type ghcpServiceHandler struct {
	ghcpv1connect.UnimplementedGitHubCommentsProxyServiceHandler
}

var _ Handler = &ghcpServiceHandler{}

func NewGhcpServiceHandler() (*ghcpServiceHandler, error) {
	return &ghcpServiceHandler{}, nil
}

func (h *ghcpServiceHandler) Name() string {
	return "GitHub Comments Proxy Handler"
}

func (h *ghcpServiceHandler) Build(opts ...connect.HandlerOption) (string, http.Handler, error) {
	path, handler := ghcpv1connect.NewGitHubCommentsProxyServiceHandler(h, opts...)
	return path, handler, nil
}
