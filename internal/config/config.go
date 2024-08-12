package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gagliardetto/solana-go"
	"github.com/iqbalbaharum/sol-stalker/internal/adapter"
	"github.com/joho/godotenv"
)

var (
	WRAPPED_SOL       = solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
	RAYDIUM_AMM_V4    = solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")
	OPENBOOK_ID       = solana.MustPublicKeyFromBase58("srmqPvymJeFKQ4zGQed1GFppgkRHL9kaELCbyksJtPX")
	RAYDIUM_AUTHORITY = solana.MustPublicKeyFromBase58("5Q544fKrFoe6tsEbD7S8EmxGTJYAKtTVhAW5Q5pge4j1")
	COMPUTE_PROGRAM   = solana.MustPublicKeyFromBase58("ComputeBudget111111111111111111111111111111")
	TRANSFER_PROGRAM  = solana.MustPublicKeyFromBase58("11111111111111111111111111111111")
)

var (
	GrpcAddr           string
	GrpcToken          string
	InsecureConnection bool
	RpcHttpUrl         string
	RpcWsUrl           string
	RedisAddr          string
	RedisPassword      string
	MysqlDSN           string
)

func InitEnv() error {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	GrpcAddr = os.Getenv("GRPC_ENDPOINT")
	GrpcToken = os.Getenv("GRPC_TOKEN")
	InsecureConnection = os.Getenv("GRPC_INSECURE") == "true"
	RpcHttpUrl = os.Getenv("RPC_HTTP_URL")
	RpcWsUrl = os.Getenv("RPC_WS_URL")
	RedisAddr = os.Getenv("REDIS_ADDR")
	RedisPassword = os.Getenv("REDIS_PASSWORD")
	MysqlDSN = os.Getenv("MYSQL_DSN")

	err := adapter.InitRedisClients(RedisAddr, RedisPassword)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to initialize Redis clients: %v", err))
	}

	err = adapter.InitMySQLClient(MysqlDSN)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to initialize MySQL clients: %v", err))
	}

	err = adapter.CreateDatabaseAndTable()

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to initialize database: %v", err))
	}

	err = adapter.CreateColumn()

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create column: %v", err))
	}

	return nil
}
