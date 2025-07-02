package mysql

import (
	"context"
	"database/sql"
	"ghorkov32/proletariat-budget-be/openapi"
)

type SavingGoalRepoImpl struct {
	db *sql.DB
}

func (s SavingGoalRepoImpl) Create(ctx context.Context, savingsGoal openapi.SavingsGoalRequest) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepoImpl) Update(ctx context.Context, id string, savingsGoal openapi.SavingsGoalRequest) error {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepoImpl) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepoImpl) GetByID(ctx context.Context, id string) (*openapi.SavingsGoal, error) {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepoImpl) List(ctx context.Context) ([]openapi.SavingsGoal, error) {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepoImpl) CreateWithdrawal(ctx context.Context, withdrawal openapi.SavingsWithdrawalRequest) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepoImpl) DeleteWithdrawal(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepoImpl) GetWithdrawalByID(ctx context.Context, id string) (*openapi.SavingsWithdrawal, error) {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepoImpl) CreateContribution(ctx context.Context, contribution openapi.SavingsContributionRequest) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepoImpl) DeleteContribution(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepoImpl) GetContributionByID(ctx context.Context, id string) (*openapi.SavingsContribution, error) {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepoImpl) ListSavingsTransactions(ctx context.Context, params openapi.ListSavingsTransactionsParams) (*openapi.SavingsTransactionList, error) {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepoImpl) FindOrCreateTags(ctx context.Context, tags []string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepoImpl) LinkTagsToSavingsGoal(ctx context.Context, tags []string, savingsGoalId string) error {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepoImpl) LinkTagsToWithdrawal(ctx context.Context, tags []string, withdrawalId string) error {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepoImpl) LinkTagsToContribution(ctx context.Context, tags []string, contributionId string) error {
	//TODO implement me
	panic("implement me")
}

func (s SavingGoalRepoImpl) CreateTag(ctx context.Context, tag openapi.TagRequest) (string, error) {
	//TODO implement me
	panic("implement me")
}
