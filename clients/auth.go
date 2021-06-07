package clients

import (
	"context"

	"github.com/perfectbui/chat/models"
	"github.com/perfectbui/chat/pb"
	"github.com/perfectbui/chat/pb/dto"
	"github.com/perfectbui/chat/pb/types"

	"google.golang.org/grpc"
)

var AuthClient *authClient

type authClient struct {
	service pb.AuthServiceClient
}

func LoadAuthClient(conn *grpc.ClientConn) {
	authServiceClient := pb.NewAuthServiceClient(conn)
	AuthClient = &authClient{service: authServiceClient}
}

func NewAuthClient(svc pb.AuthServiceClient) *authClient {
	return &authClient{svc}
}

func (client *authClient) VerifyAccessToken(ctx context.Context, in *dto.VerifyTokenRequest, opts ...grpc.CallOption) (*dto.VerifyTokenResponse, error) {
	res, err := client.service.VerifyAccessToken(ctx, &dto.VerifyTokenRequest{Token: in.Token})
	if err != nil {
		return nil, err
	}
	return &dto.VerifyTokenResponse{UserID: res.UserID, Status: types.HttpStatus_Ok}, nil
}

func (client *authClient) VerifyRefreshToken(ctx context.Context, in *dto.VerifyTokenRequest, opts ...grpc.CallOption) (*dto.VerifyTokenResponse, error) {
	res, err := client.service.VerifyRefreshToken(ctx, &dto.VerifyTokenRequest{Token: in.Token})
	if err != nil {
		return nil, err
	}
	return &dto.VerifyTokenResponse{UserID: res.UserID, Status: types.HttpStatus_Ok}, nil
}

func (client *authClient) CreateTokens(ctx context.Context, in *dto.CreateTokensRequest) (*models.Tokens, error) {
	res, err := client.service.CreateTokens(ctx, &dto.CreateTokensRequest{UserID: in.UserID})
	if err != nil {
		return nil, err
	}
	return &models.Tokens{AccessToken: res.AccessToken, RefreshToken: res.RefreshToken}, nil
}

func (client *authClient) DeleteRefreshToken(ctx context.Context, in *dto.DeleteRefreshTokenRequest, opts ...grpc.CallOption) (*dto.DeleteRefreshTokenResponse, error) {
	_, err := client.service.DeleteRefreshToken(ctx, &dto.DeleteRefreshTokenRequest{Token: in.Token})
	if err != nil {
		return &dto.DeleteRefreshTokenResponse{Status: types.HttpStatus_Error}, err
	}
	return &dto.DeleteRefreshTokenResponse{Status: types.HttpStatus_Ok}, nil
}
