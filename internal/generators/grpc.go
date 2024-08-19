package generators

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/iqbalbaharum/sol-stalker/internal/types"
	"github.com/iqbalbaharum/sol-stalker/internal/utils"
	"github.com/mr-tron/base58"
	pb "github.com/rpcpool/yellowstone-grpc/examples/golang/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

var (
	grpcAddr           string
	token              string
	jsonInput          string
	insecureConnection bool
	slots              bool
	blocks             bool
	block_meta         bool
	signature          string
	resub              uint

	accountsFilter              utils.ArrayFlags
	accountOwnersFilter         utils.ArrayFlags
	transactionsAccountsInclude utils.ArrayFlags
	transactionsAccountsExclude utils.ArrayFlags
)

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Minute, // send pings every 10 seconds if there is no activity
	Timeout:             20 * time.Second, // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

type MempoolTxn struct {
	Source               string                 `json:"source"`
	Signature            string                 `json:"signature"`
	AccountKeys          []string               `json:"accountKeys"`
	RecentBlockhash      string                 `json:"recentBlockhash"`
	Instructions         []TxInstruction        `json:"instructions"`
	InnerInstructions    []*pb.InnerInstruction `json:"innerInstructions"`
	AddressTableLookups  []TxAddressTableLookup `json:"addressTableLookups"`
	PreTokenBalances     []types.TxTokenBalance `json:"preTokenBalances"`
	PostTokenBalances    []types.TxTokenBalance `json:"postTokenBalances"`
	ComputeUnitsConsumed uint64                 `json:"computeUnitsConsumed"`
	Slot                 uint64                 `json:"slot"`
	Error                string                 `json:"error"`
}

type TxInstruction struct {
	ProgramIdIndex uint32  `json:"programIdIndex"`
	Accounts       []uint8 `json:"accounts"`
	Data           []byte  `json:"data"`
}

type TxAddressTableLookup struct {
	AccountKey      string  `json:"accountKey"`
	WritableIndexes []uint8 `json:"writableIndexes"`
	ReadonlyIndexes []uint8 `json:"readonlyIndexes"`
}

type GeyserResponse struct {
	MempoolTxns MempoolTxn `json:"mempoolTxns"`
}

var (
	conn   *grpc.ClientConn
	client pb.GeyserClient
)

func GrpcConnect(address string, plaintext bool) {
	var opts []grpc.DialOption
	if plaintext {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		pool, _ := x509.SystemCertPool()
		creds := credentials.NewClientTLSFromCert(pool, "")
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	opts = append(opts, grpc.WithKeepaliveParams(kacp))
	opts = append(opts, grpc.WithInitialWindowSize(100<<20))     // 4 MB
	opts = append(opts, grpc.WithInitialConnWindowSize(100<<20)) // 4 MB
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1<<30)))

	log.Println("Starting grpc client, connecting to", address)
	var err error
	conn, err = grpc.NewClient(address, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	client = pb.NewGeyserClient(conn)
}

func CloseConnection() error {
	if conn != nil {
		return conn.Close()
	}
	return nil
}

