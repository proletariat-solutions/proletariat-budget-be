package mysql

import (
	"errors"

	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
)

var ErrUnhandledConstraint = errors.New("unhandled constraint")

// ForeignKeyConstraint represents all foreign key constraint names in the database
type ForeignKeyConstraint string

//goland:noinspection GoCommentStart
const (
	// Account constraints
	FKAccountOwner    ForeignKeyConstraint = "fk_account_owner"
	FKAccountCurrency ForeignKeyConstraint = "fk_account_currency"

	// Transaction constraints
	FKTransactionAccount  ForeignKeyConstraint = "fk_transaction_account"
	FKTransactionCurrency ForeignKeyConstraint = "fk_transaction_currency"

	// Expenditure constraints
	FKExpenditureCategory    ForeignKeyConstraint = "fk_expenditure_category"
	FKExpenditureTransaction ForeignKeyConstraint = "fk_expenditure_transaction"

	// Expenditure tags constraints
	FKExpenditureTagsExpenditureID ForeignKeyConstraint = "fk_expenditure_tags_expenditure_id"
	FKExpenditureTagsTagID         ForeignKeyConstraint = "fk_expenditure_tags_tag_id"

	// Ingress constraints
	FKIngressCategory          ForeignKeyConstraint = "fk_ingress_category"
	FKIngressRecurrencyPattern ForeignKeyConstraint = "fk_ingress_recurrency_pattern"
	FKIngressTransaction       ForeignKeyConstraint = "fk_ingress_transaction"

	// Ingress tags constraints
	FKIngressTagsIngress ForeignKeyConstraint = "fk_ingress"
	FKIngressTagsTag     ForeignKeyConstraint = "fk_tag"

	// Ingress recurrence patterns constraints
	FKIngressRecurrencePatternAccount ForeignKeyConstraint = "fk_ingress_recurrence_pattern_account"

	// Savings goals constraints
	FKSavingsGoalAccount    ForeignKeyConstraint = "fk_savings_goal_account"
	FKSavingsGoalCategoryID ForeignKeyConstraint = "fk_savings_goal_category_id"
	FKSavingsGoalCurrency   ForeignKeyConstraint = "fk_savings_goal_currency"

	// Savings goal tags constraints
	FKSavingsGoalTagsSavingsGoalID ForeignKeyConstraint = "fk_savings_goal_tags_savings_goal_id"
	FKSavingsGoalTagsTagID         ForeignKeyConstraint = "fk_savings_goal_tags_tag_id"

	// Savings contributions constraints
	FKSavingsContributionTransfer ForeignKeyConstraint = "fk_savings_contribution_transfer"

	// Savings contribution tags constraints
	FKSavingsContributionTagSavingContribution ForeignKeyConstraint = "fk_savings_contribution_tag_saving_contribution"
	FKSavingsContributionTagTag                ForeignKeyConstraint = "fk_savings_contribution_tag_tag"

	// Savings withdrawals constraints
	FKSavingsWithdrawalSavingsGoal ForeignKeyConstraint = "fk_savings_withdrawal_savings_goal"
	FKSavingsWithdrawalTransfer    ForeignKeyConstraint = "fk_savings_withdrawal_transfer"

	// Savings withdrawal tags constraints
	FKSavingsWithdrawalTagWithdrawal ForeignKeyConstraint = "fk_savings_withdrawal_tag_withdrawal"
	FKSavingsWithdrawalTagTag        ForeignKeyConstraint = "fk_savings_withdrawal_tag_tag"

	// Exchange rates constraints
	FKExchangeRateCurrencyBase   ForeignKeyConstraint = "fk_exchange_rate_currency_base"
	FKExchangeRateCurrencyTarget ForeignKeyConstraint = "fk_exchange_rate_currency_target"

	// Transfer constraints
	FKTransferSourceAccount       ForeignKeyConstraint = "fk_transfer_source_account"
	FKTransferDestinationAccount  ForeignKeyConstraint = "fk_transfer_destination_account"
	FKTransferOutgoingTransaction ForeignKeyConstraint = "fk_transfer_outgoing_transaction"
	FKTransferIncomingTransaction ForeignKeyConstraint = "fk_transfer_incoming_transaction"

	// User roles constraints
	FKUserRolesUser ForeignKeyConstraint = "fk_user_roles_user"
	FKUserRolesRole ForeignKeyConstraint = "fk_user_roles_role"

	// Transaction rollbacks constraints
	FKTransactionRollbacksRollbackTransaction ForeignKeyConstraint = "fk_transaction_rollbacks_rollback_transaction"
	FKTransactionRollbacksTransaction         ForeignKeyConstraint = "fk_transaction_rollbacks_transaction"
)

