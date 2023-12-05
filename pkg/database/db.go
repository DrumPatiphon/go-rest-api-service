package database

import (
	"log"

	"github.com/DrumPatiphon/go-rest-api-service/config"

	_ "github.com/jackc/pgx/v5/stdlib" // manual add postgresQl driver
	"github.com/jmoiron/sqlx"
)

func DbConnect(cfg config.IDbConfig) *sqlx.DB {
	//connect
	db, err := sqlx.Connect("pgx", cfg.Url())
	if err != nil {
		log.Fatalf("connect to db failed: %v", err)
	}
	db.DB.SetMaxOpenConns(cfg.MaxOpenConns())
	return db
}
