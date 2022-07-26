package grpc

import (
	"context"
	gp "team3-task/api/grpc/gen/proto"
	"team3-task/internal/entity"
	"team3-task/internal/errors"
	"team3-task/internal/utils"
	"team3-task/pkg/logging"

	"google.golang.org/grpc"
)

type GrpcClient struct {
	client gp.AuthApiClient
	log    *logging.ZeroLogger
}

func NewGrpcClient(conn *grpc.ClientConn, log *logging.ZeroLogger) *GrpcClient {
	client := gp.NewAuthApiClient(conn)
	grpcClient := GrpcClient{client: client, log: log}
	return &grpcClient
}

func (c *GrpcClient) CheckAccess(authRequest *entity.AuthRequest) (entity.AuthResponse, error) {
	authResponse := entity.AuthResponse{}
	resp, err := c.client.Authenticate(context.Background(), &gp.AuthRequest{AccessToken: authRequest.AccessToken})
	if err != nil {
		return authResponse, err
	}
	err = utils.CheckEmail(resp.Username)
	if err != nil {
		return authResponse, errors.Newf("email is incorrect, please, relogin")
	}
	authResponse.Username = resp.Username
	authResponse.Error = resp.Error
	return authResponse, nil
}
