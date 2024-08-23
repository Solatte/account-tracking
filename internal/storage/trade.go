// storage/trade.go
package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/iqbalbaharum/sol-stalker/internal/types"
	"github.com/iqbalbaharum/sol-stalker/internal/utils"
)

type tradeStorage struct {
	client *sql.DB
}

func NewTradeStorage(client *sql.DB) *tradeStorage {
	return &tradeStorage{client: client}
}

func (s *tradeStorage) Set(trade *types.Trade) error {
	columns := utils.BuildInsertQuery(trade)

	query := fmt.Sprintf(`INSERT INTO %s`, TABLE_NAME_TRADE) + columns
	unpacked := utils.UnpackStruct(trade)

	_, err := s.client.Exec(query, unpacked...)
	if err != nil {
		log.Print(err)
		return fmt.Errorf("failed to insert trade: %w", err)
	}
	return nil
}

func (s *tradeStorage) Search(filter types.MySQLFilter) ([]*types.Trade, error) {
	ctx := context.Background()

	query, values := utils.BuildSearchQuery(TABLE_NAME_TRADE, filter)
	stmt, err := s.client.PrepareContext(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrPrepareStatement, err)
	}

	defer stmt.Close()

	rows, err := s.client.QueryContext(ctx, query, values...)

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

func (s *tradeStorage) List() ([]*types.Trade, error) {
	ctx := context.Background()

	stmt, err := s.client.PrepareContext(ctx, `SELECT * FROM `+TABLE_NAME_TRADE)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrPrepareStatement, err)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)

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

func (s *tradeStorage) DeleteAll() error {
	ctx := context.Background()
	stmt, err := s.client.PrepareContext(ctx, `TRUNCATE `+TABLE_NAME_TRADE)

	if err != nil {
		return fmt.Errorf("%s: %w", ErrPrepareStatement, err)
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx)

	if err != nil {
		return fmt.Errorf("%s: %w", ErrExecuteStatement, err)
	}

	return nil
}
