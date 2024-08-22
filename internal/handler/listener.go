package handler

import (
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/iqbalbaharum/sol-stalker/internal/adapter"
	"github.com/iqbalbaharum/sol-stalker/internal/concurrency"
	"github.com/iqbalbaharum/sol-stalker/internal/generators"
	"github.com/iqbalbaharum/sol-stalker/internal/storage"
	sub "github.com/iqbalbaharum/sol-stalker/internal/subscription"
	"github.com/iqbalbaharum/sol-stalker/internal/types"
	"github.com/iqbalbaharum/sol-stalker/internal/utils"
)

type listenerHandler struct {
	client *generators.GrpcClient
}

func NewListenerHandler() *listenerHandler {
	client, _ := adapter.GetGrpcsClient(1)
	return &listenerHandler{client}
}

func (h *listenerHandler) ListenFor(signer string, txChannel chan generators.GeyserResponse, wg *sync.WaitGroup) {

	wg.Add(2)
	go func(addr string) {
		err := h.client.GrpcSubscribeByAddresses(
			"solana-tracker",
			[]string{addr},
			[]string{},
			true,
			txChannel)

		if err != nil {
			log.Printf("Error in first gRPC subscription: %v", err)
		}

		defer func() {
			wg.Done()
			delete(sub.TxManager["failed"].Subscriptions, addr)
			log.Printf("removed failed tx listener: %s \n", addr)
		}()
	}(signer)

	go func(addr string) {
		err := h.client.GrpcSubscribeByAddresses(
			"solana-tracker",
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
			delete(sub.TxManager["success"].Subscriptions, addr)
			log.Printf("removed success tx listener: %s \n", addr)
		}()

	}(signer)

}

func (h *listenerHandler) Add(w http.ResponseWriter, r *http.Request) {
	decoded, err := utils.Decode[types.Listener](r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	result, err := storage.Listener.Add(&decoded)

	if err != nil {
		select {
		case <-ctx.Done():
			http.Error(w, ErrTimeout, http.StatusGatewayTimeout)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	h.ListenFor(decoded.Signer, concurrency.TxChannel, &concurrency.SubscribeWg)

	err = utils.Encode(w, r, http.StatusCreated, result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *listenerHandler) Remove(w http.ResponseWriter, r *http.Request) {
	signer := chi.URLParam(r, "signer")

	ctx := r.Context()
	err := storage.Listener.Remove(signer)

	sub.TxManager[sub.SUCCESS].Remove(signer)
	sub.TxManager[sub.FAILED].Remove(signer)

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

func (h *listenerHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	listeners, err := storage.Listener.GetAll()

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
