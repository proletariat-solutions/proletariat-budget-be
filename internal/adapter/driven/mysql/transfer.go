package mysql

import (
	"context"
	"database/sql"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"ghorkov32/proletariat-budget-be/openapi"
	"strconv"
	"strings"
)

type TransferRepoImpl struct {
	db *sql.DB
}

func NewTransferRepo(db *sql.DB) port.TransferRepo {
	return &TransferRepoImpl{db: db}
}

func (t TransferRepoImpl) Create(ctx context.Context, transfer openapi.Transfer, incomingTxID, outgoingTxID string) (string, error) {
	queryInsert := `Insert into transfers (source_account_id,
										   destination_account_id,
										   exchange_rate_multiplier,
										   fees,
										   outgoing_transaction_id,
										   incoming_transaction_id)
					VALUES (?,?,?,?,?,?)`
	result, errInsert := t.db.ExecContext(
		ctx,
		queryInsert,
		transfer.SourceAccountId,
		transfer.DestinationAccountId,
		transfer.DestinationAmount,
		transfer.ExchangeRate,
		transfer.Fees,
		outgoingTxID,
		incomingTxID,
	)
	if errInsert != nil {
		return "", errInsert
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(lastID, 10), nil
}

func (t TransferRepoImpl) GetByID(ctx context.Context, id string) (*openapi.Transfer, error) {
	query := `select tr.created_at,
					   tinc.description,
					   tr.destination_account_id,
					   tinc.amount,
					   tinc.currency,
					   tr.exchange_rate_multiplier,
					   tr.fees,
					   tr.id,
					   tr.source_account_id,
					   tout.amount   AS outgoing_amount,
					   tout.currency AS outgoing_currency,
					   tinc.status
				from transfers tr
						 inner join transactions tinc on tr.incoming_transaction_id = tinc.id
						 inner join transactions tout on tr.outgoing_transaction_id = tout.id
				where tr.id = ?`
	row := t.db.QueryRowContext(ctx, query, id)
	var transfer openapi.Transfer
	err := row.Scan(
		&transfer.Date,
		&transfer.Description,
		&transfer.DestinationAccountId,
		&transfer.DestinationCurrencyId,
		&transfer.DestinationAmount,
		&transfer.ExchangeRate,
		&transfer.Fees,
		&transfer.Id,
		&transfer.SourceAccountId,
		&transfer.SourceAmount,
		&transfer.SourceCurrencyId,
		&transfer.Status,
	)
	if err != nil {
		return nil, err
	}
	return &transfer, nil

}

func (t TransferRepoImpl) List(ctx context.Context, params openapi.ListTransfersParams) (*openapi.TransferList, error) {
	selectQuery := `select tr.created_at,
					   tinc.description,
					   tr.destination_account_id,
					   tinc.amount,
					   tinc.currency,
					   tr.exchange_rate_multiplier,
					   tr.fees,
					   tr.id,
					   tr.source_account_id,
					   tout.amount   AS outgoing_amount,
					   tout.currency AS outgoing_currency,
					   tinc.status
				from transfers tr
						 inner join transactions tinc on tr.incoming_transaction_id = tinc.id
						 inner join transactions tout on tr.outgoing_transaction_id = tout.id`
	countQuery := `SELECT COUNT(*) FROM transfers`

	whereClause := make([]string, 0)
	args := make([]any, 0)

	if params.SourceAccountId != nil {
		whereClause = append(whereClause, "tr.source_account_id =?")
		args = append(args, *params.SourceAccountId)
	}
	if params.DestinationAccountId != nil {
		whereClause = append(whereClause, "tr.destination_account_id =?")
		args = append(args, *params.DestinationAccountId)
	}
	if params.StartDate != nil && params.EndDate != nil {
		whereClause = append(whereClause, "tr.created_at BETWEEN? AND?")
		args = append(args, *params.StartDate, *params.EndDate)
	} else if params.StartDate != nil {
		whereClause = append(whereClause, "tr.created_at >=?")
		args = append(args, *params.StartDate)
	} else if params.EndDate != nil {
		whereClause = append(whereClause, "tr.created_at <=?")
		args = append(args, *params.EndDate)
	}
	limitOffset := ` LIMIT? OFFSET?`
	args = append(args, params.Limit, params.Offset)

	if len(whereClause) > 0 {
		selectQuery += " WHERE " + strings.Join(whereClause, " AND ")
		countQuery += " WHERE " + strings.Join(whereClause, " AND ")
	}
	query := selectQuery + limitOffset
	countQuery += limitOffset

	rows, err := t.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var transfers []openapi.Transfer
	for rows.Next() {
		var transfer openapi.Transfer
		err := rows.Scan(
			&transfer.Date,
			&transfer.Description,
			&transfer.DestinationAccountId,
			&transfer.DestinationAmount,
			&transfer.DestinationCurrencyId,
			&transfer.ExchangeRate,
			&transfer.Fees,
			&transfer.Id,
			&transfer.SourceAccountId,
			&transfer.SourceAmount,
			&transfer.SourceCurrencyId,
			&transfer.Status,
		)
		if err != nil {
			return nil, err
		}
		transfers = append(transfers, transfer)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var totalCount int
	err = t.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, err
	}
	return &openapi.TransferList{
		Transfers: &transfers,
		Metadata: &openapi.ListMetadata{
			Total:  totalCount,
			Limit:  *params.Limit,
			Offset: *params.Offset,
		},
	}, nil
}
