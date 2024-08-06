package coder

import (
	"bytes"
	"encoding/binary"

	"github.com/gagliardetto/solana-go"
)

type MarketStateLayoutV3 struct {
	Unused1                [5]byte
	Unused2                [8]byte
	OwnAddress             solana.PublicKey
	VaultSignerNonce       uint64
	BaseMint               solana.PublicKey
	QuoteMint              solana.PublicKey
	BaseVault              solana.PublicKey
	BaseDepositsTotal      uint64
	BaseFeesAccrued        uint64
	QuoteVault             solana.PublicKey
	QuoteDepositsTotal     uint64
	QuoteFeesAccrued       uint64
	QuoteDustThreshold     uint64
	RequestQueue           solana.PublicKey
	EventQueue             solana.PublicKey
	Bids                   solana.PublicKey
	Asks                   solana.PublicKey
	BaseLotSize            uint64
	QuoteLotSize           uint64
	FeeRateBps             uint64
	ReferrerRebatesAccrued uint64
	Unused3                [7]byte
}

type RaydiumMarketCoder struct{}

func NewRaydiumMarketCoder() *RaydiumMarketCoder {
	return &RaydiumMarketCoder{}
}

// Decode decodes the given byte array into an instruction.
func (coder *RaydiumMarketCoder) RaydiumMarketDecode(data []byte) (MarketStateLayoutV3, error) {
	return decodeRaydiumMarketData(data)
}

func decodeRaydiumMarketData(data []byte) (MarketStateLayoutV3, error) {
	buf := bytes.NewReader(data)
	var state MarketStateLayoutV3

	binary.Read(buf, binary.LittleEndian, &state)

	return state, nil
}
