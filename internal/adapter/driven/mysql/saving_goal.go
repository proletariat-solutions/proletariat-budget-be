package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"ghorkov32/proletariat-budget-be/openapi"
	"strings"
)

type SavingGoalRepoImpl struct {
	db       *sql.DB
	tagsRepo *port.TagsRepo
}

func NewSavingGoalRepo(
	db *sql.DB,
	tagsRepo *port.TagsRepo,
) port.SavingsGoalRepo {
	return &SavingGoalRepoImpl{db: db, tagsRepo: tagsRepo}
}

func (s SavingGoalRepoImpl) Create(
	ctx context.Context,
	savingsGoal openapi.SavingsGoalRequest,
) (
	string,
	error,
) {
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

	result, err := s.db.ExecContext(
		ctx,
		queryInsert,
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
		"ACTIVE",
	)

	if err != nil {
		return "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"%d",
		id,
	), nil
}

func (s SavingGoalRepoImpl) Update(
	ctx context.Context,
	id string,
	savingsGoal openapi.SavingsGoalRequest,
) error {
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
	_, err := s.db.ExecContext(
		ctx,
		queryUpdate,
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
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s SavingGoalRepoImpl) Delete(
	ctx context.Context,
	id string,
) error {
	// Won't delete the record, just updating the status to "inactive"
	query := `UPDATE savings_goals SET status =?, updated_at = NOW() WHERE id =?`
	_, err := s.db.ExecContext(
		ctx,
		query,
		"inactive",
		id,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to delete savings goal: %w",
			err,
		)
	}
	return nil
}

func (s SavingGoalRepoImpl) GetByID(
	ctx context.Context,
	id string,
) (
	*openapi.SavingsGoal,
	error,
) {
	query := `select sg.id,
					   sg.name,
					   sg.description,
					   sg.target_amount,
					   sg.currency,
					   sg.target_date,
					   sg.initial_amount,
					   sg.current_amount,
					   sg.percent_complete,
					   sg.account_id,
					   sg.priority,
					   sg.auto_contribute,
					   sg.auto_contribute_amount,
					   sg.auto_contribute_frequency,
					   sg.status,
					   sg.projected_completion_date,
					   sg.created_at,
					   sg.updated_at,
					   GROUP_CONCAT(sgt.tag_id ORDER BY sgt.tag_id SEPARATOR ',') as tags,
					   c.id,
					   c.name,
					   c.description,
					   c.color,
					   c.background_color,
					   c.active,
					   c.category_type
				from savings_goals sg
						 inner join categories c ON sg.category_id = c.id
						 left join savings_goal_tags sgt ON sg.id = sgt.savings_goal_id
				where sg.id =?
				group by sg.id, sg.name, sg.description, sg.target_amount, sg.currency, sg.target_date, sg.initial_amount,
						 sg.current_amount, sg.percent_complete, sg.account_id, sg.priority, sg.auto_contribute,
						 sg.auto_contribute_amount, sg.auto_contribute_frequency, sg.status, sg.projected_completion_date,
						 sg.created_at, sg.updated_at
`
	var savingsGoal openapi.SavingsGoal
	var tags string
	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&savingsGoal.Id,
		&savingsGoal.Name,
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
		&savingsGoal.ProjectedCompletionDate,
		&savingsGoal.CreatedAt,
		&savingsGoal.UpdatedAt,
		&tags,
		&savingsGoal.Category.Id,
		&savingsGoal.Category.Name,
		&savingsGoal.Category.Description,
		&savingsGoal.Category.Color,
		&savingsGoal.Category.BackgroundColor,
		&savingsGoal.Category.Active,
	)
	if errors.Is(
		err,
		sql.ErrNoRows,
	) {
		return nil, fmt.Errorf("savings goal not found")
	} else if err != nil {
		return nil, fmt.Errorf(
			"failed to select savings goal: %w",
			err,
		)
	}

	/*	savingsGoal.Tags, err = (*s.tagsRepo).GetByIDs(ctx, strings.Split(tags, ","))
		if err != nil {
			return nil, fmt.Errorf("failed to get tags: %w", err)
		}*/

	return &savingsGoal, nil
}

