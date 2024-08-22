package handler

import (
	"net/http"

	"github.com/iqbalbaharum/sol-stalker/internal/storage"
	"github.com/iqbalbaharum/sol-stalker/internal/types"
	"github.com/iqbalbaharum/sol-stalker/internal/utils"
)

type tradeHandler struct {
}

func NewTradeHandler() *tradeHandler {
	return &tradeHandler{}
}

func (h *tradeHandler) Get(w http.ResponseWriter, r *http.Request) {
	decoded, err := utils.Decode[types.MySQLFilter](r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	trades, err := storage.Trade.Search(decoded)

	if err != nil {
		select {
		case <-ctx.Done():
			http.Error(w, ErrTimeout, http.StatusGatewayTimeout)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	utils.Encode(w, r, http.StatusOK, trades)
}

func (h *tradeHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := storage.Trade.DeleteAll()

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
