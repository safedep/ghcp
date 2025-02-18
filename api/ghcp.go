package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"buf.build/gen/go/safedep/api/connectrpc/go/safedep/services/ghcp/v1/ghcpv1connect"
	ghcpv1 "buf.build/gen/go/safedep/api/protocolbuffers/go/safedep/services/ghcp/v1"
	"connectrpc.com/connect"
	"github.com/safedep/dry/log"
	"github.com/safedep/ghcp/services"
)

type serviceSignature = services.Service[*ghcpv1.CreatePullRequestCommentRequest,
	*ghcpv1.CreatePullRequestCommentResponse]

type ghcpServiceHandler struct {
	ghcpv1connect.UnimplementedGitHubCommentsProxyServiceHandler

	ghcpService serviceSignature
}

var _ Handler = &ghcpServiceHandler{}

func NewGhcpServiceHandler(ghcpService serviceSignature) (*ghcpServiceHandler, error) {
	return &ghcpServiceHandler{
		ghcpService: ghcpService,
	}, nil
}

func (h *ghcpServiceHandler) Name() string {
	return "GitHub Comments Proxy Handler"
}

func (h *ghcpServiceHandler) Build(opts ...connect.HandlerOption) (string, http.Handler, error) {
	path, handler := ghcpv1connect.NewGitHubCommentsProxyServiceHandler(h, opts...)
	return path, handler, nil
}

func (h *ghcpServiceHandler) CreatePullRequestComment(ctx context.Context,
	req *connect.Request[ghcpv1.CreatePullRequestCommentRequest]) (*connect.Response[ghcpv1.CreatePullRequestCommentResponse], error) {
	log.Debugf("CreatePullRequestComment request received: %v", req.Msg)
	if req.Msg == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument,
			errors.New("request message is nil"))
	}

	res, err := h.ghcpService.Execute(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GHCP service: %w", err)
	}

	return connect.NewResponse(res), nil
}
