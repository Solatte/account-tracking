package bot

import (
	"database/sql"
	"errors"
	"log"
	"math/big"
	"slices"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/iqbalbaharum/sol-stalker/internal/coder"
	"github.com/iqbalbaharum/sol-stalker/internal/config"
	"github.com/iqbalbaharum/sol-stalker/internal/generators"
	"github.com/iqbalbaharum/sol-stalker/internal/liquidity"
	"github.com/iqbalbaharum/sol-stalker/internal/types"
)

func getPublicKeyFromTx(pos int, tx generators.MempoolTxn, instruction generators.TxInstruction) (*solana.PublicKey, error) {
	accountIndexes := instruction.Accounts
	if len(accountIndexes) == 0 {
		return nil, errors.New("no account indexes provided")
	}

	lookupsForAccountKeyIndex := GenerateTableLookup(tx.AddressTableLookups)
	var ammId *solana.PublicKey
	accountIndex := int(accountIndexes[pos])

	if accountIndex >= len(tx.AccountKeys) {
		lookupIndex := accountIndex - len(tx.AccountKeys)
		lookup := lookupsForAccountKeyIndex[lookupIndex]
		table, err := GetLookupTable(solana.MustPublicKeyFromBase58(lookup.LookupTableKey))
		if err != nil {
			return nil, err
		}

		if int(lookup.LookupTableIndex) >= len(table.Addresses) {
			return nil, errors.New("lookup table index out of range")
		}

		ammId = &table.Addresses[lookup.LookupTableIndex]

	} else {
		key := solana.MustPublicKeyFromBase58(tx.AccountKeys[accountIndex])
		ammId = &key
	}

	return ammId, nil
}

func ProcessSwapBaseIn(ins generators.TxInstruction, tx generators.GeyserResponse, computeLimit uint32, computePrice uint32, tip string, tipAmount int64, status string) {
	var ammId *solana.PublicKey

	var err error
	ammId, err = getPublicKeyFromTx(1, tx.MempoolTxns, ins)
	if err != nil {
		return
	}

	if ammId == nil {
		return
	}

	openbookId, err := getPublicKeyFromTx(7, tx.MempoolTxns, ins)

	if err != nil {
		return
	}

	var signerAccountIndex int

	if openbookId.String() == config.OPENBOOK_ID.String() {
		signerAccountIndex = 17
	} else {
		signerAccountIndex = 16
	}

	signerTokenAccount, _ := getPublicKeyFromTx(signerAccountIndex, tx.MempoolTxns, ins)

	pKey, err := liquidity.GetPoolKeys(ammId)
	if err != nil {
		return
	}

	mint, _, err := liquidity.GetMint(pKey)
	if err != nil {
		return
	}

	preAmount := GetPreBalanceFromTransaction(tx.MempoolTxns.PreTokenBalances, tx.MempoolTxns.PostTokenBalances, mint)
	amount := GetBalanceFromTransaction(tx.MempoolTxns.PreTokenBalances, tx.MempoolTxns.PostTokenBalances, mint)

	var action string = "SELL"

	if amount.Cmp(big.NewInt(0)) != 0 {
		if amount.Sign() == 1 {
			action = "BUY"
		}
	}

	if tip == "" {
		tip = sql.NullString{}.String
	}

	trade := &types.Trade{
		AmmId:        ammId,
		Mint:         &mint,
		Action:       action,
		ComputeLimit: uint64(computeLimit),
		ComputePrice: uint64(computePrice),
		Amount:       preAmount.String(),
		Signature:    tx.MempoolTxns.Signature,
		Timestamp:    time.Now().Unix(),
		Tip:          tip,
		TipAmount:    tipAmount,
		Status:       status,
		Signer:       signerTokenAccount.String(),
	}

	err = SetTrade(trade)
	if err != nil {
		log.Print(err)
	}

	log.Printf("%s | %s | %s | %d | %d | %d", ammId, tx.MempoolTxns.Signature, action, computeLimit, computePrice, amount)
}

var (
	JitoTipAccounts     []string
	bloxRouteTipAccount = "HWEoBxYs7ssKuudEjzjmpfJVX7Dvi7wescFsVx2L5yoY"
	tipAccount          = []string{"jito", "bloxroute"}
)

func ProcessResponse(response generators.GeyserResponse) {
	c := coder.NewRaydiumAmmInstructionCoder()

	var (
		isProcess    bool
		ix           generators.TxInstruction
		res          generators.GeyserResponse
		computeLimit uint32
		computePrice uint32
		tipAmount    int64
		tip          string
		status       = "success"
	)

	for _, ins := range response.MempoolTxns.Instructions {
		programId := response.MempoolTxns.AccountKeys[ins.ProgramIdIndex]

		if programId == config.RAYDIUM_AMM_V4.String() {
			decodedIx, err := c.Decode(ins.Data)
			if err != nil {
				continue
			}

			switch decodedIx.(type) {
			case coder.SwapBaseIn:
				isProcess = true
				ix = ins
				res = response
			case coder.SwapBaseOut:
			default:
				log.Println("Unknown instruction type")
			}
		}

		if programId == config.COMPUTE_PROGRAM.String() {
			computeDecoded, err := c.DecodeCompute(ins.Data)
			if err != nil {
				continue
			}

			if computeDecoded.Instruction == 2 {
				computeLimit = computeDecoded.Value
			}

			if computeDecoded.Instruction == 3 {
				computePrice = computeDecoded.Value
			}
		}

		if programId == config.TRANSFER_PROGRAM.String() {
			transfer, err := c.DecodeTransfer(ins.Data)

			if err != nil {
				continue
			}

			var accounts []string

			for _, idx := range ins.Accounts {
				accountKeysLength := len(response.MempoolTxns.AccountKeys)

				if idx >= uint8(accountKeysLength) {
					accounts = append(accounts, response.MempoolTxns.AccountKeys[accountKeysLength-1])
				} else {
					accounts = append(accounts, response.MempoolTxns.AccountKeys[idx])

				}
			}

			destination := accounts[1]

			if destination == "" {
				continue
			}

			isJitoTipAccount := slices.Contains(JitoTipAccounts, destination)
			isBloxRouteTipAccount := strings.EqualFold(destination, bloxRouteTipAccount)

			if isJitoTipAccount {
				tip = tipAccount[0]
				tipAmount = transfer.Amount
				log.Printf("len AccountKeys: %s | %v | %d \n", response.MempoolTxns.Signature, tip, tipAmount)
			} else if isBloxRouteTipAccount {
				tip = tipAccount[1]
				tipAmount = transfer.Amount
				log.Printf("len AccountKeys: %s | %v | %d \n", response.MempoolTxns.Signature, tip, tipAmount)
			}
		}

	}

	if response.MempoolTxns.Error != "" {
		status = "failed"
	}

	if isProcess {
		ProcessSwapBaseIn(ix, res, computeLimit, computePrice, tip, tipAmount, status)
	}
}
