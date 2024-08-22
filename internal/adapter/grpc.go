package adapter

import (
	"fmt"
	"log"

	"github.com/iqbalbaharum/sol-stalker/internal/generators"
	"github.com/iqbalbaharum/sol-stalker/internal/types"
)

var (
	grpcs = make(map[int]*generators.GrpcClient)
)

func InitGrpcsClients(configs []types.GrpcConfig) error {

	var initError error
	for i := 0; i < len(configs); i++ {

		client, err := generators.GrpcConnect(configs[i].Addr, configs[i].Token, configs[i].InsecureConnection)

		if err != nil {
			panic(err)
		}

		grpcs[i] = client
	}

	return initError
}

func GetGrpcsClient(db int) (*generators.GrpcClient, error) {
	grpc, exists := grpcs[db]
	if !exists {
		return nil, fmt.Errorf("grpc client %d is not initialized. call InitGrpcClient first", db)
	}
	return grpc, nil
}

func CloseGrpcsConnection() {
	for idx, grpc := range grpcs {
		if err := grpc.CloseConnection(); err != nil {
			log.Printf("Error closing gRPC %d connection: %v", idx, err)
		}
	}
}
