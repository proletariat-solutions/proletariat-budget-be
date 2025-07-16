package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"ghorkov32/proletariat-budget-be/openapi"
	"strconv"
)

type TransactionRepoImpl struct {
	db       *sql.DB
	tagsRepo *port.TagsRepo
}

func NewTransactionRepo(
	db *sql.DB,
	tagsRepo port.TagsRepo,
) port.TransactionRepo {
	return &TransactionRepoImpl{db: db, tagsRepo: &tagsRepo}
}

func (t TransactionRepoImpl) Create(
	ctx context.Context,
	transaction domain.Transaction,
) (
	string,
	error,
) {
	queryInsert := `insert into transactions
					(account_id,
					 amount,
					 currency,
					 transaction_date,
					 description,
					 transaction_type,
					 balance_after, 
					 status)
					VALUES (?,
							?,
							?,
							?,
							?,
							?,
							?,
							?)`
	result, errInsert := t.db.ExecContext(
		ctx,
		queryInsert,
		transaction.AccountID,
		transaction.Amount,
		transaction.Currency,
		transaction.TransactionDate,
		transaction.Description,
		transaction.TransactionType,
		transaction.BalanceAfter,
		transaction.Status,
	)
	if errInsert != nil {
		return "", translateError(errInsert)
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return "", translateError(errInsert)
	}
	return strconv.FormatInt(
		lastID,
		10,
	), nil
}

func (t TransactionRepoImpl) GetByID(
	ctx context.Context,
	id string,
) (
	*domain.Transaction,
	error,
) {
	querySelect := `SELECT id, 
					   account_id, 
					   amount, 
					   currency, 
					   transaction_date, 
					   description, 
					   transaction_type, 
					   balance_after, 
					   status 
					FROM transactions WHERE id=?`
	var transaction domain.Transaction
	err := t.db.QueryRowContext(
		ctx,
		querySelect,
		id,
	).Scan(
		&transaction.ID,
		&transaction.AccountID,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.TransactionDate,
		&transaction.Description,
		&transaction.TransactionType,
		&transaction.BalanceAfter,
		&transaction.Status,
	)
	if errors.Is(
		err,
		sql.ErrNoRows,
	) {
		return nil, port.ErrRecordNotFound // TODO return correct error from domain
	}
	if err != nil {
		return nil, translateError(err)
	}
	return &transaction, nil

}

