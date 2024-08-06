package types

type TxTokenBalance struct {
	Mint    string `json:"mint"`
	Owner   string `json:"owner"`
	Amount  string `json:"amount"`
	Decimal uint32 `json:"decimal"`
}
