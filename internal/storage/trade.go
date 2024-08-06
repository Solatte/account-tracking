// storage/trade.go
package storage

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/iqbalbaharum/sol-stalker/internal/types"
)

type TradeStorage struct {
	client *sql.DB
}

func NewTradeStorage(db *sql.DB) *TradeStorage {
	return &TradeStorage{client: db}
}

func (s *TradeStorage) SetTrade(trade *types.Trade) error {
	query := `
		INSERT INTO trade (amm_id, mint, action, compute_limit, compute_price, amount, signature, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.client.Exec(query, trade.AmmId.String(), trade.Mint.String(), trade.Action, trade.ComputeLimit, trade.ComputePrice, trade.Amount, trade.Signature, trade.Timestamp)
	if err != nil {
		log.Print(err)
		return fmt.Errorf("failed to insert trade: %w", err)
	}
	return nil
}

func (s *TradeStorage) GetTrade(ammId *solana.PublicKey) (*types.Trade, error) {
	query := `
		SELECT *
		FROM trade
		WHERE ammId = ?
	`
	row := s.client.QueryRow(query, ammId.String())

	var trade types.Trade
	var mint string

	err := row.Scan(&ammId, &mint, &trade.Action, &trade.ComputeLimit, &trade.ComputePrice, &trade.Amount, &trade.Signature, &trade.Timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No trade found
		}
		return nil, fmt.Errorf("failed to retrieve trade: %w", err)
	}

	trade.AmmId = ammId

	mintPubKey, err := solana.PublicKeyFromBase58(mint)
	if err != nil {
		return nil, fmt.Errorf("invalid mint: %w", err)
	}
	trade.Mint = &mintPubKey

	return &trade, nil
}
