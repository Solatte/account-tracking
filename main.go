package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/iqbalbaharum/sol-stalker/internal/adapter"
	"github.com/iqbalbaharum/sol-stalker/internal/concurrency"
	"github.com/iqbalbaharum/sol-stalker/internal/config"
	"github.com/iqbalbaharum/sol-stalker/internal/handler"
	bot "github.com/iqbalbaharum/sol-stalker/internal/library"
	"github.com/iqbalbaharum/sol-stalker/internal/storage"
	"github.com/iqbalbaharum/sol-stalker/internal/types"
)

type Server struct {
	Router *chi.Mux
}

func CreateServer() *Server {
	server := &Server{
		Router: handler.CreateRoutes(),
	}

	return server
}

const (
	PORT = 5000
)

var (
	wg sync.WaitGroup
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

	// Create a worker pool
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for response := range concurrency.TxChannel {
				bot.ProcessResponse(response)
			}
		}()
	}

	bot.JitoTipAccounts = adapter.GetJitoTipAccounts()

	err = adapter.InitGrpcsClients([]types.GrpcConfig{
		config.GRPC1,
		config.GRPC2,
	})

	if err != nil {
		panic(err)
	}

	err = adapter.InitMySQLClient(config.MySqlDsn)

	if err != nil {
		panic(err)
	}

	client, err := adapter.GetMySQLClient()

	if err != nil {
		panic(err)
	}

	storage.Init(client)
	concurrency.SubscribeWg.Wait()

	server := CreateServer()
	port := fmt.Sprintf(":%d", PORT)
	fmt.Printf("server running on port%s \n", port)

	http.ListenAndServe(port, server.Router)

	defer func() { adapter.CloseGrpcsConnection() }()
}
