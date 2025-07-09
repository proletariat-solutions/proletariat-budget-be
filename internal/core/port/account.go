package port

import (
	"context"
	"ghorkov32/proletariat-budget-be/openapi"
)

type AccountRepo interface {
	Create(ctx context.Context, account openapi.AccountRequest) (string, error)
	GetByID(ctx context.Context, id string) (*openapi.Account, error)
	Update(ctx context.Context, account openapi.Account) error
	Deactivate(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, params openapi.ListAccountsParams) (*openapi.AccountList, error)
	HasTransactions(ctx context.Context, id string) (bool, error)
}
