package database

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"

	"github.com/kwhite17/Neighbors/pkg/assets"
	_ "github.com/mattn/go-sqlite3"
)

type DbManager interface {
	ReadAllEntities(ctx context.Context) (*sql.Rows, error)
	ReadEntity(ctx context.Context, id int64) (*sql.Rows, error)
	WriteEntity(ctx context.Context, values []interface{}) (sql.Result, error)
	DeleteEntity(ctx context.Context, id string) (sql.Result, error)
}

type Datasource interface {
	ExecuteReadQuery(ctx context.Context, query string, arguments []interface{}) (*sql.Rows, error)
	ExecuteWriteQuery(ctx context.Context, query string, arguments []interface{}) (sql.Result, error)
}

type NeighborsDatasource struct {
	Database *sql.DB
}

func InitDatabase(host string, developmentMode bool) *sql.DB {
	db, err := sql.Open("sqlite3", host)
	if err != nil {
		log.Fatalf("ERROR - dbInit: Connect - %v\n", err)
	}
	if developmentMode {
		_, err = db.Exec(loadMigration())
		if err != nil {
			log.Fatalf("ERROR - dbInit: Table Creation - %v\n", err)
		}
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

func (nd NeighborsDatasource) ExecuteReadQuery(ctx context.Context, query string, arguments []interface{}) (*sql.Rows, error) {
	resultSet, err := nd.Database.QueryContext(ctx, query, arguments...)
	if err != nil {
		log.Printf("ERROR - ReadQuery - %v\n", err)
		return nil, err
	}
	return resultSet, nil
}

func (nd NeighborsDatasource) ExecuteWriteQuery(ctx context.Context, query string, arguments []interface{}) (sql.Result, error) {
	result, err := nd.Database.ExecContext(ctx, query, arguments...)
	if err != nil {
		log.Printf("ERROR - WriteQuery - %v\n", err)
		return nil, err
	}
	return result, nil
}