func GrpcSubscribeByAddresses(grpcToken string, accountInclude []string, accountExclude []string, failed bool, txChannel chan<- GeyserResponse) error {
	if client == nil {
		return errors.New("GRPC not connected")
	}

	// defer close(txChannel)

	var subscription pb.SubscribeRequest = pb.SubscribeRequest{
		Slots:        make(map[string]*pb.SubscribeRequestFilterSlots),
		Blocks:       make(map[string]*pb.SubscribeRequestFilterBlocks),
		BlocksMeta:   make(map[string]*pb.SubscribeRequestFilterBlocksMeta),
		Accounts:     make(map[string]*pb.SubscribeRequestFilterAccounts),
		Transactions: make(map[string]*pb.SubscribeRequestFilterTransactions),
		Entry:        make(map[string]*pb.SubscribeRequestFilterEntry),
		Commitment:   pb.CommitmentLevel_PROCESSED.Enum(),
	}

	// subscription.Slots["slots"] = &pb.SubscribeRequestFilterSlots{}
	subscription.Blocks = make(map[string]*pb.SubscribeRequestFilterBlocks)
	subscription.BlocksMeta = make(map[string]*pb.SubscribeRequestFilterBlocksMeta)
	subscription.Accounts = make(map[string]*pb.SubscribeRequestFilterAccounts)
	subscription.Transactions = make(map[string]*pb.SubscribeRequestFilterTransactions)

	// Subscribe to generic transaction stream
	if len(accountInclude) > 0 {
		subscription.Transactions[accountInclude[0]] = &pb.SubscribeRequestFilterTransactions{
			Vote:           utils.BoolPointer(false),
			Failed:         utils.BoolPointer(failed),
			AccountInclude: accountInclude,
			AccountExclude: accountExclude,
		}
	}

	subscriptionJson, err := json.Marshal(&subscription)
	if err != nil {
		log.Printf("Failed to marshal subscription request: %v", subscriptionJson)
		return err
	}
	log.Printf("Subscription request: %s", string(subscriptionJson))

	// Set up the subscription request
	ctx := context.Background()
	if grpcToken != "" {
		md := metadata.New(map[string]string{"x-token": grpcToken})
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	stream, err := client.Subscribe(ctx,
		grpc.MaxCallRecvMsgSize(100<<20),
	)

	if err != nil {
		log.Fatalf("%v", err)
		return err
	}

	err = stream.Send(&subscription)
	if err != nil {
		log.Fatalf("%v", err)
		return err
	}

	key := "success"

	if failed {
		key = "failed"
	}

	TxSubscriptionManager[key].Add(accountInclude[0], stream, failed)

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-TxSubscriptionManager[key].Subscriptions[accountInclude[0]].Done:
			stream.CloseSend()
			return nil
		default:
			resp, _ := stream.Recv()

			if err == io.EOF {
				// return
			}

			if err != nil {
				log.Printf("Error occurred in receiving update: %v", err)
				break
			}

			if resp.GetTransaction() != nil {

				message := resp.GetTransaction().Transaction.Transaction.Message
				meta := resp.GetTransaction().Transaction.Meta

				var errorString string

				if meta.Err != nil && len(meta.Err.Err) >= 10 {
					relevantByte := meta.Err.Err[9]
					errorString = fmt.Sprintf("0x%x", relevantByte)
				}

				response := &GeyserResponse{
					MempoolTxns: MempoolTxn{
						Source:               "grpc",
						Signature:            base58.Encode(resp.GetTransaction().Transaction.Signature),
						AccountKeys:          convertAccountKeys(message.AccountKeys),
						RecentBlockhash:      base58.Encode(message.RecentBlockhash),
						Instructions:         convertInstructions(message.Instructions),
						AddressTableLookups:  convertAddressTableLookups(message.AddressTableLookups),
						PreTokenBalances:     convertTokenBalances(meta.PreTokenBalances),
						PostTokenBalances:    convertTokenBalances(meta.PostTokenBalances),
						ComputeUnitsConsumed: *resp.GetTransaction().Transaction.GetMeta().ComputeUnitsConsumed,
						Slot:                 resp.GetTransaction().Slot,
						Error:                errorString,
					},
				}

				txChannel <- *response

			}
		}
	}

	return nil
}

func convertAccountKeys(accountKeys [][]byte) []string {
	encodedKeys := make([]string, len(accountKeys))
	for i, key := range accountKeys {
		encodedKeys[i] = base58.Encode(key)
	}
	return encodedKeys
}

func convertInstructions(instructions []*pb.CompiledInstruction) []TxInstruction {
	convertedInstructions := make([]TxInstruction, len(instructions))
	for i, instr := range instructions {
		convertedInstructions[i] = TxInstruction{
			ProgramIdIndex: instr.ProgramIdIndex,
			Accounts:       instr.Accounts,
			Data:           instr.Data,
		}
	}
	return convertedInstructions
}

func convertAddressTableLookups(lookups []*pb.MessageAddressTableLookup) []TxAddressTableLookup {
	convertedLookups := make([]TxAddressTableLookup, len(lookups))
	for i, lookup := range lookups {
		convertedLookups[i] = TxAddressTableLookup{
			AccountKey:      base58.Encode(lookup.AccountKey),
			WritableIndexes: lookup.WritableIndexes,
			ReadonlyIndexes: lookup.ReadonlyIndexes,
		}
	}
	return convertedLookups
}

func convertTokenBalances(tokenBalances []*pb.TokenBalance) []types.TxTokenBalance {
	convertedBalances := make([]types.TxTokenBalance, len(tokenBalances))
	for i, balance := range tokenBalances {
		convertedBalances[i] = types.TxTokenBalance{
			Mint:    balance.Mint,
			Owner:   balance.Owner,
			Amount:  balance.UiTokenAmount.Amount,
			Decimal: balance.UiTokenAmount.Decimals,
		}
	}
	return convertedBalances
}

func GetBlockhash() (solana.Hash, error) {
	if client == nil {
		return solana.Hash{}, errors.New("GRPC not connected")
	}

	ctx := context.Background()
	block, err := client.GetLatestBlockhash(ctx, &pb.GetLatestBlockhashRequest{
		Commitment: pb.CommitmentLevel_CONFIRMED.Enum(),
	})

	log.Print(err)

	if err != nil {
		return solana.Hash{}, err
	}

	hash, err := solana.HashFromBase58(block.Blockhash)
	if err != nil {
		return solana.Hash{}, err
	}

	return hash, nil
}

var TxSubscriptionManager map[string]*SubscriptionManager

func init() {

	TxSubscriptionManager = make(map[string]*SubscriptionManager)
	TxSubscriptionManager["success"] = NewSubscriptionManager()
	TxSubscriptionManager["failed"] = NewSubscriptionManager()
}
