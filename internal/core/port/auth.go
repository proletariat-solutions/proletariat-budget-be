package port

import (
	"context"

	"ghorkov32/proletariat-budget-be/internal/core/domain"
)

type AuthRepo interface {
	CreateUser(ctx context.Context, user domain.User) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, id string, user domain.User) error

	// Token management
	CreateToken(ctx context.Context, userID string) (*domain.AuthToken, error)
	ValidateToken(ctx context.Context, token string) (*domain.User, error)
	RevokeToken(ctx context.Context, token string) error
}
