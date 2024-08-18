package types

import "github.com/gagliardetto/solana-go"

type Trade struct {
	AmmId        *solana.PublicKey
	Mint         *solana.PublicKey
	Action       string
	ComputeLimit uint64
	ComputePrice uint64
	Amount       string
	Signature    string
	Timestamp    int64
	Tip          string
	TipAmount    int64
	Status       string
}
