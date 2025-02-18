package api

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	protovalidate "github.com/bufbuild/protovalidate-go"
	"github.com/safedep/dry/log"
	"google.golang.org/protobuf/proto"
)

type validatorInterceptor struct {
	validator protovalidate.Validator
}

func NewValidatorInterceptor() (connect.Interceptor, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, err
	}

	return &validatorInterceptor{validator: validator}, nil
}

func (v *validatorInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		log.Debugf("Validator Interceptor: Validating request: %s", req.Spec().Procedure)
		msg, ok := req.Any().(proto.Message)
		if !ok {
			return nil, connect.NewError(connect.CodeInvalidArgument,
				errors.New("request is not a proto message"))
		}

		if err := v.validator.Validate(msg); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		return next(ctx, req)
	}
}

func (v *validatorInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		return nil
	}
}

func (v *validatorInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, stream connect.StreamingHandlerConn) error {
		return fmt.Errorf("not implemented")
	}
}
