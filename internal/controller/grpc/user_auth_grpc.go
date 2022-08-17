package grpc

import (
	"context"
	gp "team3-task/api/grpc/gen/proto"
	"team3-task/internal/entity"
	"team3-task/internal/errors"
	"team3-task/internal/utils"

	"google.golang.org/grpc"
)

type GClient struct {
	client gp.AuthApiClient
}

func NewGrpcClient(conn *grpc.ClientConn) *GClient {
	client := gp.NewAuthApiClient(conn)
	grpcClient := GClient{client: client}

	return &grpcClient
}

func (c *GClient) CheckAccess(authRequest *entity.AuthRequest) (entity.AuthResponse, error) {
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
