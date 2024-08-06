package types

import "github.com/gagliardetto/solana-go"

type Tracker struct {
	AmmId       *solana.PublicKey
	Status      string
	LastUpdated int64
}