func (s SavingGoalRepoImpl) List(ctx context.Context) (
	[]openapi.SavingsGoal,
	error,
) {
	query := `select sg.id,
					   sg.name,
					   sg.description,
					   sg.target_amount,
					   sg.currency,
					   sg.target_date,
					   sg.initial_amount,
					   sg.current_amount,
					   sg.percent_complete,
					   sg.account_id,
					   sg.priority,
					   sg.auto_contribute,
					   sg.auto_contribute_amount,
					   sg.auto_contribute_frequency,
					   sg.status,
					   sg.projected_completion_date,
					   sg.created_at,
					   sg.updated_at,
					   GROUP_CONCAT(sgt.tag_id ORDER BY sgt.tag_id SEPARATOR ',') as tags,
					   c.id,
					   c.name,
					   c.description,
					   c.color,
					   c.background_color,
					   c.active,
					   c.category_type
				from savings_goals sg
						 inner join categories c ON sg.category_id = c.id
						 left join savings_goal_tags sgt ON sg.id = sgt.savings_goal_id
				where sg.status != 'inactive'
				group by sg.id, sg.name, sg.description, sg.target_amount, sg.currency, sg.target_date, sg.initial_amount,
						 sg.current_amount, sg.percent_complete, sg.account_id, sg.priority, sg.auto_contribute,
						 sg.auto_contribute_amount, sg.auto_contribute_frequency, sg.status, sg.projected_completion_date,
						 sg.created_at, sg.updated_at
`
	var savingsGoals []openapi.SavingsGoal
	var ids []string
	tagsByID := make(map[string][]string)
	rows, err := s.db.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to select savings goals: %w",
			err,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var savingsGoal openapi.SavingsGoal
		var tags string
		err := rows.Scan(
			&savingsGoal.Id,
			&savingsGoal.Name,
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
			&savingsGoal.ProjectedCompletionDate,
			&savingsGoal.CreatedAt,
			&savingsGoal.UpdatedAt,
			&tags,
			&savingsGoal.Category.Id,
			&savingsGoal.Category.Name,
			&savingsGoal.Category.Description,
			&savingsGoal.Category.Color,
			&savingsGoal.Category.BackgroundColor,
			&savingsGoal.Category.Active,
		)

		tagsByID[*savingsGoal.Id] = strings.Split(
			tags,
			",",
		)

		if err != nil {
			return nil, fmt.Errorf(
				"failed to scan row: %w",
				err,
			)
		}
		ids = append(
			ids,
			*savingsGoal.Id,
		)
	}
	/*
		tags, err := (*s.tagsRepo).ListByType(
			ctx,
			"savings_goal",
			&ids,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get tags: %w",
				err,
			)
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
		}*/

	return savingsGoals, nil
}

func (s SavingGoalRepoImpl) MarkAsCompleted(
	ctx context.Context,
	id string,
) error {
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
		return fmt.Errorf(
			"failed to mark savings goal as completed: %w",
			err,
		)
	}
	return nil
}

func (s SavingGoalRepoImpl) MarkAsAbandoned(
	ctx context.Context,
	id string,
) error {
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
		return fmt.Errorf(
			"failed to mark savings goal as abandoned: %w",
			err,
		)
	}
	return nil
}

func (s SavingGoalRepoImpl) CreateWithdrawal(
	ctx context.Context,
	withdrawal openapi.SavingsWithdrawalRequest,
	transferID, goalID string,
) (
	string,
	error,
) {
	queryInsert := `insert into savings_withdrawals 
						(savings_goal_id, date, reason, transfer_id) 
					VALUES (?,?,?,?)`
	result, errInsert := s.db.ExecContext(
		ctx,
		queryInsert,
		goalID,
		withdrawal.Date,
		withdrawal.Reason,
		transferID,
	)
	if errInsert != nil {
		return "", fmt.Errorf(
			"failed to insert savings withdrawal: %w",
			errInsert,
		)
	}
	withdrawalID, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf(
			"failed to get last insert ID: %w",
			err,
		)
	}

	return fmt.Sprintf(
		"%d",
		withdrawalID,
	), nil
}

func (s SavingGoalRepoImpl) DeleteWithdrawal(
	ctx context.Context,
	id string,
) error {
	queryDelete := `DELETE FROM savings_withdrawals WHERE id=?`
	_, errDelete := s.db.ExecContext(
		ctx,
		queryDelete,
		id,
	)
	if errDelete != nil {
		return fmt.Errorf(
			"failed to delete savings withdrawal: %w",
			errDelete,
		)
	}
	return nil
}

