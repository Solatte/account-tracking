package handler

import (
	"net/http"

	"github.com/iqbalbaharum/sol-stalker/internal/database"
	"github.com/iqbalbaharum/sol-stalker/internal/utils"
)

func GetTrade(w http.ResponseWriter, r *http.Request) {
	decoded, err := utils.Decode[database.MySQLFilter](r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	trades, err := database.SearchTrade(ctx, decoded)

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

func DeleteAllTrade(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := database.DeleteTrade(ctx)

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
