package usecase

import (
	"context"
	"errors"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
)

type ExpenditureUseCase struct {
	expenditureRepo *port.ExpenditureRepo
	accountRepo     *port.AccountRepo
	tagsRepo        *port.TagsRepo
	categoryRepo    *port.CategoryRepo
	transactionRepo *port.TransactionRepo
}

func NewExpenditureUseCase(
	expenditureRepo *port.ExpenditureRepo,
	accountRepo *port.AccountRepo,
	tagsRepo *port.TagsRepo,
	categoryRepo *port.CategoryRepo,
	transactionRepo *port.TransactionRepo,
) *ExpenditureUseCase {
	return &ExpenditureUseCase{
		expenditureRepo: expenditureRepo,
		accountRepo:     accountRepo,
		tagsRepo:        tagsRepo,
		categoryRepo:    categoryRepo,
		transactionRepo: transactionRepo,
	}
}

func (u *ExpenditureUseCase) Create(
	ctx context.Context,
	expenditure domain.Expenditure,
) (
	*domain.Expenditure,
	error,
) {
	// Checking if account exists and is active
	account, err := (*u.accountRepo).GetByID(
		ctx,
		expenditure.Transaction.AccountID,
	)
	if err != nil {
		if errors.Is(
			err,
			port.ErrRecordNotFound,
		) {
			return nil, domain.ErrAccountNotFound
		}
		return nil, err
	} else if !account.Active {
		return nil, domain.ErrAccountInactive
	}

	// Checking if category exists
	category, err := (*u.categoryRepo).GetByID(
		ctx,
		expenditure.Category.ID,
	)
	if err != nil {
		if errors.Is(
			err,
			port.ErrRecordNotFound,
		) {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, err
	} else if !category.Active {
		return nil, domain.ErrCategoryInactive
	}

	// Generate transaction
	if !account.HasSufficientBalance(expenditure.Transaction.Amount) {
		return nil, domain.ErrInsufficientBalance
	}

	account.DebitBalance(expenditure.Transaction.Amount)
	expenditure.Transaction.BalanceAfter = &account.CurrentBalance
	statusCompleted := domain.TransactionStatusCompleted
	expenditure.Transaction.Status = &statusCompleted

	txID, err := (*u.transactionRepo).Create(
		ctx,
		*expenditure.Transaction,
	)
	if err != nil {
		return nil, err
	} else {
		expenditure.Transaction.ID = &txID
	}

	// Update account balance
	err = (*u.accountRepo).Update(
		ctx,
		*account,
	)
	if err != nil {
		return nil, err
	}

	// Create expenditure
	expID, err := (*u.expenditureRepo).Create(
		ctx,
		expenditure,
	)
	if err != nil {
		return nil, err
	}

	if expenditure.Tags != nil && len(*expenditure.Tags) > 0 {
		// Update tags
		err = (*u.tagsRepo).LinkTagsToType(
			ctx,
			expID,
			expenditure.Tags,
		)

		if err != nil {
			if errors.Is(
				err,
				port.ErrForeignKeyViolation,
			) {
				return nil, domain.ErrTagNotFound
			}
			return nil, err
		}
	}

	return (*u.expenditureRepo).GetByID(
		ctx,
		expID,
	)
}

func (u *ExpenditureUseCase) Get(
	ctx context.Context,
	id string,
) (
	*domain.Expenditure,
	error,
) {
	expenditure, err := (*u.expenditureRepo).GetByID(
		ctx,
		id,
	)
	if err != nil {
		if errors.Is(
			err,
			port.ErrRecordNotFound,
		) {
			return nil, domain.ErrExpenditureNotFound
		}
		return nil, err
	}
	return expenditure, nil
}

func (u *ExpenditureUseCase) List(
	ctx context.Context,
	params domain.ExpenditureListParams,
) (
	*domain.ExpenditureList,
	error,
) {
	return (*u.expenditureRepo).FindExpenditures(
		ctx,
		params,
	)
}
