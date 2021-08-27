package grpc_bugsnag

import (
	"context"
	"fmt"

	"github.com/bugsnag/bugsnag-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		bugsnag.OnBeforeNotify(
			func(event *bugsnag.Event, config *bugsnag.Configuration) error {
				event.MetaData.AddStruct("RPC Method", info.FullMethod)
				event.MetaData.AddStruct("RPC Message", fmt.Sprintf("%+v", req))
				return nil
			},
		)

		h, err := handler(ctx, req)
		if err != nil {
			return h, errors.WithStack(err)
		}

		return h, err
	}
}
