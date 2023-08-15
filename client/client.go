package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	BaseURL              = "http://localhost:8080"
	REQUEST_MAX_DURATION = 300 * time.Millisecond
)

type CurrencyExchange struct {
	Bid string `json:"bid"`
}

var logger *slog.Logger = get_logger()

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), REQUEST_MAX_DURATION)
	defer cancel()
	url := BaseURL + "/cotacao"
	resp, err := doNewRequestWithContext(ctx, url)
	if err != nil {
		logger.ErrorContext(ctx, err.Error())
		return
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error("status code is not 200")
		return
	}

	currencies := CurrencyExchange{}

	err = json.NewDecoder(resp.Body).Decode(&currencies)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	fmt.Printf("USDBRL: %+v\n", currencies)
	err = saveCotacaoToFile(currencies)
	if err != nil {
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

func saveCotacaoToFile(cotacao CurrencyExchange) error {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	stringToWrite := fmt.Sprintf("Dólar: %s\n", cotacao.Bid)
	_, err = file.WriteString(stringToWrite)
	if err != nil {
		return err
	}
	return nil
}

func get_logger() *slog.Logger {
	replace := func(groups []string, a slog.Attr) slog.Attr {
		// Remove the directory from the source's filename.
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
			source.Function = filepath.Base(source.Function)
		}
		return a
	}

	opts := &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: replace,
	}

	var handler slog.Handler = slog.NewJSONHandler(os.Stdout, opts)

	logger := slog.New(handler)
	return logger
}