func (s SavingGoalRepoImpl) GetWithdrawalByID(
	ctx context.Context,
	id string,
) (
	*openapi.SavingsWithdrawal,
	error,
) {
	query := `select sw.id,
					   sw.savings_goal_id,
					   sw.date,
					   sw.reason,
					   sw.created_at,
					   sw.updated_at,
					   tr.destination_amount,
					   tr.destination_account_id,
					   t.description,
					   GROUP_CONCAT(swt.tag_id ORDER BY swt.tag_id SEPARATOR ',') as tags
				from savings_withdrawals sw
						 inner join transfers tr ON sw.transfer_id = tr.id
						 inner join transactions t ON tr.transaction_id = t.id
						 left join savings_withdrawal_tags swt ON sw.id = swt.withdrawal_id
				group by sw.id, sw.savings_goal_id, sw.date, sw.reason, sw.transfer_id, sw.created_at, sw.updated_at,
						 tr.destination_amount, tr.destination_account_id, t.description`
	row := s.db.QueryRowContext(
		ctx,
		query,
		id,
	)
	var withdrawal openapi.SavingsWithdrawal
	var tags string
	err := row.Scan(
		&withdrawal.Id,
		&withdrawal.SavingsGoalId,
		&withdrawal.Date,
		&withdrawal.Reason,
		&withdrawal.CreatedAt,
		&withdrawal.UpdatedAt,
		&withdrawal.Amount,
		&withdrawal.DestinationAccountId,
		&withdrawal.Notes,
		&tags,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to scan row: %w",
			err,
		)
	}

	/*	withdrawal.Tags, err = (*s.tagsRepo).GetByIDs(ctx, strings.Split(tags, ","))
		if err != nil {
			return nil, fmt.Errorf("failed to get tags: %w", err)
		}*/
	return &withdrawal, nil
}

func (s SavingGoalRepoImpl) CreateContribution(
	ctx context.Context,
	contribution openapi.SavingsContributionRequest,
	transferID, goalID string,
) (
	string,
	error,
) {
	queryInsert := `insert into savings_contributions (savings_goal_id, date, transfer_id) VALUES (?,?,?)`
	result, errInsert := s.db.ExecContext(
		ctx,
		queryInsert,
		goalID,
		contribution.Date,
		transferID,
	)
	if errInsert != nil {
		return "", fmt.Errorf(
			"failed to insert savings contribution: %w",
			errInsert,
		)
	}
	contributionID, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf(
			"failed to get last insert ID: %w",
			err,
		)
	}
	return fmt.Sprintf(
		"%d",
		contributionID,
	), nil
}

func (s SavingGoalRepoImpl) DeleteContribution(
	ctx context.Context,
	id string,
) error {
	queryDelete := `DELETE FROM savings_contributions WHERE id=?`
	_, errDelete := s.db.ExecContext(
		ctx,
		queryDelete,
		id,
	)
	if errDelete != nil {
		return fmt.Errorf(
			"failed to delete savings contribution: %w",
			errDelete,
		)
	}
	return nil
}

func (s SavingGoalRepoImpl) GetContributionByID(
	ctx context.Context,
	id string,
) (
	*openapi.SavingsContribution,
	error,
) {
	query := `select sc.id,
					   sc.savings_goal_id,
					   sc.date,
					   sc.created_at,
					   sc.updated_at,
					   tr.destination_amount,
					   tr.source_account_id,
					   GROUP_CONCAT(sct.tag_id ORDER BY sct.tag_id SEPARATOR ',') as tags
				from savings_contributions sc
						 inner join transfers tr on sc.transfer_id = tr.id
						 inner join transactions t on tr.transaction_id = t.id
						 left join savings_contribution_tags sct on sc.id = sct.contribution_id
				where sc.id = ?
				group by sc.id,
						 sc.savings_goal_id,
						 sc.date,
						 sc.created_at,
						 sc.updated_at,
						 tr.destination_amount,
						 tr.source_account_id`
	row := s.db.QueryRowContext(
		ctx,
		query,
		id,
	)
	var contribution openapi.SavingsContribution
	var tags string
	err := row.Scan(
		&contribution.Id,
		&contribution.SavingsGoalId,
		&contribution.Date,
		&contribution.CreatedAt,
		&contribution.UpdatedAt,
		&contribution.Amount,
		&contribution.SourceAccountId,
		&tags,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to scan row: %w",
			err,
		)
	}
	/*	contribution.Tags, err = (*s.tagsRepo).GetByIDs(ctx, strings.Split(tags, ","))
		if err != nil {
			return nil, fmt.Errorf("failed to get tags: %w", err)
		}*/
	return &contribution, nil
}