func (t TransactionRepoImpl) List(
	ctx context.Context,
	params openapi.ListTransactionsParams,
) (
	*openapi.TransactionList,
	error,
) {
	/*// Holy shit this was a big fucking query to make.
		querySelect := `WITH transaction_details AS (
					SELECT
						t.id,
						t.balance_after,
						t.currency,
						t.transaction_date,
						t.description,
						t.status,
						t.transaction_type,
						t.account_id,
						t.amount,
								tf.id as transfer_id,
						tf.source_account_id,
						tf.destination_account_id,
						tf.fees,
						tf.incoming_transaction_id,
						tf.outgoing_transaction_id,
								tr.transaction_id as original_transaction_id,
						tr.rollback_reason,
						orig_t.transaction_type as orig_transaction_type,
						orig_t.amount as orig_amount,
						orig_tf.source_account_id as orig_source_account_id,
						orig_tf.destination_account_id as orig_destination_account_id,
						orig_tf.fees as orig_fees,
						orig_tf.incoming_transaction_id as orig_incoming_transaction_id,
						orig_tf.outgoing_transaction_id as orig_outgoing_transaction_id,
								e.id as expenditure_id,
						i.id as ingress_id
					FROM transactions t
							 LEFT JOIN transaction_rollbacks tr ON t.transaction_type = 'rollback' AND tr.rollback_transaction_id = t.id
							 LEFT JOIN transactions orig_t ON tr.transaction_id = orig_t.id
							 LEFT JOIN transfers tf ON t.transaction_type = 'transfer' AND
													   (tf.outgoing_transaction_id = t.id OR tf.incoming_transaction_id = t.id)
							 LEFT JOIN transfers orig_tf ON orig_t.transaction_type = 'transfer' AND
															(orig_tf.outgoing_transaction_id = orig_t.id OR orig_tf.incoming_transaction_id = orig_t.id)
							 LEFT JOIN expenditures e ON t.id = e.transaction_id
							 LEFT JOIN ingresses i ON t.id = i.transaction_id
				),
					 transaction_tags AS (
						 SELECT
							 t.id as transaction_id,
							 GROUP_CONCAT(
									 CASE
										 WHEN t.transaction_type = 'expenditure' THEN et.tag_id
										 WHEN t.transaction_type = 'ingress' THEN it.tag_id
										 END
									 ORDER BY COALESCE(et.tag_id, it.tag_id)
									 SEPARATOR ','
							 ) as tags
						 FROM transactions t
								  LEFT JOIN expenditures e ON t.id = e.transaction_id AND t.transaction_type = 'expenditure'
								  LEFT JOIN ingresses i ON t.id = i.transaction_id AND t.transaction_type = 'ingress'
								  LEFT JOIN expenditure_tags et ON e.id = et.expenditure_id
								  LEFT JOIN ingress_tags it ON i.id = it.ingress_id
						 WHERE t.transaction_type IN ('expenditure', 'ingress')
						 GROUP BY t.id
					 )
				SELECT
					td.balance_after,
						CASE td.transaction_type
						WHEN 'ingress' THEN td.amount
						WHEN 'transfer' THEN
							CASE WHEN td.incoming_transaction_id = td.id THEN td.amount END
						WHEN 'rollback' THEN
							CASE td.orig_transaction_type
								WHEN 'ingress' THEN td.amount
								WHEN 'transfer' THEN
									CASE WHEN td.orig_incoming_transaction_id = td.original_transaction_id THEN td.amount END
								END
						END AS credit,

					td.currency,
					td.transaction_date,

						CASE td.transaction_type
						WHEN 'expenditure' THEN td.amount
						WHEN 'transfer' THEN
							CASE WHEN td.outgoing_transaction_id = td.id THEN td.amount END
						WHEN 'rollback' THEN
							CASE td.orig_transaction_type
								WHEN 'expenditure' THEN td.amount
								WHEN 'transfer' THEN
									CASE WHEN td.orig_outgoing_transaction_id = td.original_transaction_id THEN td.amount END
								END
						END AS debit,

					td.description,

						CASE
						WHEN td.transaction_type = 'transfer' THEN td.fees
						WHEN td.transaction_type = 'rollback' AND td.orig_transaction_type = 'transfer' THEN td.orig_fees
						END AS fees,

						CASE td.transaction_type
						WHEN 'transfer' THEN td.source_account_id
						WHEN 'expenditure' THEN td.account_id
						WHEN 'rollback' AND td.orig_transaction_type = 'transfer' THEN td.orig_source_account_id
						END AS from_account,

					td.id,
					td.status,

						CASE td.transaction_type
						WHEN 'expenditure' THEN td.expenditure_id
						WHEN 'ingress' THEN td.ingress_id
						WHEN 'transfer' THEN td.transfer_id
						END AS related_entity_id,

						CASE td.transaction_type
						WHEN 'transfer' THEN td.destination_account_id
						WHEN 'ingress' THEN td.account_id
						WHEN 'rollback' AND td.orig_transaction_type = 'transfer' THEN td.orig_destination_account_id
						END AS to_account,

					td.transaction_type,
					tt.tags,
					td.original_transaction_id,
					td.rollback_reason

				FROM transaction_details td
						 LEFT JOIN transaction_tags tt ON td.id = tt.transaction_id
				`

		orderBy := " ORDER BY td.transaction_date DESC, td.id DESC"

		queryCount := `WITH transaction_details AS (
	    SELECT
	        t.id,
	        t.balance_after,
	        t.currency,
	        t.transaction_date,
	        t.description,
	        t.status,
	        t.transaction_type,
	        t.account_id,
	        t.amount,
	        tf.id as transfer_id,
	        tf.source_account_id,
	        tf.destination_account_id,
	        tf.fees,
	        tf.incoming_transaction_id,
	        tf.outgoing_transaction_id,
	        tr.transaction_id as original_transaction_id,
	        tr.rollback_reason,
	        orig_t.transaction_type as orig_transaction_type,
	        orig_t.amount as orig_amount,
	        orig_tf.source_account_id as orig_source_account_id,
	        orig_tf.destination_account_id as orig_destination_account_id,
	        orig_tf.fees as orig_fees,
	        orig_tf.incoming_transaction_id as orig_incoming_transaction_id,
	        orig_tf.outgoing_transaction_id as orig_outgoing_transaction_id,
	        e.id as expenditure_id,
	        i.id as ingress_id
	    FROM transactions t
	             LEFT JOIN transaction_rollbacks tr ON t.transaction_type = 'rollback' AND tr.rollback_transaction_id = t.id
	             LEFT JOIN transactions orig_t ON tr.transaction_id = orig_t.id
	             LEFT JOIN transfers tf ON t.transaction_type = 'transfer' AND
	                                       (tf.outgoing_transaction_id = t.id OR tf.incoming_transaction_id = t.id)
	             LEFT JOIN transfers orig_tf ON orig_t.transaction_type = 'transfer' AND
	                                            (orig_tf.outgoing_transaction_id = orig_t.id OR orig_tf.incoming_transaction_id = orig_t.id)
	             LEFT JOIN expenditures e ON t.id = e.transaction_id
	             LEFT JOIN ingresses i ON t.id = i.transaction_id
	),
	     transaction_tags AS (
	         SELECT
	             t.id as transaction_id,
	             GROUP_CONCAT(
	                     CASE
	                         WHEN t.transaction_type = 'expenditure' THEN et.tag_id
	                         WHEN t.transaction_type = 'ingress' THEN it.tag_id
	                         END
	                     ORDER BY COALESCE(et.tag_id, it.tag_id)
	                     SEPARATOR ','
	             ) as tags
	         FROM transactions t
	                  LEFT JOIN expenditures e ON t.id = e.transaction_id AND t.transaction_type = 'expenditure'
	                  LEFT JOIN ingresses i ON t.id = i.transaction_id AND t.transaction_type = 'ingress'
	                  LEFT JOIN expenditure_tags et ON e.id = et.expenditure_id
	                  LEFT JOIN ingress_tags it ON i.id = it.ingress_id
	         WHERE t.transaction_type IN ('expenditure', 'ingress')
	         GROUP BY t.id
	     )
	SELECT COUNT(*) as total_count
	FROM transaction_details td
	         LEFT JOIN transaction_tags tt ON td.id = tt.transaction_id `

		whereClause := make([]string, 0)
		args := make([]any, 0)

		if params.StartDate != nil && params.EndDate != nil {
			whereClause = append(whereClause, "td.transaction_date BETWEEN? AND?")
			args = append(args, *params.StartDate, *params.EndDate)
		} else if params.StartDate != nil {
			whereClause = append(whereClause, "td.transaction_date >=?")
			args = append(args, *params.StartDate)
		} else if params.EndDate != nil {
			whereClause = append(whereClause, "td.transaction_date <=?")
			args = append(args, *params.EndDate)
		}

		if params.TransactionType != nil {
			whereClause = append(whereClause, "td.transaction_type =?")
			args = append(args, *params.TransactionType)
		}

		if params.AccountId != nil {
			whereClause = append(whereClause, "td.account_id = ?")
			args = append(args, params.AccountId)
		}

		if params.MinAmount != nil && params.MaxAmount != nil {
			whereClause = append(whereClause, "td.amount BETWEEN? AND?")
			args = append(args, *params.MinAmount, *params.MaxAmount)
		} else if params.MinAmount != nil {
			whereClause = append(whereClause, "td.amount >=?")
			args = append(args, *params.MinAmount)
		} else if params.MaxAmount != nil {
			whereClause = append(whereClause, "td.amount <=?")
			args = append(args, *params.MaxAmount)
		}

		if params.Currency != nil {
			whereClause = append(whereClause, "td.currency =?")
			args = append(args, *params.Currency)
		}

		if len(whereClause) > 0 {
			querySelect += " WHERE "
			queryCount += " WHERE "
		}
		for i, clause := range whereClause {
			if i > 0 {
				querySelect += " AND "
				queryCount += " AND "
			}
			querySelect += clause
			queryCount += clause
		}
		querySelect += orderBy
		querySelect += fmt.Sprintf(" LIMIT %d OFFSET %d", params.Limit, params.Offset)

		stmt, errQueryStmt := t.db.PrepareContext(ctx, querySelect)
		if errQueryStmt != nil {
			return nil, fmt.Errorf("failed to prepare select statement: %w", errQueryStmt)
		}

		rows, errQueryRows := stmt.QueryContext(ctx, args...)
		if errQueryRows != nil {
			return nil, fmt.Errorf("failed to select transactions: %w", errQueryRows)
		}
		defer rows.Close()
		var transactions []openapi.Transaction
		for rows.Next() {
			var transaction openapi.Transaction
			var tags string
			err := rows.Scan(
				&transaction.BalanceAfter,
				&transaction.Currency,
				&transaction.Date,
				&transaction.Debit,
				&transaction.Description,
				&transaction.Fees,
				&transaction.FromAccountId,
				&transaction.Id,
				&transaction.Status,
				&transaction.RelatedEntityId,
				&transaction.ToAccountId,
				&transaction.TransactionType,
				&tags,
				&transaction.OriginalTransactionId,
				&transaction.RollbackReason,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to scan transaction row: %w", err)
			}

			tagList, err := (*t.tagsRepo).GetByIDs(ctx, strings.Split(tags, ","))
			if err != nil {
				return nil, fmt.Errorf("failed to get expenditure tags: %w", err)
			}
			transaction.Tags = tagList
			break
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("failed to iterate over transaction rows: %w", err)
		}
		stmt, errQueryStmt = t.db.PrepareContext(ctx, queryCount)
		if errQueryStmt != nil {
			return nil, fmt.Errorf("failed to prepare count statement: %w", errQueryStmt)
		}

		var count int
		errQueryCount := stmt.QueryRowContext(ctx, args...).Scan(&count)
		if errQueryCount != nil {
			return nil, fmt.Errorf("failed to count transactions: %w", errQueryCount)
		}
		return &openapi.TransactionList{
			Metadata: openapi.ListMetadata{
				Total:  count,
				Offset: params.Offset,
				Limit:  params.Limit,
			},
			Transactions: transactions,
		}, nil
	*/

	return nil, fmt.Errorf("not implemented") // Placeholder for actual implementation of the function.
}
