package handler

import (
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/iqbalbaharum/sol-stalker/internal/concurrency"
	"github.com/iqbalbaharum/sol-stalker/internal/config"
	"github.com/iqbalbaharum/sol-stalker/internal/database"
	"github.com/iqbalbaharum/sol-stalker/internal/generators"
	"github.com/iqbalbaharum/sol-stalker/internal/utils"
)

const (
	ErrTimeout = "request timed out"
)

func ListenFor(signer string, txChannel chan generators.GeyserResponse, wg *sync.WaitGroup) {

	wg.Add(2)
	go func(addr string) {
		err := generators.GrpcSubscribeByAddresses(
			config.GrpcToken,
			[]string{addr},
			[]string{},
			true,
			txChannel)

		if err != nil {
			log.Printf("Error in first gRPC subscription: %v", err)
		}

		defer func() {
			wg.Done()
			delete(generators.TxSubscriptionManager["failed"].Subscriptions, addr)
			log.Printf("removed failed tx listener: %s \n", addr)
		}()
	}(signer)

	go func(addr string) {
		err := generators.GrpcSubscribeByAddresses(
			config.GrpcToken,
			[]string{addr},
			[]string{},
			false,
			txChannel,
		)

		if err != nil {
			log.Printf("Error in first gRPC subscription: %v", err)
		}

		defer func() {
			wg.Done()
			delete(generators.TxSubscriptionManager["success"].Subscriptions, addr)
			log.Printf("removed success tx listener: %s \n", addr)
		}()

	}(signer)

}

func AddListener(w http.ResponseWriter, r *http.Request) {
	decoded, err := utils.Decode[database.Listener](r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	result, err := database.AddListener(ctx, &decoded)

	if err != nil {
		select {
		case <-ctx.Done():
			http.Error(w, ErrTimeout, http.StatusGatewayTimeout)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	ListenFor(decoded.Signer, concurrency.TxChannel, &concurrency.SubscribeWg)

	err = utils.Encode(w, r, http.StatusCreated, result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RemoveListener(w http.ResponseWriter, r *http.Request) {
	signer := chi.URLParam(r, "signer")

	ctx := r.Context()
	err := database.RemoveListener(ctx, signer)

	generators.TxSubscriptionManager["success"].Remove(signer)
	generators.TxSubscriptionManager["failed"].Remove(signer)

	if err != nil {
		select {
		case <-ctx.Done():
			http.Error(w, ErrTimeout, http.StatusGatewayTimeout)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetListener(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	listeners, err := database.GetAllListener(ctx)

	if err != nil {
		select {
		case <-ctx.Done():
			http.Error(w, ErrTimeout, http.StatusGatewayTimeout)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = utils.Encode(w, r, http.StatusOK, listeners)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
