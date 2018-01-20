package database

import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "log"
import "context"

var NeighborsDatabase = NeighborsDatasource{database: initDatabase()}

type Datasource interface {
	ExecuteReadQuery(ctx context.Context, query string, arguments []interface{}) *sql.Rows
	ExecuteWriteQuery(ctx context.Context, query string, arguments []interface{}) sql.Result
}

type NeighborsDatasource struct {
	database *sql.DB
}

func initDatabase() *sql.DB {
	db, err := sql.Open("mysql", "neighbors_dba:neighbors_dba@/neighbors")
	if err != nil {
		log.Printf("ERROR - dbInit - %v\n", err)
	}
	return db
}

func (nd NeighborsDatasource) ExecuteReadQuery(ctx context.Context, query string, arguments []interface{}) *sql.Rows {
	resultSet, err := nd.database.QueryContext(ctx, query, arguments...)
	if err != nil {
		log.Printf("ERROR - ReadQuery - %v\n", err)
		return nil
	}
	return resultSet
}

func (nd NeighborsDatasource) ExecuteWriteQuery(ctx context.Context, query string, arguments []interface{}) sql.Result {
	result, err := nd.database.ExecContext(ctx, query, arguments...)
	if err != nil {
		log.Printf("ERROR - WriteQuery - %v\n", err)
		return nil
	}
	return result
}
