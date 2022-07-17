package grpc

import (
	"context"
	"fmt"

	pb "gitlab.com/g6834/team41/analytics/api/auth"

	"gitlab.com/g6834/team41/analytics/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

func (c *Client) Validate(ctx context.Context, login string, tokens models.TokenPair) (models.TokenPair, error) {
	conn, err := grpc.Dial(c.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("can't connect to auth: %w", err)
	}
	defer conn.Close()

	auth := pb.NewAuthServiceClient(conn)

	in := pb.ValidateRequest{
		Login:        login,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	r, err := auth.Validate(ctx, &in)
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("auth validate error: %w", err)
	}

	newTokens := models.TokenPair{
		AccessToken:  r.AccessToken,
		RefreshToken: r.RefreshToken,
	}
	return newTokens, nil
}
