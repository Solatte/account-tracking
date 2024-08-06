package liquidity

import (
	"errors"
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/iqbalbaharum/sol-stalker/internal/adapter"
	"github.com/iqbalbaharum/sol-stalker/internal/config"
	"github.com/iqbalbaharum/sol-stalker/internal/rpc"
	"github.com/iqbalbaharum/sol-stalker/internal/storage"
	"github.com/iqbalbaharum/sol-stalker/internal/types"
)

// Return pool keys from storage if available, otherwise fetch from RPC and store in storage
func GetPoolKeys(ammId *solana.PublicKey) (*types.RaydiumPoolKeys, error) {
	redisClient, err := adapter.GetRedisClient(4)

	storedPoolKey, err := storage.GetPoolKeys(redisClient, ammId)

	if err != nil && err.Error() != "key not found" {
		log.Print(err.Error(), err.Error() == "key not found")
		return nil, err
	}

	if !storedPoolKey.ID.IsZero() {
		return storedPoolKey, nil
	}

	state, err := rpc.GetLiquidityState(ammId)
	if err != nil {
		return &types.RaydiumPoolKeys{}, err
	}

	authority, err := getAssociatedAuthority(config.RAYDIUM_AMM_V4)
	if err != nil {
		return &types.RaydiumPoolKeys{}, err
	}

	pKey := &types.RaydiumPoolKeys{
		ID:                 *ammId,
		BaseMint:           state.BaseMint,
		QuoteMint:          state.QuoteMint,
		LpMint:             state.LpMint,
		BaseDecimals:       int(state.BaseDecimal),
		QuoteDecimals:      int(state.QuoteDecimal),
		Authority:          authority,
		OpenOrders:         state.OpenOrders,
		TargetOrders:       state.TargetOrders,
		BaseVault:          state.BaseVault,
		QuoteVault:         state.QuoteVault,
		Version:            3,
		MarketProgramID:    state.MarketProgramId,
		MarketID:           state.MarketId,
		MarketAuthority:    authority,
		WithdrawQueue:      state.WithdrawQueue,
		LpVault:            state.LpVault,
		LookupTableAccount: solana.PublicKey{},
	}

	marketInfo, err := rpc.GetMarketState(&state.MarketId)

	if err != nil {
		return &types.RaydiumPoolKeys{}, err
	}

	pKey.MarketBaseVault = marketInfo.BaseVault
	pKey.MarketQuoteVault = marketInfo.QuoteVault
	pKey.MarketBids = marketInfo.Bids
	pKey.MarketAsks = marketInfo.Asks
	pKey.MarketEventQueue = marketInfo.EventQueue

	storage.SetPoolKeys(redisClient, pKey)

	return pKey, nil
}

func GetMint(pKey *types.RaydiumPoolKeys) (solana.PublicKey, bool, error) {

	var mint solana.PublicKey = solana.PublicKey{}
	var swap = false
	if pKey.BaseMint == config.WRAPPED_SOL {
		mint = pKey.QuoteMint
		swap = true
	} else {
		if pKey.QuoteMint == config.WRAPPED_SOL {
			mint = pKey.BaseMint
		} else {
			return solana.PublicKey{}, false, errors.New("neither BaseMint nor QuoteMint is WRAPPED_SOL")
		}
	}

	return mint, swap, nil
}

func GetPoolSolBalance(pKey *types.RaydiumPoolKeys) (uint64, error) {

	_, swap, err := GetMint(pKey)
	if err != nil {
		return 0, err
	}

	var value uint64
	if !swap {
		value, err = rpc.GetBalance(pKey.QuoteVault)
	} else {
		value, err = rpc.GetBalance(pKey.BaseVault)
	}

	if err != nil {
		return 0, err
	}

	return value, nil
}

func getAssociatedAuthority(programId solana.PublicKey) (solana.PublicKey, error) {
	seed := []byte{97, 109, 109, 32, 97, 117, 116, 104, 111, 114, 105, 116, 121}
	programAddress, _, err := solana.FindProgramAddress([][]byte{seed}, programId)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return programAddress, nil
}
