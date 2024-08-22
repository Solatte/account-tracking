package types

import "github.com/gagliardetto/solana-go"

type Trade struct {
	AmmId        *solana.PublicKey `json:"amm_id"`
	Mint         *solana.PublicKey `json:"mint"`
	Action       string            `json:"action"`
	ComputeLimit uint64            `json:"compute_limit"`
	ComputePrice uint64            `json:"compute_price"`
	Amount       string            `json:"amount"`
	Signature    string            `json:"signature"`
	Timestamp    int64             `json:"timestamp"`
	Tip          string            `json:"tip"`
	TipAmount    int64             `json:"tip_amount"`
	Status       string            `json:"status"`
	Signer       string            `json:"signer"`
}
