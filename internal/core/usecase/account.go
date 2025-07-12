package usecase

import (
	"context"
	"errors"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"strings"
)

type AccountUseCase struct {
	accountRepo *port.AccountRepo
}

func NewAccountUseCase(accountRepo *port.AccountRepo) *AccountUseCase {
	return &AccountUseCase{
		accountRepo: accountRepo,
	}
}

func (a *AccountUseCase) Create(ctx context.Context, account domain.Account) (*string, error) {
	ID, err := (*a.accountRepo).Create(ctx, account)
	if err != nil {
		if errors.Is(err, port.ErrForeignKeyViolation) {
			if strings.Contains(err.Error(), "currency") {
				return nil, domain.ErrInvalidCurrency
			}
			if strings.Contains(err.Error(), "fk_accounts_owner") {
				return nil, domain.ErrAccountOwnerNotFound
			}
		}
		return nil, err
	}
	return ID, nil
}

func (a *AccountUseCase) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	account, err := (*a.accountRepo).GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, port.ErrRecordNotFound) {
			return nil, domain.ErrAccountNotFound
		}
		return nil, err
	}
	return account, nil
}

func (a *AccountUseCase) Update(ctx context.Context, account domain.Account) (*domain.Account, error) {
	err := (*a.accountRepo).Update(ctx, account)
	if err != nil {
		if errors.Is(err, port.ErrForeignKeyViolation) {
			if strings.Contains(err.Error(), "currency") {
				return nil, domain.ErrInvalidCurrency
			}
			if strings.Contains(err.Error(), "fk_accounts_owner") {
				return nil, domain.ErrAccountOwnerNotFound
			}
		}
		return nil, err
	}
	updatedAccount, errGet := (*a.accountRepo).GetByID(ctx, *account.ID)
	if errGet != nil {
		if errors.Is(errGet, port.ErrRecordNotFound) {
			return nil, domain.ErrAccountNotFound
		}
		return nil, errGet
	}
	return updatedAccount, nil
}

func (a *AccountUseCase) Deactivate(ctx context.Context, id string) error {
	account, err := (*a.accountRepo).GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, port.ErrRecordNotFound) {
			return domain.ErrAccountNotFound
		}
	}
	if !*account.Active {
		return domain.ErrAccountAlreadyInactive
	}
	err = (*a.accountRepo).SetActive(ctx, id, false)
	if err != nil {
		return err
	}
	return nil
}

func (a *AccountUseCase) Activate(ctx context.Context, id string) error {
	account, err := (*a.accountRepo).GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, port.ErrRecordNotFound) {
			return domain.ErrAccountNotFound
		}
	}
	if *account.Active {
		return domain.ErrAccountAlreadyActive
	}
	err = (*a.accountRepo).SetActive(ctx, id, true)
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
	err := (*a.accountRepo).Delete(ctx, id)
	if err != nil {
		if errors.Is(err, port.ErrRecordNotFound) {
			return domain.ErrAccountNotFound
		}
		return err
	}
	return nil
}

func (a *AccountUseCase) List(ctx context.Context, params domain.AccountListParams) (*domain.AccountList, error) {
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
