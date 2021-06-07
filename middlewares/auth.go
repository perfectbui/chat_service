package middlewares

import (
	"context"

	"github.com/perfectbui/chat/clients"
	"github.com/perfectbui/chat/errors"
	"github.com/perfectbui/chat/pb/dto"
)

type Token struct {
	RefreshToken string
	AccessToken  string
}

type RequestInfo struct {
	UserID int64
}

func CheckAuth(accessToken, refreshToken string) (*RequestInfo, *Token, error) {
	ctx := context.Background()
	if len(accessToken) > 0 {
		if verifyResp, err := clients.AuthClient.VerifyAccessToken(ctx, &dto.VerifyTokenRequest{Token: accessToken}); err == nil {
			return &RequestInfo{UserID: verifyResp.UserID}, &Token{AccessToken: accessToken, RefreshToken: refreshToken}, nil
		}
	}
	return doRefreshToken(ctx, refreshToken)
}

func doRefreshToken(ctx context.Context, refreshToken string) (*RequestInfo, *Token, error) {
	if len(refreshToken) > 0 {
		verifyResp, err := clients.AuthClient.VerifyRefreshToken(ctx, &dto.VerifyTokenRequest{Token: refreshToken})
		if err != nil {
			return nil, nil, err
		}
		_, err = clients.AuthClient.DeleteRefreshToken(ctx, &dto.DeleteRefreshTokenRequest{Token: refreshToken})
		if err != nil {
			return nil, nil, err
		}
		createResp, err := clients.AuthClient.CreateTokens(ctx, &dto.CreateTokensRequest{UserID: verifyResp.UserID})
		if err != nil {
			return nil, nil, err
		}
		return &RequestInfo{UserID: verifyResp.UserID}, &Token{AccessToken: createResp.AccessToken, RefreshToken: createResp.RefreshToken}, nil
	}
	return nil, nil, errors.ErrorTokenIsNotValid
}
