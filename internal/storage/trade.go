// storage/trade.go
package storage

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/iqbalbaharum/sol-stalker/internal/adapter"
	"github.com/iqbalbaharum/sol-stalker/internal/types"
	"github.com/iqbalbaharum/sol-stalker/internal/utils"
)

type TradeStorage struct {
	client    *sql.DB
	tableName string
}

func NewTradeStorage(db *sql.DB) *TradeStorage {
	tName := adapter.TableName
	return &TradeStorage{client: db, tableName: tName}
}

func (s *TradeStorage) SetTrade(trade *types.Trade) error {

	column := "("
	value := " VALUES ("

	for _, c := range adapter.Column {
		column += fmt.Sprintf("%s,", c.Field)
		value += "?,"
	}

	column = utils.ReplaceLastComma(column, ")")
	value = utils.ReplaceLastComma(value, ")")

	query := fmt.Sprintf(`INSERT INTO %s`, s.tableName) + column + value
	unpacked := utils.UnpackStruct(trade)

	_, err := s.client.Exec(query, unpacked...)
	if err != nil {
		log.Print(err)
		return fmt.Errorf("failed to insert trade: %w", err)
	}
	return nil
}

func (s *TradeStorage) GetTrade(ammId *solana.PublicKey) (*types.Trade, error) {
	query := fmt.Sprintf(`
		SELECT *
		FROM %s 
		WHERE ammId = ?
	`, s.tableName)
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
