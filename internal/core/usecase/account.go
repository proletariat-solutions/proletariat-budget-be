package usecase

import (
	"context"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"ghorkov32/proletariat-budget-be/openapi"
)

type AccountUseCase struct {
	accountRepo *port.AccountRepo
}

func NewAccountUseCase(accountRepo *port.AccountRepo) *AccountUseCase {
	return &AccountUseCase{
		accountRepo: accountRepo,
	}
}

func (a *AccountUseCase) Create(ctx context.Context, account openapi.AccountRequest) (string, error) {
	ID, err := (*a.accountRepo).Create(ctx, account)
	if err != nil {
		return "", err
	}
	return ID, nil
}

func (a *AccountUseCase) GetByID(ctx context.Context, id string) (*openapi.Account, error) {
	account, err := (*a.accountRepo).GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (a *AccountUseCase) Update(ctx context.Context, account openapi.Account) error {
	err := (*a.accountRepo).Update(ctx, account)
	if err != nil {
		return err
	}
	return nil
}

func (a *AccountUseCase) Delete(ctx context.Context, id string) error {
	hasTransactions, errHasTransactions := (*a.accountRepo).HasTransactions(ctx, id)
	if errHasTransactions != nil {
		return errHasTransactions
	}
	if hasTransactions {
		return domain.ErrAccountHasTransactions
	}
	err := (*a.accountRepo).Deactivate(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (a *AccountUseCase) List(ctx context.Context, params openapi.ListAccountsParams) (*openapi.AccountList, error) {
	accounts, err := (*a.accountRepo).List(ctx, params)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (a *AccountUseCase) HasTransactions(ctx context.Context, id string) (bool, error) {
	hasTransactions, err := (*a.accountRepo).HasTransactions(ctx, id)
	if err != nil {
		return false, err
	}
	return hasTransactions, nil
}
