package mysql

import (
	"context"
	"database/sql"
	"ghorkov32/proletariat-budget-be/openapi"
)

type SavingGoalRepo struct {
	db *sql.DB
}

func (s SavingGoalRepo) ListCategories(ctx context.Context) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepo) GetCategory(ctx context.Context, id string) (*openapi.SavingsCategory, error) {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepo) CreateCategory(ctx context.Context, category openapi.SavingsCategoryRequest) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepo) UpdateCategory(ctx context.Context, id string, category openapi.SavingsCategoryRequest) error {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepo) DeleteCategory(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepo) FindOrCreateTags(ctx context.Context, tags []string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepo) LinkTagsToSavingsGoal(ctx context.Context, tags []string, savingsGoalId string) error {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepo) LinkTagsToWithdrawal(ctx context.Context, tags []string, withdrawalId string) error {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepo) LinkTagsToContribution(ctx context.Context, tags []string, contributionId string) error {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepo) CreateTag(ctx context.Context, tag openapi.TagRequest) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepo) Create(ctx context.Context, savingsGoal openapi.SavingsGoalRequest) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepo) Update(ctx context.Context, id string, savingsGoal openapi.SavingsGoalRequest) error {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepo) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepo) GetByID(ctx context.Context, id string) (*openapi.SavingsGoal, error) {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepo) List(ctx context.Context) ([]openapi.SavingsGoal, error) {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepo) CreateWithdrawal(ctx context.Context, withdrawal openapi.SavingsWithdrawalRequest) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepo) DeleteWithdrawal(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepo) GetWithdrawalByID(ctx context.Context, id string) (*openapi.SavingsWithdrawal, error) {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepo) CreateContribution(ctx context.Context, contribution openapi.SavingsContributionRequest) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepo) DeleteContribution(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepo) GetContributionByID(ctx context.Context, id string) (*openapi.SavingsContribution, error) {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepo) ListSavingsTransactions(ctx context.Context, params openapi.ListSavingsTransactionsParams) (*openapi.SavingsTransactionList, error) {
	//TODO implement me
	panic("implement me")
}
