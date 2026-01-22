package database

import (
	"crypto-sync-bot/internal/models"
	"database/sql"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitSQLite(path string) error {
	var err error
	DB, err = sql.Open("sqlite", path)
	if err != nil {
		return err
	}

	query := `
	CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		exchange TEXT,
		symbol TEXT,
		order_id TEXT,
		status TEXT,
		error_message TEXT,
		timestamp INTEGER
	);`

	_, err = DB.Exec(query)
	return err
}

func SaveOrderResult(res *models.OrderResult) error {
	query := `INSERT INTO orders (exchange, symbol, order_id, status, error_message, timestamp) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, res.Exchange, res.Symbol, res.OrderID, res.Status, res.ErrorMessage, res.Timestamp)
	return err
}
