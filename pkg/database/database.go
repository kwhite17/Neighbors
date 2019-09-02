package database

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"
	"regexp"

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

type Datasource interface {
	ExecuteReadQuery(ctx context.Context, query string, arguments []interface{}) (*sql.Rows, error)
	ExecuteWriteQuery(ctx context.Context, query string, arguments []interface{}) (sql.Result, error)
	ExecuteSingleReadQuery(ctx context.Context, query string, arguments []interface{}) *sql.Row
}

type NeighborsDatasource struct {
	Database *sql.DB
	Config   *DbConfig
}

func InitDatabase(dbConfig *DbConfig) *sql.DB {
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

func (nd NeighborsDatasource) ExecuteSingleReadQuery(ctx context.Context, query string, arguments []interface{}) *sql.Row {
	finalQuery := query
	if nd.Config.Driver != "postgres" {
		postgresParam := regexp.MustCompile("$[0-9]")
		finalQuery = postgresParam.ReplaceAllLiteralString(query, "?")
	}
	return nd.Database.QueryRowContext(ctx, finalQuery, arguments...)
}

func (nd NeighborsDatasource) ExecuteReadQuery(ctx context.Context, query string, arguments []interface{}) (*sql.Rows, error) {
	finalQuery := query
	if nd.Config.Driver != "postgres" {
		postgresParam := regexp.MustCompile("$[0-9]")
		finalQuery = postgresParam.ReplaceAllLiteralString(query, "?")
	}

	resultSet, err := nd.Database.QueryContext(ctx, finalQuery, arguments...)
	if err != nil {
		log.Printf("ERROR - ReadQuery: %s, Args: %v, Error: %v\n", query, arguments, err)
		return nil, err
	}
	return resultSet, nil
}

func (nd NeighborsDatasource) ExecuteWriteQuery(ctx context.Context, query string, arguments []interface{}) (sql.Result, error) {
	result, err := nd.Database.ExecContext(ctx, query, arguments...)
	if err != nil {
		log.Printf("ERROR - WriteQuery: %s, Args: %v, Error: %v\n", query, arguments, err)
		return nil, err
	}
	return result, nil
}
