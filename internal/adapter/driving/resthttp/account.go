package resthttp

import (
	"context"
	"errors"

	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"ghorkov32/proletariat-budget-be/openapi"
	"github.com/rs/zerolog/log"
)

func (c *Controller) ListAccounts(ctx context.Context, request openapi.ListAccountsRequestObject) (openapi.ListAccountsResponseObject, error) {
	accounts, err := c.useCases.Account.List(ctx, *FromOAPIAccountListParams(&request.Params))
	if err != nil {
		log.Err(err).Msg("Failed to list accounts")

		return openapi.ListAccounts500JSONResponse{
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Failed to list accounts",
			},
		}, nil
	}

	return openapi.ListAccounts200JSONResponse(
		*ToOAPIAccountList(accounts),
	), nil
}

func (c *Controller) CreateAccount(ctx context.Context, request openapi.CreateAccountRequestObject) (openapi.CreateAccountResponseObject, error) {
	id, err := c.useCases.Account.Create(ctx, *FromOAPIAccountRequest(request.Body))
	if err != nil {
		log.Err(err).Msg("Failed to create account")
		if errors.Is(err, domain.ErrMemberNotFound) || errors.Is(err, domain.ErrInvalidCurrency) || errors.Is(err, domain.ErrMemberInactive) {
			return openapi.CreateAccount400JSONResponse{
				N400JSONResponse: openapi.N400JSONResponse{
					Message: err.Error(),
				},
			}, nil
		}

		return openapi.CreateAccount500JSONResponse{ // coverage-ignore
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Failed to create account",
			},
		}, nil
	}

	account, err := c.useCases.Account.GetByID(ctx, *id)
	if err != nil {
		log.Err(err).Msg("Failed to get created account")

		return openapi.CreateAccount500JSONResponse{ // coverage-ignore
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Failed to get created account",
			},
		}, nil
	}

	return openapi.CreateAccount201JSONResponse(*ToOAPIAccount(*account)), nil
}

func (c *Controller) DeleteAccount(ctx context.Context, request openapi.DeleteAccountRequestObject) (openapi.DeleteAccountResponseObject, error) {
	err := c.useCases.Account.Delete(ctx, request.Id)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			return openapi.DeleteAccount404JSONResponse{
				N404JSONResponse: openapi.N404JSONResponse{
					Message: err.Error(),
				},
			}, nil
		} else if errors.Is(err, domain.ErrAccountHasTransactions) {
			return openapi.DeleteAccount409JSONResponse{
				N409JSONResponse: openapi.N409JSONResponse{
					Message: err.Error(),
				},
			}, nil
		} else {
			log.Err(err).Msg("Failed to delete account")

			return openapi.DeleteAccount500JSONResponse{ // coverage-ignore
				N500JSONResponse: openapi.N500JSONResponse{
					Message: "Failed to delete account",
				},
			}, nil
		}
	}

	return openapi.DeleteAccount204Response{}, nil
}

func (c *Controller) GetAccount(ctx context.Context, request openapi.GetAccountRequestObject) (openapi.GetAccountResponseObject, error) {
	account, err := c.useCases.Account.GetByID(ctx, request.Id)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			return openapi.GetAccount404JSONResponse{
				N404JSONResponse: openapi.N404JSONResponse{
					Message: err.Error(),
				},
			}, nil
		} else {
			log.Err(err).Msg("Failed to get account")

			return openapi.GetAccount500JSONResponse{ // coverage-ignore
				N500JSONResponse: openapi.N500JSONResponse{
					Message: "Failed to get account",
				},
			}, nil
		}
	}

	return openapi.GetAccount200JSONResponse(*ToOAPIAccount(*account)), nil
}

func (c *Controller) UpdateAccount(ctx context.Context, request openapi.UpdateAccountRequestObject) (openapi.UpdateAccountResponseObject, error) {
	account, err := c.useCases.Account.Update(ctx, *FromOAPIAccount(request.Body))
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			return openapi.UpdateAccount404JSONResponse{
				N404JSONResponse: openapi.N404JSONResponse{
					Message: err.Error(),
				},
			}, nil
		} else if errors.Is(err, domain.ErrMemberNotFound) ||
			errors.Is(err, domain.ErrInvalidCurrency) ||
			errors.Is(err, port.ErrInvalidDataFormat) {
			return openapi.UpdateAccount400JSONResponse{
				N400JSONResponse: openapi.N400JSONResponse{
					Message: err.Error(),
				},
			}, nil
		} else {
			log.Err(err).Msg("Failed to update account")

			return openapi.UpdateAccount500JSONResponse{ // coverage-ignore
				N500JSONResponse: openapi.N500JSONResponse{
					Message: "Failed to update account",
				},
			}, nil
		}
	}

	return openapi.UpdateAccount200JSONResponse(*ToOAPIAccount(*account)), nil
}

func (c *Controller) GetBalances(ctx context.Context, request openapi.GetBalancesRequestObject) (openapi.GetBalancesResponseObject, error) {
	return openapi.GetBalances501Response{}, nil
}

func (c *Controller) CanDeleteAccount(ctx context.Context, request openapi.CanDeleteAccountRequestObject) (openapi.CanDeleteAccountResponseObject, error) {
	hasTransactions, err := c.useCases.Account.HasTransactions(ctx, request.Id)
	if err != nil {
		log.Err(err).Msg("Failed to check if account has transactions")

		return openapi.CanDeleteAccount500JSONResponse{ // coverage-ignore
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Failed to check if account has transactions",
			},
		}, nil
	}
	var reason string
	if hasTransactions {
		reason = "Account has transactions"
	}

	return openapi.CanDeleteAccount200JSONResponse{
		CanDelete: !hasTransactions,
		Reason:    &reason,
	}, nil
}

func (c *Controller) ActivateAccount(ctx context.Context, request openapi.ActivateAccountRequestObject) (openapi.ActivateAccountResponseObject, error) {
	err := c.useCases.Account.Activate(ctx, request.Id)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			return openapi.ActivateAccount404JSONResponse{
				N404JSONResponse: openapi.N404JSONResponse{
					Message: err.Error(),
				},
			}, nil
		} else if errors.Is(err, domain.ErrAccountAlreadyActive) {
			return openapi.ActivateAccount400JSONResponse{
				N400JSONResponse: openapi.N400JSONResponse{
					Message: err.Error(),
				},
			}, nil
		}
		log.Err(err).Msg("Failed to activate account")

		return openapi.ActivateAccount500JSONResponse{ // coverage-ignore
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Failed to activate account",
			},
		}, nil
	}

	return openapi.ActivateAccount200Response{}, nil
}

func (c *Controller) DeactivateAccount(ctx context.Context, request openapi.DeactivateAccountRequestObject) (openapi.DeactivateAccountResponseObject, error) {
	err := c.useCases.Account.Deactivate(ctx, request.Id)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			return openapi.DeactivateAccount404JSONResponse{
				N404JSONResponse: openapi.N404JSONResponse{
					Message: err.Error(),
				},
			}, nil
		} else if errors.Is(err, domain.ErrAccountAlreadyInactive) {
			return openapi.DeactivateAccount400JSONResponse{
				N400JSONResponse: openapi.N400JSONResponse{
					Message: err.Error(),
				},
			}, nil
		}
		log.Err(err).Msg("Failed to deactivate account")

		return openapi.DeactivateAccount500JSONResponse{ // coverage-ignore
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Failed to deactivate account",
			},
		}, nil
	}

	return openapi.DeactivateAccount204Response{}, nil
}
