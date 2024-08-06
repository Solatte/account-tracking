package liquidity

type LiquidityPoolInfo struct {
	Status       uint64
	BaseDecimals int
	QuoteDecimals int
	LpDecimals   int
	BaseReserve  uint64
	QuoteReserve uint64
	LpSupply     uint64
	StartTime    uint64
}
