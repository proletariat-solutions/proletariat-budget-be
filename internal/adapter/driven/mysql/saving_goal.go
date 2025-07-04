package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"ghorkov32/proletariat-budget-be/openapi"
	"sort"
	"strings"
)

type SavingGoalRepoImpl struct {
	db           *sql.DB
	categoryRepo *port.CategoryRepo
	tagsRepo     *port.TagsRepo
}

func (s SavingGoalRepoImpl) Create(ctx context.Context, savingsGoal openapi.SavingsGoalRequest) (string, error) {
	queryInsert := `INSERT INTO savings_goals
						(name,
						 category_id,
						 description,
						 target_amount, 
						 currency, 
						 target_date, 
						 initial_amount, 
						 current_amount, 
						 percent_complete, 
						 account_id, 
						 priority,
						 auto_contribute, 
						 auto_contribute_amount, 
						 auto_contribute_frequency, 
						 status, 
						 created_at,
						 updated_at)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`
	percentComplete := float32(0)
	if savingsGoal.InitialAmount != nil {
		percentComplete = (*savingsGoal.InitialAmount / savingsGoal.TargetAmount) * 100
	}

	result, err := s.db.ExecContext(ctx, queryInsert,
		savingsGoal.Name,
		savingsGoal.Category.Id,
		savingsGoal.Description,
		savingsGoal.TargetAmount,
		savingsGoal.Currency,
		savingsGoal.TargetDate,
		savingsGoal.InitialAmount,
		savingsGoal.InitialAmount,
		percentComplete,
		savingsGoal.AccountId,
		savingsGoal.Priority,
		savingsGoal.AutoContribute,
		savingsGoal.AutoContributeAmount,
		savingsGoal.AutoContributeFrequency,
		"ACTIVE")

	if err != nil {
		return "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", id), nil
}

func (s SavingGoalRepoImpl) Update(ctx context.Context, id string, savingsGoal openapi.SavingsGoalRequest) error {
	queryUpdate := `UPDATE savings_goals
					SET name=?,
						category_id=?,
						description=?,
						target_amount=?,
						currency=?,
						target_date=?,
						initial_amount=?,
						account_id=?,
						priority=?,
						auto_contribute=?,
						auto_contribute_amount=?,
						auto_contribute_frequency=?,
						status=?,
						updated_at=NOW()
					WHERE id = ?`
	_, err := s.db.ExecContext(ctx, queryUpdate,
		savingsGoal.Name,
		savingsGoal.Category.Id,
		savingsGoal.Description,
		savingsGoal.TargetAmount,
		savingsGoal.Currency,
		savingsGoal.TargetDate,
		savingsGoal.InitialAmount,
		savingsGoal.AccountId,
		savingsGoal.Priority,
		savingsGoal.AutoContribute,
		savingsGoal.AutoContributeAmount,
		savingsGoal.AutoContributeFrequency,
		"active",
		id)
	if err != nil {
		return err
	}
	return nil
}

func (s SavingGoalRepoImpl) Delete(ctx context.Context, id string) error {
	// Won't delete the record, just updating the status to "inactive"
	query := `UPDATE savings_goals SET status =?, updated_at = NOW() WHERE id =?`
	_, err := s.db.ExecContext(
		ctx,
		query,
		"inactive",
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete savings goal: %w", err)
	}
	return nil
}

func (s SavingGoalRepoImpl) GetByID(ctx context.Context, id string) (*openapi.SavingsGoal, error) {
	query := `SELECT id, name, category_id, description, target_amount, currency, target_date, initial_amount, current_amount, percent_complete, account_id, priority, auto_contribute, auto_contribute_amount, auto_contribute_frequency, status, created_at, updated_at FROM savings_goals WHERE id=? AND status='active'`
	var savingsGoal openapi.SavingsGoal
	var categoryID string
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&savingsGoal.Id,
		&savingsGoal.Name,
		&categoryID,
		&savingsGoal.Description,
		&savingsGoal.TargetAmount,
		&savingsGoal.Currency,
		&savingsGoal.TargetDate,
		&savingsGoal.InitialAmount,
		&savingsGoal.CurrentAmount,
		&savingsGoal.PercentComplete,
		&savingsGoal.AccountId,
		&savingsGoal.Priority,
		&savingsGoal.AutoContribute,
		&savingsGoal.AutoContributeAmount,
		&savingsGoal.AutoContributeFrequency,
		&savingsGoal.Status,
		&savingsGoal.CreatedAt,
		&savingsGoal.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("savings goal not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to select savings goal: %w", err)
	}

	category, err := (*s.categoryRepo).GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	savingsGoal.Category = *category
	return &savingsGoal, nil
}

