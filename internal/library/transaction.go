package bot

import (
	"math/big"

	"github.com/gagliardetto/solana-go"
	"github.com/iqbalbaharum/sol-stalker/internal/config"
	"github.com/iqbalbaharum/sol-stalker/internal/types"
)

func GetBalanceFromTransaction(preTokenBalances, postTokenBalances []types.TxTokenBalance, mint solana.PublicKey) *big.Int {
	var tokenPreAccount, tokenPostAccount *types.TxTokenBalance

	for _, account := range preTokenBalances {
		if account.Mint == mint.String() && account.Owner == config.RAYDIUM_AUTHORITY.String() {
			tokenPreAccount = &account
			break
		}
	}

	for _, account := range postTokenBalances {
		if account.Mint == mint.String() && account.Owner == config.RAYDIUM_AUTHORITY.String() {
			tokenPostAccount = &account
			break
		}
	}

	if tokenPreAccount == nil || tokenPostAccount == nil {
		return big.NewInt(0)
	}

	preAmount, success := new(big.Int).SetString(tokenPreAccount.Amount, 10)
	if !success {
		return big.NewInt(0)
	}

	postAmount, success := new(big.Int).SetString(tokenPostAccount.Amount, 10)
	if !success {
		return big.NewInt(0)
	}

	tokenAmount := new(big.Int).Sub(preAmount, postAmount)
	return tokenAmount
}

func GetPreBalanceFromTransaction(preTokenBalances, postTokenBalances []types.TxTokenBalance, mint solana.PublicKey) *big.Int {
	var tokenPreAccount, tokenPostAccount *types.TxTokenBalance

	for _, account := range preTokenBalances {
		if account.Mint == mint.String() && account.Owner == config.RAYDIUM_AUTHORITY.String() {
			tokenPreAccount = &account
			break
		}
	}

	for _, account := range postTokenBalances {
		if account.Mint == mint.String() && account.Owner == config.RAYDIUM_AUTHORITY.String() {
			tokenPostAccount = &account
			break
		}
	}

	if tokenPreAccount == nil || tokenPostAccount == nil {
		return big.NewInt(0)
	}

	preAmount, success := new(big.Int).SetString(tokenPreAccount.Amount, 10)
	if !success {
		return big.NewInt(0)
	}

	return preAmount
}
