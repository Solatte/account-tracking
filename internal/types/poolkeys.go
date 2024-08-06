package types

import "github.com/gagliardetto/solana-go"

type RaydiumPoolKeys struct {
	ID                 solana.PublicKey
	BaseMint           solana.PublicKey
	QuoteMint          solana.PublicKey
	LpMint             solana.PublicKey
	BaseDecimals       int
	QuoteDecimals      int
	LpDecimals         int
	Version            int
	ProgramID          solana.PublicKey
	Authority          solana.PublicKey
	OpenOrders         solana.PublicKey
	TargetOrders       solana.PublicKey
	BaseVault          solana.PublicKey
	QuoteVault         solana.PublicKey
	WithdrawQueue      solana.PublicKey
	LpVault            solana.PublicKey
	MarketProgramID    solana.PublicKey
	MarketID           solana.PublicKey
	MarketAuthority    solana.PublicKey
	MarketBaseVault    solana.PublicKey
	MarketQuoteVault   solana.PublicKey
	MarketBids         solana.PublicKey
	MarketAsks         solana.PublicKey
	MarketEventQueue   solana.PublicKey
	LookupTableAccount solana.PublicKey
}
