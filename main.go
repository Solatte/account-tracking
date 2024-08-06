package main

import (
	"errors"
	"log"
	"math/big"
	"runtime"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/iqbalbaharum/sol-stalker/internal/coder"
	"github.com/iqbalbaharum/sol-stalker/internal/config"
	"github.com/iqbalbaharum/sol-stalker/internal/generators"
	bot "github.com/iqbalbaharum/sol-stalker/internal/library"
	"github.com/iqbalbaharum/sol-stalker/internal/liquidity"
	"github.com/iqbalbaharum/sol-stalker/internal/types"
	pb "github.com/rpcpool/yellowstone-grpc/examples/golang/proto"
)

var (
	client *pb.GeyserClient
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	numCPU := runtime.NumCPU() * 2
	maxProcs := runtime.GOMAXPROCS(0)
	log.Printf("Number of logical CPUs available: %d", numCPU)
	log.Printf("Number of CPUs being used: %d", maxProcs)

	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Printf("Initialized .env")
	err := config.InitEnv()
	if err != nil {
		log.Print(err)
		return
	}

	generators.GrpcConnect(config.GrpcAddr, config.InsecureConnection)

	txChannel := make(chan generators.GeyserResponse)

	var wg sync.WaitGroup

	// Create a worker pool
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for response := range txChannel {
				processResponse(response)
			}
		}()
	}

	var addrs []string = []string{
		"5niysgHXFoa8apmrgeBNRXJ6yPiz4WnMnVnAobUXoaMh",
	}

	log.Print("Tracking ", addrs)

	var subscribeWg sync.WaitGroup
	subscribeWg.Add(2)

	go func() {
		defer subscribeWg.Done()
		err := generators.GrpcSubscribeByAddresses(
			config.GrpcToken,
			addrs,
			[]string{},
			true,
			txChannel)
		if err != nil {
			log.Printf("Error in first gRPC subscription: %v", err)
		}
	}()

	go func() {
		defer subscribeWg.Done()
		err := generators.GrpcSubscribeByAddresses(
			config.GrpcToken,
			addrs,
			[]string{},
			false,
			txChannel)
		if err != nil {
			log.Printf("Error in second gRPC subscription: %v", err)
		}
	}()

	// Wait for both subscriptions to complete
	subscribeWg.Wait()

	wg.Wait()

	defer func() {
		if err := generators.CloseConnection(); err != nil {
			log.Printf("Error closing gRPC connection: %v", err)
		}
	}()
}

func processResponse(response generators.GeyserResponse) {
	c := coder.NewRaydiumAmmInstructionCoder()

	var computeLimit uint32
	var computePrice uint32
	// var tip uint32

	for _, ins := range response.MempoolTxns.Instructions {
		programId := response.MempoolTxns.AccountKeys[ins.ProgramIdIndex]

		if programId == config.RAYDIUM_AMM_V4.String() {
			decodedIx, err := c.Decode(ins.Data)
			if err != nil {
				continue
			}

			switch decodedIx.(type) {
			case coder.SwapBaseIn:
				processSwapBaseIn(ins, response, computeLimit, computePrice)
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
			log.Print(ins.Data)
			log.Print(ins.Accounts)
		}
	}
}

func getPublicKeyFromTx(pos int, tx generators.MempoolTxn, instruction generators.TxInstruction) (*solana.PublicKey, error) {
	accountIndexes := instruction.Accounts
	if len(accountIndexes) == 0 {
		return nil, errors.New("no account indexes provided")
	}

	lookupsForAccountKeyIndex := bot.GenerateTableLookup(tx.AddressTableLookups)
	var ammId *solana.PublicKey
	accountIndex := int(accountIndexes[pos])

	if accountIndex >= len(tx.AccountKeys) {
		lookupIndex := accountIndex - len(tx.AccountKeys)
		lookup := lookupsForAccountKeyIndex[lookupIndex]
		table, err := bot.GetLookupTable(solana.MustPublicKeyFromBase58(lookup.LookupTableKey))
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

func processSwapBaseIn(ins generators.TxInstruction, tx generators.GeyserResponse, computeLimit uint32, computePrice uint32) {
	var ammId *solana.PublicKey

	var err error
	ammId, err = getPublicKeyFromTx(1, tx.MempoolTxns, ins)
	if err != nil {
		return
	}

	if ammId == nil {
		return
	}

	pKey, err := liquidity.GetPoolKeys(ammId)
	if err != nil {
		return
	}

	mint, _, err := liquidity.GetMint(pKey)
	if err != nil {
		return
	}

	preAmount := bot.GetPreBalanceFromTransaction(tx.MempoolTxns.PreTokenBalances, tx.MempoolTxns.PostTokenBalances, mint)
	amount := bot.GetBalanceFromTransaction(tx.MempoolTxns.PreTokenBalances, tx.MempoolTxns.PostTokenBalances, mint)

	var action string = "SELL"

	if amount.Cmp(big.NewInt(0)) != 0 {
		if amount.Sign() == 1 {
			action = "BUY"
		}
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
	}

	err = bot.SetTrade(trade)
	if err != nil {
		log.Print(err)
	}

	log.Printf("%s | %s | %s | %d | %d | %d", ammId, tx.MempoolTxns.Signature, action, computeLimit, computePrice, amount)
}

// XDiryUaLKQ2VNejPFJ7AiZCnNPBd7Sa4thuf85oSXk7
//2024/08/03 01:53:05.972767
//08-02-2024 17:44:13 - 2024/08/03 02:00:03.449094

//
