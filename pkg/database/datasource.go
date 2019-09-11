package database

import (
	"context"
	"database/sql"
	"log"
	"math"
	"regexp"
)

type Datasource interface {
	ExecuteBatchReadQuery(ctx context.Context, query string, arguments []interface{}) (*sql.Rows, error)
	ExecuteWriteQuery(ctx context.Context, query string, arguments []interface{}) (sql.Result, error)
	ExecuteSingleReadQuery(ctx context.Context, query string, arguments []interface{}) *sql.Row
	finalizeQuery(query string, isWriteQuery bool) string
}

type StandardDatasource struct {
	Database *sql.DB
	Datasource
}

type PostgresDatasource struct {
	Database *sql.DB
	Datasource
}

type postgresResult struct {
	lastInsertID int64
	rowsAffected int64
	sql.Result
}

type dbConfig struct {
	Driver          string
	Host            string
	DevelopmentMode bool
}

var SQLITE3 = &dbConfig{
	Driver:          "sqlite3",
	Host:            "file::memory:?mode=memory&cache=shared",
	DevelopmentMode: true,
}

func (sd StandardDatasource) ExecuteSingleReadQuery(ctx context.Context, query string, arguments []interface{}) *sql.Row {
	return sd.Database.QueryRowContext(ctx, sd.finalizeQuery(query, false), arguments...)
}

func (sd StandardDatasource) ExecuteBatchReadQuery(ctx context.Context, query string, arguments []interface{}) (*sql.Rows, error) {
	resultSet, err := sd.Database.QueryContext(ctx, sd.finalizeQuery(query, false), arguments...)
	if err != nil {
		log.Printf("ERROR - ReadQuery: %s, Args: %v, Error: %v\n", query, arguments, err)
		return nil, err
	}
	return resultSet, nil
}

func (sd StandardDatasource) ExecuteWriteQuery(ctx context.Context, query string, arguments []interface{}) (sql.Result, error) {
	result, err := sd.Database.ExecContext(ctx, sd.finalizeQuery(query, true), arguments...)
	if err != nil {
		log.Printf("ERROR - WriteQuery: %s, Args: %v, Error: %v\n", query, arguments, err)
		return nil, err
	}
	return result, nil
}

func (sd StandardDatasource) finalizeQuery(query string, isWriteQuery bool) string {
	postgresParam := regexp.MustCompile("$[0-9]")
	finalQuery := postgresParam.ReplaceAllLiteralString(query, "?")

	return finalQuery
}

func (pd PostgresDatasource) ExecuteSingleReadQuery(ctx context.Context, query string, arguments []interface{}) *sql.Row {
	return pd.Database.QueryRowContext(ctx, pd.finalizeQuery(query, false), arguments...)
}

func (pd PostgresDatasource) ExecuteBatchReadQuery(ctx context.Context, query string, arguments []interface{}) (*sql.Rows, error) {
	resultSet, err := pd.Database.QueryContext(ctx, pd.finalizeQuery(query, false), arguments...)
	if err != nil {
		log.Printf("ERROR - ReadQuery: %s, Args: %v, Error: %v\n", query, arguments, err)
		return nil, err
	}
	return resultSet, nil
}

func (pd PostgresDatasource) ExecuteWriteQuery(ctx context.Context, query string, arguments []interface{}) (sql.Result, error) {
	rows, err := pd.Database.QueryContext(ctx, pd.finalizeQuery(query, true), arguments...)
	if err != nil {
		log.Printf("ERROR - WriteQuery: %s, Args: %v, Error: %v\n", query, arguments, err)
		return nil, err
	}

	return pd.buildPostgresResult(rows), nil
}

func (pd PostgresDatasource) finalizeQuery(query string, isWriteQuery bool) string {
	finalQuery := query
	if isWriteQuery {
		finalQuery = query + " RETURNING id"
	}

	return finalQuery + ";"
}

func (pd PostgresDatasource) buildPostgresResult(rows *sql.Rows) postgresResult {
	lastInsertID := int64(-1)
	rowsAffected := int64(0)
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {

		}
		rowsAffected++
		lastInsertID = int64(math.Max(float64(lastInsertID), float64(id)))
	}
	return postgresResult{lastInsertID: int64(lastInsertID), rowsAffected: int64(rowsAffected)}
}

func (pr postgresResult) LastInsertId() (int64, error) {
	return pr.lastInsertID, nil
}

func (pr postgresResult) RowsAffected() (int64, error) {
	return pr.rowsAffected, nil
}

func BuildDatasource(driver string, host string, developmentMode bool) Datasource {
	config := buildConfig(driver, host, developmentMode)
	if config.Driver == "postgres" {
		return PostgresDatasource{Database: InitDatabase(config)}
	}
	return StandardDatasource{Database: InitDatabase(config)}
}

func buildConfig(driver string, host string, developmentMode bool) *dbConfig {
	return &dbConfig{Driver: driver, Host: host, DevelopmentMode: developmentMode}
}
