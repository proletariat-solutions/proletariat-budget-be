package domain

import "errors"

// Savings domain errors
var (
	ErrSavingsGoalNotFound             = errors.New("savings goal not found")
	ErrSavingsGoalHasActiveWithdrawals = errors.New("savings goal has active withdrawals and cannot be deleted")
	ErrSavingsWithdrawalNotFound       = errors.New("savings withdrawal not found")
	ErrSavingsContributionNotFound     = errors.New("savings contribution not found")
)
