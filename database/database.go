package database

import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "log"
import "context"

var neighborsDatabase = initDatabase()

func initDatabase() *sql.DB {
	db, err := sql.Open("mysql", "neighbors_dba:neighbors_dba@/neighbors")
	if err != nil {
		log.Printf("ERROR - dbInit - %v\n", err)
	}
	return db
}

func ExecuteReadQuery(ctx context.Context, query string, arguments []interface{}) *sql.Rows {
	resultSet, err := neighborsDatabase.QueryContext(ctx, query, arguments...)
	if err != nil {
		log.Printf("ERROR - ReadQuery - %v\n", err)
		return nil
	}
	return resultSet
}

func ExecuteWriteQuery(ctx context.Context, query string, arguments []interface{}) sql.Result {
	result, err := neighborsDatabase.ExecContext(ctx, query, arguments...)
	if err != nil {
		log.Printf("ERROR - WriteQuery - %v\n", err)
		return nil
	}
	return result
}
