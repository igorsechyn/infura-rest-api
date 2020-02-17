package routes

import (
	"encoding/json"
	"errors"
	"ethereum-api/pkg/app"
	apperrors "ethereum-api/pkg/app/errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func getResponseCode(err error) int {
	if errors.Is(err, apperrors.BadInputErr) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func RegisterGetTransactionByBlockNumberAndIndex(router *chi.Mux, app *app.App) {
	router.Get("/transactions/blockNumber/{blockNumber}/index/{index}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		blockNumber := chi.URLParam(r, "blockNumber")
		indexNumber := chi.URLParam(r, "index")
		transaction, err := app.TransactionsService.GetTransactionByBlockNumberAndIndex(blockNumber, indexNumber)
		if err != nil {
			w.WriteHeader(getResponseCode(err))
			w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, err))) //nolint:errcheck
			return
		}

		payload, err := json.Marshal(transaction)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, err))) //nolint:errcheck
			return
		}

		if transaction == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "transaction not found"}`)) //nolint:errcheck
			return
		}

		w.Write(payload) //nolint:errcheck
	})
}
