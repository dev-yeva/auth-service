package client

import (
	"context"
	"errors"
	"time"

	"github.com/dev-yeva/auth_protos/gen"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	api gen.AuthClient
}

func MustNew(address string, timeout time.Duration, retryCount int) *Client {

	retryOpts := []retry.CallOption{
		retry.WithCodes(codes.Aborted, codes.Canceled, codes.DeadlineExceeded),
		retry.WithMax(uint(retryCount)),
		retry.WithPerRetryTimeout(timeout),
	}

	connection, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(retry.UnaryClientInterceptor(retryOpts...)),
	)
	if err != nil {
		panic(err)
	}
	return &Client{api: gen.NewAuthClient(connection)}
}

func (c *Client) Register(ctx context.Context, email, password string) (int64, error) {
	resp, err := c.api.Register(ctx, &gen.RegisterRequest{Email: email, Password: password})
	if err != nil {
		return 0, FormatError(err)
	}
	return resp.UserId, nil
}

func (c *Client) Login(ctx context.Context, email, password string, appId int32) (string, error) {
	resp, err := c.api.Login(ctx, &gen.LoginRequest{Email: email, Password: password, AppId: appId})
	if err != nil {
		return "", FormatError(err)
	}
	return resp.Token, nil
}

func (c *Client) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	resp, err := c.api.IsAdmin(ctx, &gen.IsAdminRequest{UserId: userId})
	if err != nil {
		return false, FormatError(err)
	}
	return resp.IsAdmin, nil
}

func FormatError(err error) error {
	// st - структура с ошибкой от gRPC. содержит статус код, сообщение, детали
	if st, ok := status.FromError(err); ok {
		return errors.New(st.Message())
	}
	return err
}
