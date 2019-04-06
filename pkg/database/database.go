package database

import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "log"
import "context"

var NeighborsDatabase = NeighborsDatasource{Database: initDatabase()}

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

func initDatabase() *sql.DB {
	db, err := sql.Open("mysql", "neighbors_dba:neighbors_dba@/neighbors")
	if err != nil {
		log.Printf("ERROR - dbInit - %v\n", err)
	}
	return db
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
