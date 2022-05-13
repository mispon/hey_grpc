package reflection

import (
	"context"
	"net"

	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// NewClient creates and returns client
func NewClient(ctx context.Context, hostPort string) (*grpcreflect.Client, error) {
	client, _, err := NewClientConn(ctx, hostPort)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// NewClientConn creates client and returns with conn
func NewClientConn(ctx context.Context, hostPort string) (*grpcreflect.Client, *grpc.ClientConn, error) {
	_, _, err := net.SplitHostPort(hostPort)
	if err != nil {
		return nil, nil, err
	}

	conn, err := grpc.DialContext(ctx, hostPort, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	client := grpcreflect.NewClient(ctx, reflectpb.NewServerReflectionClient(conn))
	return client, conn, nil
}
