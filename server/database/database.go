package database

import (
	"context"
	"database/sql"
	"log/slog"
	"strconv"
	"time"

	log "github.com/GabrieldeFreire/client_server_api/server/log"
	schema "github.com/GabrieldeFreire/client_server_api/server/schema"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DATABASE_TIMEOUT = 10 * time.Millisecond
	COTACOES_DB      = "./cotacoes.db"
)

var logger *slog.Logger = log.GetInstance()

func AddCurrency(currency schema.CurrencyExchange) error {
	err := createTable()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	db, err := sql.Open("sqlite3", COTACOES_DB)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	currencyBid, err := strconv.ParseFloat(currency.Bid, 32)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), DATABASE_TIMEOUT)
	defer cancel()

	insertQuery := `Insert into cotacoes (id, code, codein, bid) values(?, ?, ?, ?)`
	_, err = db.ExecContext(ctx, insertQuery, uuid.New().String(), currency.Code, currency.Codein, currencyBid)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "db.ExecContext")
		return err
	}
	return nil
}

func createTable() error {
	db, err := sql.Open("sqlite3", COTACOES_DB)
	defer db.Close()
	if err != nil {
		return err
	}

	createCotacoesTableQuery := `
	CREATE TABLE IF NOT EXISTS cotacoes (
		id varchar(36) PRIMARY KEY,
		code text DEFAULT NULL,
		codein text DEFAULT NULL,
		bid REAL DEFAULT NULL
	);
	`

	_, err = db.Exec(createCotacoesTableQuery)

	if err != nil {
		return err
	}

	return nil
}
