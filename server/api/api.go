package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	db "github.com/GabrieldeFreire/client_server_api/server/database"
	log_ "github.com/GabrieldeFreire/client_server_api/server/log"
	schema "github.com/GabrieldeFreire/client_server_api/server/schema"
	"github.com/julienschmidt/httprouter"
)

const REQUEST_MAX_DURATION = 400 * time.Millisecond

var logger *slog.Logger = log_.GetInstance()

func Start(port string) {
	router := httprouter.New()
	router.GET("/cotacao", getCotacaoUSDBRL)
	logger.Info("Starting server on port " + port)
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", port), router))
}

func getCotacaoUSDBRL(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx, cancel := context.WithTimeout(context.Background(), REQUEST_MAX_DURATION)
	defer cancel()

	resp, err := doNewRequestWithContext(ctx, "https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.ErrorContext(ctx, err.Error())
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
		} else {
			err = fmt.Errorf("some unknown error happened: %w", err)
			logger.Error(err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	currencies := map[string]schema.CurrencyExchange{}

	err = json.NewDecoder(resp.Body).Decode(&currencies)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = db.AddCurrency(currencies["USDBRL"])
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	logger.Info("Currency added to database", "currency", currencies["USDBRL"])

	err = json.NewEncoder(w).Encode(currencies["USDBRL"])
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Error(err.Error())
		return
	}
}

func doNewRequestWithContext(ctx context.Context, url string) (*http.Response, error) {
	res, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(res)
	if err != nil {
		return nil, err
	}
	return resp, err
}
