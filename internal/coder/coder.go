package coder

type Initialize2 struct {
	Nonce          byte
	OpenTime       uint64
	InitPcAmount   uint64
	InitCoinAmount uint64
}

type Withdraw struct {
	Amount uint64
}

type SwapBaseIn struct {
	AmountIn         uint64
	MinimumAmountOut uint64
}

type SwapBaseOut struct {
	MaxAmountIn uint64
	AmountOut   uint64
}

type Compute struct {
	Instruction uint8
	Value       uint32
}

type Transfer struct {
	Instruction uint32
	Amount      int64
}
