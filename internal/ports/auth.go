package ports

import (
	"context"

	"gitlab.com/g6834/team41/analytics/internal/models"
)

type AuthService interface {
	Validate(ctx context.Context, login string, tokens models.TokenPair) (models.TokenPair, error)
}