func (s SavingGoalRepoImpl) ListSavingsTransactions(
	ctx context.Context,
	params openapi.ListSavingsTransactionsParams,
	savingGoalID string,
) (
	*openapi.SavingsTransactionList,
	error,
) {
	querySelect := `select IF(sc.id is not null, tr.source_account_id, tr.destination_account_id) as account_id,
					   tr.destination_amount,
					   tr.created_at,
					   t.description,
					   IF(sc.id is not null, GROUP_CONCAT(sct.tag_id ORDER BY sct.tag_id SEPARATOR ','),
						  GROUP_CONCAT(swt.tag_id ORDER BY swt.tag_id SEPARATOR ','))         as tags,
					   IF(sc.id is not null, 'contribution', 'withdrawal')                    as transaction_type,
					   t.id as transaction_id
				from transfers tr
						 inner join transactions t on t.id = tr.transaction_id
						 left join proletariat_budget.savings_contributions sc on tr.id = sc.transfer_id
						 left join proletariat_budget.savings_withdrawals sw on tr.id = sw.transfer_id
						 left join proletariat_budget.savings_contribution_tags sct on sc.id = sct.contribution_id
						 left join proletariat_budget.savings_withdrawal_tags swt on sw.id = swt.withdrawal_id
				where (sc.savings_goal_id = ? or sw.savings_goal_id = ?)
				`
	queryCount := `select COUNT(*) from transfers tr
						 inner join transactions t on t.id = tr.transaction_id
						 left join proletariat_budget.savings_contributions sc on tr.id = sc.transfer_id
						 left join proletariat_budget.savings_withdrawals sw on tr.id = sw.transfer_id
						 left join proletariat_budget.savings_contribution_tags sct on sc.id = sct.contribution_id
						 left join proletariat_budget.savings_withdrawal_tags swt on sw.id = swt.withdrawal_id
					where (sc.savings_goal_id = ? or sw.savings_goal_id = ?)`

	groupBy := ` group by tr.id, sc.id, sw.id `
	orderBy := ` order by tr.created_at DESC `

	limit := fmt.Sprintf(
		" LIMIT %d OFFSET %d",
		*params.Limit,
		*params.Offset,
	)

	whereClause := make(
		[]string,
		0,
	)
	args := make(
		[]any,
		0,
	)
	args = append(
		args,
		savingGoalID,
		savingGoalID,
	)

	if params.StartDate != nil && params.EndDate != nil {
		whereClause = append(
			whereClause,
			"tr.created_at BETWEEN? AND?",
		)
		args = append(
			args,
			*params.StartDate,
			*params.EndDate,
		)
	} else if params.StartDate != nil {
		whereClause = append(
			whereClause,
			"tr.created_at >=?",
		)
		args = append(
			args,
			*params.StartDate,
		)
	} else if params.EndDate != nil {
		whereClause = append(
			whereClause,
			"tr.created_at <=?",
		)
		args = append(
			args,
			*params.EndDate,
		)
	}

	if params.MinAmount != nil && params.MaxAmount != nil {
		whereClause = append(
			whereClause,
			"tr.destination_amount BETWEEN? AND?",
		)
		args = append(
			args,
			*params.MinAmount,
			*params.MaxAmount,
		)
	} else if params.MinAmount != nil {
		whereClause = append(
			whereClause,
			"tr.destination_amount >=?",
		)
		args = append(
			args,
			*params.MinAmount,
		)
	} else if params.MaxAmount != nil {
		whereClause = append(
			whereClause,
			"tr.destination_amount <=?",
		)
		args = append(
			args,
			*params.MaxAmount,
		)
	}

	if params.Type != nil {
		whereClause = append(
			whereClause,
			"IF(sc.id is not null, 'contribution', 'withdrawal') =?",
		)
		args = append(
			args,
			*params.Type,
		)
	}

	for i, clause := range whereClause {
		if i > 0 {
			querySelect += " AND "
			queryCount += " AND "
		}
		querySelect += clause
		queryCount += clause
	}

	querySelect += groupBy + orderBy + limit

	rows, err := s.db.QueryContext(
		ctx,
		querySelect,
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to select savings transactions: %w",
			err,
		)
	}
	defer rows.Close()
	var transactions []openapi.SavingsTransaction
	for rows.Next() {
		var transaction openapi.SavingsTransaction
		var tags string
		err := rows.Scan(
			&transaction.AccountId,
			&transaction.Amount,
			&transaction.Date,
			&transaction.Description,
			&transaction.Id,
			&tags,
			&transaction.Type,
			&transaction.TransactionId,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to scan row: %w",
				err,
			)
		}
		transaction.SavingsGoalId = savingGoalID
		/*		transaction.Tags, err = (*s.tagsRepo).GetByIDs(ctx, strings.Split(tags, ","))
				if err != nil {
					return nil, fmt.Errorf("failed to get tags: %w", err)
				}*/
		transactions = append(
			transactions,
			transaction,
		)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to iterate over rows: %w",
			err,
		)
	}

	var totalCount int
	err = s.db.QueryRowContext(
		ctx,
		queryCount,
		args...,
	).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to count rows: %w",
			err,
		)
	}

	return &openapi.SavingsTransactionList{
		Metadata: &openapi.ListMetadata{
			Total:  totalCount,
			Limit:  *params.Limit,
			Offset: *params.Offset,
		},
		Transactions: transactions,
	}, nil
}
