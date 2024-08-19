package main

import (
	"log"
	"runtime"
	"sync"

	"github.com/iqbalbaharum/sol-stalker/external/api/concurrency"
	"github.com/iqbalbaharum/sol-stalker/internal/adapter"
	"github.com/iqbalbaharum/sol-stalker/internal/config"
	"github.com/iqbalbaharum/sol-stalker/internal/generators"
	bot "github.com/iqbalbaharum/sol-stalker/internal/library"
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
	bot.JitoTipAccounts = adapter.GetJitoTipAccounts()

	txChannel := make(chan generators.GeyserResponse)

	var wg sync.WaitGroup

	// Create a worker pool
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for response := range txChannel {
				bot.ProcessResponse(response)
			}
		}()
	}

	var addrs []string = config.Addresses
	log.Print("Tracking ", addrs)

	for _, addr := range addrs {
		concurrency.SubscribeWg.Add(2)
		go func(addrs []string) {
			defer concurrency.SubscribeWg.Done()
			err := generators.GrpcSubscribeByAddresses(
				config.GrpcToken,
				addrs,
				[]string{},
				true,
				txChannel)
			if err != nil {
				log.Printf("Error in first gRPC subscription: %v", err)
			}
		}([]string{addr})

		go func(addrs []string) {
			defer concurrency.SubscribeWg.Done()
			err := generators.GrpcSubscribeByAddresses(
				config.GrpcToken,
				addrs,
				[]string{},
				false,
				txChannel)
			if err != nil {
				log.Printf("Error in second gRPC subscription: %v", err)
			}
		}([]string{addr})

	}

	// Wait for both subscriptions to complete
	concurrency.SubscribeWg.Wait()

	wg.Wait()

	defer func() {
		if err := generators.CloseConnection(); err != nil {
			log.Printf("Error closing gRPC connection: %v", err)
		}
	}()
}
