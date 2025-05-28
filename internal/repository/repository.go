package repository

import (
	"context"
	"database/sql"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/gofiber/fiber/v2/log"
	"nebulaLive/internal/config"
	"nebulaLive/internal/entity/ent"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "nebulaLive/pkg/sqlite"
)

// Client is the Entgo client for database operations.
type Client struct {
	*ent.Client
}

// NewClient creates a new Entgo client.
func NewClient(cfg *config.Config) *Client {
	databaseCfg := cfg.Database

	var (
		drv *entsql.Driver
		db  *sql.DB
		err error
	)
	switch databaseCfg.Type {
	case "sqlite3", "sqlite":
		drv, err = entsql.Open(dialect.SQLite, databaseCfg.Connection)
	case "pgx", "postgres":
		db, err = sql.Open("pgx", databaseCfg.Connection)
		drv = entsql.OpenDB(dialect.Postgres, db)
	default:
		log.Fatalf("Unknown database driver: %s", databaseCfg.Type)
	}
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	client := ent.NewClient(
		ent.Driver(drv),
		ent.Log(func(a ...any) {
			log.Debug(a...)
		}),
	)
	if databaseCfg.Migrate {
		if err := client.Schema.Create(context.Background()); err != nil {
			log.Fatalf("Failed to create database schema: %v", err)
		}
	}
	return &Client{Client: client}
}
