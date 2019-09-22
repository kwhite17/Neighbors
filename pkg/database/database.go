package database

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"

	"github.com/kwhite17/Neighbors/pkg/assets"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type DbManager interface {
	ReadAllEntities(ctx context.Context) (*sql.Rows, error)
	ReadEntity(ctx context.Context, id interface{}) (*sql.Rows, error)
	WriteEntity(ctx context.Context, values []interface{}) (sql.Result, error)
	DeleteEntity(ctx context.Context, id interface{}) (sql.Result, error)
}

func InitDatabase(dbConfig *dbConfig) *sql.DB {
	db, err := sql.Open(dbConfig.Driver, dbConfig.Host)
	if err != nil {
		log.Fatalf("ERROR - dbInit: Connect - %v\n", err)
	}
	if dbConfig.DevelopmentMode {
		_, err = db.Exec(loadMigration())
		if err != nil {
			log.Fatalf("ERROR - dbInit: Table Creation - %v\n", err)
		}
	}
	if dbConfig.Driver == SQLITE3.Driver {
		db.SetMaxOpenConns(1) //ax this when I switch to production db
	}
	return db
}

func loadMigration() string {
	fullPath := filepath.Join("assets", "scripts", "neighbors_db.sql")
	migrationFile, err := assets.Asset(fullPath)
	if err != nil {
		log.Fatalf("ERROR - migration - %v\n", err)
	}

	return string(migrationFile)
}
