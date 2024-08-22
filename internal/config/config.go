package config

import (
	"log"
	"math/rand"
	"os"

	"github.com/gagliardetto/solana-go"
	"github.com/iqbalbaharum/sol-stalker/internal/types"
	"github.com/joho/godotenv"
)

var (
	WRAPPED_SOL       = solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
	RAYDIUM_AMM_V4    = solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")
	OPENBOOK_ID       = solana.MustPublicKeyFromBase58("srmqPvymJeFKQ4zGQed1GFppgkRHL9kaELCbyksJtPX")
	RAYDIUM_AUTHORITY = solana.MustPublicKeyFromBase58("5Q544fKrFoe6tsEbD7S8EmxGTJYAKtTVhAW5Q5pge4j1")
	COMPUTE_PROGRAM   = solana.MustPublicKeyFromBase58("ComputeBudget111111111111111111111111111111")
	TRANSFER_PROGRAM  = solana.MustPublicKeyFromBase58("11111111111111111111111111111111")
	BLOXROUTE_TIP     = solana.MustPublicKeyFromBase58("HWEoBxYs7ssKuudEjzjmpfJVX7Dvi7wescFsVx2L5yoY")
	GRPC1             types.GrpcConfig
	GRPC2             types.GrpcConfig
)

var (
	Grpc1Addr          string
	GrpcToken          string
	Grpc2Addr          string
	InsecureConnection bool
	RpcHttpUrl         string
	RpcWsUrl           string
	RedisAddr          string
	RedisPassword      string
	MySqlDsn           string
	MySqlDbName        string
)

func InitEnv() error {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	Grpc1Addr = os.Getenv("GRPC1_ENDPOINT")
	GrpcToken = os.Getenv("GRPC_TOKEN")
	Grpc2Addr = os.Getenv("GRPC2_ENDPOINT")
	InsecureConnection = os.Getenv("GRPC_INSECURE") == "true"
	RpcHttpUrl = os.Getenv("RPC_HTTP_URL")
	RpcWsUrl = os.Getenv("RPC_WS_URL")
	RedisAddr = os.Getenv("REDIS_ADDR")
	RedisPassword = os.Getenv("REDIS_PASSWORD")
	MySqlDsn = os.Getenv("MYSQL_DSN")
	MySqlDbName = os.Getenv("MYSQL_DBNAME")

	GRPC1 = types.GrpcConfig{
		Addr:               Grpc1Addr,
		Token:              GrpcToken,
		InsecureConnection: false,
	}

	GRPC2 = types.GrpcConfig{
		Addr:               Grpc2Addr,
		Token:              "",
		InsecureConnection: true,
	}

	return nil
}

func GetJitoTipAddress() solana.PublicKey {

	var mainnetTipAccounts = []solana.PublicKey{
		solana.MustPublicKeyFromBase58("96gYZGLnJYVFmbjzopPSU6QiEV5fGqZNyN9nmNhvrZU5"),
		solana.MustPublicKeyFromBase58("HFqU5x63VTqvQss8hp11i4wVV8bD44PvwucfZ2bU7gRe"),
		solana.MustPublicKeyFromBase58("Cw8CFyM9FkoMi7K7Crf6HNQqf4uEMzpKw6QNghXLvLkY"),
		solana.MustPublicKeyFromBase58("ADaUMid9yfUytqMBgopwjb2DTLSokTSzL1zt6iGPaS49"),
		solana.MustPublicKeyFromBase58("DfXygSm4jCyNCybVYYK6DwvWqjKee8pbDmJGcLWNDXjh"),
		solana.MustPublicKeyFromBase58("ADuUkR4vqLUMWXxW9gh6D6L8pMSawimctcNZ5pGwDcEt"),
		solana.MustPublicKeyFromBase58("DttWaMuVvTiduZRnguLF7jNxTgiMBZ1hyAumKUiL2KRL"),
		solana.MustPublicKeyFromBase58("3AVi9Tg9Uo68tJfuvoKvqKNWKkC5wPdSSdeBnizKZ6jT"),
	}

	randomIndex := rand.Intn(len(mainnetTipAccounts))
	return mainnetTipAccounts[randomIndex]
}
