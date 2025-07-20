package port

import (
	"context"

	"ghorkov32/proletariat-budget-be/internal/core/domain"
)

type AccountRepo interface {
	Create(ctx context.Context, account domain.Account) (*string, error)
	GetByID(ctx context.Context, id string) (*domain.Account, error)
	Update(ctx context.Context, account domain.Account) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, params domain.AccountListParams) (*domain.AccountList, error)
	HasTransactions(ctx context.Context, id string) (bool, error)
}