// ForeignKeyErrorMap maps foreign key constraints to MySQL error codes and their corresponding domain errors
var ForeignKeyErrorMap = map[ForeignKeyConstraint]map[uint16]error{
	// Account constraints
	FKAccountOwner: {
		1451: domain.ErrMemberHasActiveAccounts,
		1452: domain.ErrMemberNotFound,
	},
	FKAccountCurrency: {
		1451: &port.InfrastructureError{
			Type:    "unknown_constraint_error",
			Message: "depending foreign key violation on 'accounts' table for key 'currency_id'",
			Cause:   ErrUnhandledConstraint,
		}, // shouldn't happen because there's no option to delete a currency, but we'll create a custom error anyways just in case
		1452: domain.ErrInvalidCurrency,
	},

	// Transaction constraints
	FKTransactionAccount: {
		1451: &port.InfrastructureError{
			Type:    "unknown_constraint_error",
			Message: "depending foreign key violation on 'transactions' table for key 'account_id'",
			Cause:   ErrUnhandledConstraint,
		}, // shouldn't happen because there's no way to delete a single TX, but we'll create a custom error anyways just in case
		1452: domain.ErrAccountNotFound,
	},
	FKTransactionCurrency: {
		1451: &port.InfrastructureError{
			Type:    "unknown_constraint_error",
			Message: "depending foreign key violation on 'transactions' table for key 'currency_id'",
			Cause:   ErrUnhandledConstraint,
		}, // shouldn't happen because there's no way to delete a currency, but we'll create a custom error anyways just in case
		1452: domain.ErrInvalidCurrency,
	},

	// Expenditure constraints
	FKExpenditureCategory: {
		1451: &port.InfrastructureError{
			Type:    "unknown_constraint_error",
			Message: "depending foreign key violation on 'expenditures' table for key 'category_id'",
			Cause:   ErrUnhandledConstraint,
		}, // shouldn't happen because expenditures are immutable
		1452: domain.ErrCategoryNotFound,
	},
	FKExpenditureTransaction: {
		1451: &port.InfrastructureError{
			Type:    "unknown_constraint_error",
			Message: "depending foreign key violation on 'expenditures' table for key 'transaction_id'",
			Cause:   ErrUnhandledConstraint,
		}, // shouldn't happen because expenditures are immutable
		1452: domain.ErrTransactionNotFound,
	},

	// Expenditure tags constraints
	FKExpenditureTagsExpenditureID: {
		1451: &port.InfrastructureError{
			Type:    "unknown_constraint_error",
			Message: "depending foreign key violation on 'expenditure_tags' table for key 'expenditure_id'",
			Cause:   ErrUnhandledConstraint,
		}, // shouldn't happen because expenditures are immutable
		1452: domain.ErrExpenditureNotFound,
	},
	FKExpenditureTagsTagID: {
		1451: domain.ErrTagInUse,
		1452: domain.ErrTagNotFound,
	},

	// Ingress constraints
	FKIngressCategory: {
		1451: &port.InfrastructureError{
			Type:    "unknown_constraint_error",
			Message: "depending foreign key violation on 'ingresses' table for key 'category_id'",
			Cause:   ErrUnhandledConstraint,
		}, // shouldn't happen because ingresses are immutable
		1452: domain.ErrCategoryNotFound,
	},
	FKIngressRecurrencyPattern: {
		1451: &port.InfrastructureError{
			Type:    "unknown_constraint_error",
			Message: "depending foreign key violation on 'ingresses' table for key 'recurrency_pattern_id'",
			Cause:   ErrUnhandledConstraint,
		}, // shouldn't happen because ingresses are immutable
		1452: domain.ErrRecurrencePatternNotFound,
	},
	FKIngressTransaction: {
		1451: &port.InfrastructureError{
			Type:    "unknown_constraint_error",
			Message: "depending foreign key violation on 'ingresses' table for key 'transaction_id'",
			Cause:   ErrUnhandledConstraint,
		}, // shouldn't happen because ingresses are immutable
		1452: domain.ErrTransactionNotFound,
	},

	// Ingress tags constraints
	FKIngressTagsIngress: {
		1451: &port.InfrastructureError{
			Type:    "unknown_constraint_error",
			Message: "depending foreign key violation on 'ingress_tags' table for key 'ingress_id'",
			Cause:   ErrUnhandledConstraint,
		}, // shouldn't happen because ingresses are immutable
		1452: domain.ErrIngressNotFound,
	},
	FKIngressTagsTag: {
		1451: domain.ErrTagInUse,
		1452: domain.ErrTagNotFound,
	},

	// Ingress recurrence patterns constraints
	FKIngressRecurrencePatternAccount: {
		1451: domain.ErrAccountHasActiveRecurrencePatterns,
		1452: domain.ErrAccountNotFound,
	},

	// Savings goals constraints
	FKSavingsGoalAccount: {
		1451: domain.ErrAccountHasActiveSavingsGoals,
		1452: domain.ErrAccountNotFound,
	},
	FKSavingsGoalCategoryID: {
		1451: domain.ErrCategoryHasActiveSavingsGoals,
		1452: domain.ErrCategoryNotFound,
	},
	FKSavingsGoalCurrency: {
		1451: &port.InfrastructureError{
			Type:    "unknown_constraint_error",
			Message: "depending foreign key violation on 'savings_goals' table for key 'currency_id'",
			Cause:   ErrUnhandledConstraint,
		}, // shouldn't happen because there's no way to delete a currency
		1452: domain.ErrInvalidCurrency,
	},

	// Savings goal tags constraints
	FKSavingsGoalTagsSavingsGoalID: {
		1451: &port.InfrastructureError{
			Type:    "unknown_constraint_error",
			Message: "depending foreign key violation on 'savings_goal_tags' table for key 'savings_goal_id'",
			Cause:   ErrUnhandledConstraint,
		}, // shouldn't happen because savings goals are immutable
		1452: domain.ErrSavingsGoalNotFound,
	},
	FKSavingsGoalTagsTagID: {
		1451: domain.ErrTagInUse,
		1452: domain.ErrTagNotFound,
	},

	// Savings contributions constraints
	FKSavingsContributionTransfer: {
		1451: &port.InfrastructureError{
			Type:    "unknown_constraint_error",
			Message: "depending foreign key violation on 'savings_contributions' table for key 'transfer_id'",
			Cause:   ErrUnhandledConstraint,
		}, // shouldn't happen because savings contributions are immutable
		1452: domain.ErrTransferNotFound,
	},

	// Savings contribution tags constraints
	FKSavingsContributionTagSavingContribution: {
		1451: &port.InfrastructureError{
			Type:    "unknown_constraint_error",
			Message: "depending foreign key violation on 'savings_contribution_tags' table for key 'contribution_id'",
			Cause:   ErrUnhandledConstraint,
		}, // shouldn't happen because savings contributions are immutable
		1452: domain.ErrSavingsContributionNotFound,
	},
	FKSavingsContributionTagTag: {
		1451: domain.ErrTagInUse,
		1452: domain.ErrTagNotFound,
	},

	// Savings withdrawals constraints
	FKSavingsWithdrawalSavingsGoal: {
		1451: domain.ErrSavingsGoalHasActiveWithdrawals,
		1452: domain.ErrSavingsGoalNotFound,
	},
	FKSavingsWithdrawalTransfer: {
		1451: &port.InfrastructureError{
			Type:    "unknown_constraint_error",
			Message: "depending foreign key violation on 'savings_withdrawals' table for key 'transfer_id'",
			Cause:   ErrUnhandledConstraint,
		}, // shouldn't happen because savings withdrawals are immutable
		1452: domain.ErrTransferNotFound,
	},

	// Savings withdrawal tags constraints
	FKSavingsWithdrawalTagWithdrawal: {
		1451: &port.InfrastructureError{
			Type:    "unknown_constraint_error",
			Message: "depending foreign key violation on 'savings_withdrawal_tags' table for key 'withdrawal_id'",
			Cause:   ErrUnhandledConstraint,
		}, // shouldn't happen because withdrawals are immutable
		1452: domain.ErrSavingsWithdrawalNotFound,
	},

	// Additional constraints from update
	FKSavingsWithdrawalTagTag: {
		1451: domain.ErrTagInUse,
		1452: domain.ErrTagNotFound,
	},
}

// String returns the string representation of the foreign key constraint
func (fk ForeignKeyConstraint) String() string {
	return string(fk)
}