func (s SavingGoalRepoImpl) List(ctx context.Context) ([]openapi.SavingsGoal, error) {
	query := `SELECT 
				id, 
				name, 
				category_id, 
				description, 
				target_amount, 
				currency, 
				target_date, 
				initial_amount, 
				current_amount, 
				percent_complete, 
				account_id, 
				priority, 
				auto_contribute, 
				auto_contribute_amount, 
				auto_contribute_frequency, 
				status, 
				created_at, 
				updated_at,
			    GROUP_CONCAT(sgt.tag_id ORDER BY sgt.tag_id SEPARATOR ',') as tags
			FROM savings_goals
			left join savings_goal_tags sgt ON savings_goals.id = sgt.savings_goal_id
			WHERE status != 'inactive'`
	var savingsGoals []openapi.SavingsGoal
	var ids []string
	tagsByID := make(map[string][]string)
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to select savings goals: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var savingsGoal openapi.SavingsGoal
		var categoryID string
		var tags string
		err := rows.Scan(
			&savingsGoal.Id,
			&savingsGoal.Name,
			&categoryID,
			&savingsGoal.Description,
			&savingsGoal.TargetAmount,
			&savingsGoal.Currency,
			&savingsGoal.TargetDate,
			&savingsGoal.InitialAmount,
			&savingsGoal.CurrentAmount,
			&savingsGoal.PercentComplete,
			&savingsGoal.AccountId,
			&savingsGoal.Priority,
			&savingsGoal.AutoContribute,
			&savingsGoal.AutoContributeAmount,
			&savingsGoal.AutoContributeFrequency,
			&savingsGoal.Status,
			&savingsGoal.CreatedAt,
			&savingsGoal.UpdatedAt,
			&tags,
		)

		tagsByID[*savingsGoal.Id] = strings.Split(tags, ",")

		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		ids = append(ids, *savingsGoal.Id)
		category, err := (*s.categoryRepo).GetByID(ctx, categoryID)
		if err != nil {
			return nil, fmt.Errorf("failed to get category: %w", err)
		}
		savingsGoal.Category = *category
	}

	tags, err := (*s.tagsRepo).ListByType(ctx, "savings_goal", &ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	for _, goal := range savingsGoals {
		goalTags := make([]openapi.Tag, 0, len(tagsByID[*goal.Id]))
		for _, tagID := range tagsByID[*goal.Id] {
			idx := sort.Search(len(tags), func(i int) bool { return tags[i].Id == tagID })
			if idx >= 0 {
				goalTags = append(goalTags, tags[idx])
			}
		}
		goal.Tags = &goalTags
	}

	return savingsGoals, nil
}

func (s SavingGoalRepoImpl) MarkAsCompleted(ctx context.Context, id string) error {
	query := `
	UPDATE
	savings_goals
	SET
	status =?, updated_at = NOW()
	WHERE
	id =?`
	_, err := s.db.ExecContext(
		ctx,
		query,
		"completed",
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to mark savings goal as completed: %w", err)
	}
	return nil
}

func (s SavingGoalRepoImpl) MarkAsAbandoned(ctx context.Context, id string) error {
	query := `
	UPDATE
	savings_goals
	SET
	status =?, updated_at = NOW()
	WHERE
	id =?`
	_, err := s.db.ExecContext(
		ctx,
		query,
		"abandoned",
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to mark savings goal as abandoned: %w", err)
	}
	return nil
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
