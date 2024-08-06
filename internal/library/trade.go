package bot

import (
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/iqbalbaharum/sol-stalker/internal/adapter"
	"github.com/iqbalbaharum/sol-stalker/internal/storage"
	"github.com/iqbalbaharum/sol-stalker/internal/types"
)

func SetTrade(trade *types.Trade) error {
	db, err := adapter.GetMySQLClient()
	if err != nil {
		log.Print("Failed to get initialize mysql instance: %v", err)
		return err
	}

	tradeStorage := storage.NewTradeStorage(db)
	return tradeStorage.SetTrade(trade)
}

func GetTrade(ammId *solana.PublicKey) (*types.Trade, error) {
	db, err := adapter.GetMySQLClient()
	if err != nil {
		return nil, err
	}

	tradeStorage := storage.NewTradeStorage(db)
	return tradeStorage.GetTrade(ammId)
}
