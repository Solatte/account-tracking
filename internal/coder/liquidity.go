package coder

import (
	"bytes"
	"encoding/binary"

	"github.com/gagliardetto/solana-go"
)

type LiquidityState struct {
	Status                 uint64
	Nonce                  uint64
	MaxOrder               uint64
	Depth                  uint64
	BaseDecimal            uint64
	QuoteDecimal           uint64
	State                  uint64
	ResetFlag              uint64
	MinSize                uint64
	VolMaxCutRatio         uint64
	AmountWaveRatio        uint64
	BaseLotSize            uint64
	QuoteLotSize           uint64
	MinPriceMultiplier     uint64
	MaxPriceMultiplier     uint64
	SystemDecimalValue     uint64
	MinSeparateNumerator   uint64
	MinSeparateDenominator uint64
	TradeFeeNumerator      uint64
	TradeFeeDenominator    uint64
	PnlNumerator           uint64
	PnlDenominator         uint64
	SwapFeeNumerator       uint64
	SwapFeeDenominator     uint64
	BaseNeedTakePnl        uint64
	QuoteNeedTakePnl       uint64
	QuoteTotalPnl          uint64
	BaseTotalPnl           uint64
	PoolOpenTime           uint64
	PunishPcAmount         uint64
	PunishCoinAmount       uint64
	OrderbookToInitTime    uint64
	SwapBaseInAmount       [16]byte
	SwapQuoteOutAmount     [16]byte
	SwapBase2QuoteFee      uint64
	SwapQuoteInAmount      [16]byte
	SwapBaseOutAmount      [16]byte
	SwapQuote2BaseFee      uint64
	BaseVault              solana.PublicKey
	QuoteVault             solana.PublicKey
	BaseMint               solana.PublicKey
	QuoteMint              solana.PublicKey
	LpMint                 solana.PublicKey
	OpenOrders             solana.PublicKey
	MarketId               solana.PublicKey
	MarketProgramId        solana.PublicKey
	TargetOrders           solana.PublicKey
	WithdrawQueue          solana.PublicKey
	LpVault                solana.PublicKey
	Owner                  solana.PublicKey
	LpReserve              uint64
	Padding                [3]uint64
}

type RaydiumLiquidityCoder struct{}

func NewRaydiumLiquidityCoder() *RaydiumLiquidityCoder {
	return &RaydiumLiquidityCoder{}
}

// Decode decodes the given byte array into an instruction.
func (coder *RaydiumLiquidityCoder) RaydiumLiquidityDecode(data []byte) (LiquidityState, error) {
	return decodeRaydiumLiqudityData(data)
}

func decodeRaydiumLiqudityData(data []byte) (LiquidityState, error) {
	buf := bytes.NewReader(data)
	var state LiquidityState

	binary.Read(buf, binary.LittleEndian, &state)

	return state, nil
}
