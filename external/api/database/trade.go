package database

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/iqbalbaharum/sol-stalker/internal/types"
)

func SearchTrade(ctx context.Context, filter MySQLFilter) ([]*types.Trade, error) {
	query := `SELECT * FROM trade`
	var values []any

	for idx, q := range filter.Query {
		if idx == 0 {
			query += " WHERE "
		}

		query += fmt.Sprintf("%s %s ?", q.Column, q.Op)
		values = append(values, q.Query)

		if idx < len(filter.Query)-1 {
			query += " AND "
		} else {
			query += ";"
		}
	}

	fmt.Println(query, values)

	s, err := mysqlClient.PrepareContext(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrPrepareStatement, err)
	}

	defer s.Close()

	rows, err := mysqlClient.QueryContext(ctx, query, values...)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrExecuteQuery, err)
	}

	defer rows.Close()

	var trades []*types.Trade

	var ammId string
	var mint string

	for rows.Next() {
		var t types.Trade

		err = rows.Scan(
			&ammId,
			&mint,
			&t.Action,
			&t.ComputeLimit,
			&t.ComputePrice,
			&t.Amount,
			&t.Signature,
			&t.Timestamp,
			&t.Tip,
			&t.TipAmount,
			&t.Status,
			&t.Signer,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", ErrScanData, err)
		}

		ammIdPk, err := solana.PublicKeyFromBase58(ammId)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", ErrScanData, err)
		}

		mintPk, err := solana.PublicKeyFromBase58(mint)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", ErrScanData, err)
		}

		t.AmmId = &ammIdPk
		t.Mint = &mintPk

		trades = append(trades, &t)
	}

	return trades, nil
}

func ListTrade(ctx context.Context) ([]*types.Trade, error) {
	s, err := mysqlClient.PrepareContext(ctx, `SELECT * FROM trade`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrPrepareStatement, err)
	}

	defer s.Close()

	rows, err := s.QueryContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrExecuteQuery, err)
	}

	defer rows.Close()

	var trades []*types.Trade

	for rows.Next() {
		var l types.Trade

		err = rows.Scan(&l.Signer)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", ErrScanData, err)
		}

		trades = append(trades, &l)
	}

	if len(trades) == 0 {
		return nil, ErrTradeNotFound
	}

	return trades, nil
}

func DeleteTrade(ctx context.Context) error {
	s, err := mysqlClient.PrepareContext(ctx, `TRUNCATE trade;`)

	if err != nil {
		return fmt.Errorf("%s: %w", ErrPrepareStatement, err)
	}

	defer s.Close()

	_, err = s.ExecContext(ctx)

	if err != nil {
		return fmt.Errorf("%s: %w", ErrExecuteStatement, err)
	}

	return nil
}
