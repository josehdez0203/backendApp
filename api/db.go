package main

import (
	"database/sql"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/josehdez0203/backendApp/logger"
)

func openDB(dns string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dns)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (app *application) connectDB() (*sql.DB, error) {
	connection, err := openDB(app.DNS)
	if err != nil {
		return nil, err
	}
	logger.L_Info("Conn Postgres 👷")
	return connection, nil
}
