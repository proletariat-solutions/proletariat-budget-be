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
	// Validate account
	account, err := u.validateAccount(
		ctx,
		expenditure.Transaction.AccountID,
	)
	if err != nil {
		return nil, err
	}

	// Validate category
	err = u.validateCategory(
		ctx,
		expenditure.Category.ID,
	)
	if err != nil {
		return nil, err
	}

	// Process transaction
	err = u.processTransaction(
		ctx,
		account,
		&expenditure,
	)
	if err != nil {
		return nil, err
	}

	// Create expenditure record
	expID, err := u.createExpenditureRecord(
		ctx,
		expenditure,
	)
	if err != nil {
		return nil, err
	}

	// Link tags if present
	err = u.linkTags(
		ctx,
		expID,
		expenditure.Tags,
	)
	if err != nil {
		return nil, err
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

func (u *ExpenditureUseCase) validateAccount(
	ctx context.Context,
	accountID string,
) (
	*domain.Account,
	error,
) {
	account, err := (*u.accountRepo).GetByID(
		ctx,
		accountID,
	)
	if err != nil {
		if errors.Is(
			err,
			port.ErrRecordNotFound,
		) {
			return nil, domain.ErrAccountNotFound
		}
		return nil, err
	}

	if !account.Active {
		return nil, domain.ErrAccountInactive
	}

	return account, nil
}

func (u *ExpenditureUseCase) validateCategory(
	ctx context.Context,
	categoryID string,
) error {
	category, err := (*u.categoryRepo).GetByID(
		ctx,
		categoryID,
	)
	if err != nil {
		if errors.Is(
			err,
			port.ErrRecordNotFound,
		) {
			return domain.ErrCategoryNotFound
		}
		return err
	}

	if !category.Active {
		return domain.ErrCategoryInactive
	}

	return nil
}

func (u *ExpenditureUseCase) processTransaction(
	ctx context.Context,
	account *domain.Account,
	expenditure *domain.Expenditure,
) error {
	if !account.HasSufficientBalance(expenditure.Transaction.Amount) {
		return domain.ErrInsufficientBalance
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
		return err
	}

	expenditure.Transaction.ID = &txID

	return (*u.accountRepo).Update(
		ctx,
		*account,
	)
}

func (u *ExpenditureUseCase) createExpenditureRecord(
	ctx context.Context,
	expenditure domain.Expenditure,
) (
	string,
	error,
) {
	return (*u.expenditureRepo).Create(
		ctx,
		expenditure,
	)
}

func (u *ExpenditureUseCase) linkTags(
	ctx context.Context,
	expID string,
	tags *[]*domain.Tag,
) error {
	if tags == nil || len(*tags) == 0 {
		return nil
	}

	err := (*u.tagsRepo).LinkTagsToType(
		ctx,
		expID,
		tags,
	)
	if err != nil {
		if errors.Is(
			err,
			port.ErrForeignKeyViolation,
		) {
			return domain.ErrTagNotFound
		}
		return err
	}

	return nil
}
