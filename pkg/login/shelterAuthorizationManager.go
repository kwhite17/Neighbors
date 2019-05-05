package login

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kwhite17/Neighbors/pkg/database"
)

var createAuthQuery = "INSERT INTO shelterAuth (ShelterID, Username, Password, AccessID, CreatedAt, ExpiresAt) VALUES (?, ?, ?, ?, ?, ?)"
var deleteAuthQuery = "DELETE FROM shelterAuth WHERE ShelterID=?"
var getSingleAuthQuery = "SELECT ShelterID, Username, AccessID, CreatedAt, ExpiresAt from shelterAuth where ShelterID=?"

type ShelterAuthManager struct {
	ds database.Datasource
	database.DbManager
}

type ShelterAuth struct {
	ShelterID int64
	Username  string
	Password  string
	AccessID  string
	CreatedAt int64
	ExpiresAt int64
}

func (sam *ShelterAuthManager) GetShelter(ctx context.Context, id int64) (*ShelterAuth, error) {
	result, err := sam.ReadEntity(ctx, id)
	if err != nil {
		return nil, err
	}
	auth, err := sam.buildAuth(result)
	if err != nil {
		return nil, err
	}
	return auth[0], nil
}

func (sam *ShelterAuthManager) WriteAuth(ctx context.Context, auth *ShelterAuth) (int64, error) {
	values := []interface{}{auth.ShelterID, auth.Username, auth.Password, auth.AccessID, auth.CreatedAt, auth.ExpiresAt}
	result, err := sam.WriteEntity(ctx, values)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func (sam *ShelterAuthManager) DeleteShelter(ctx context.Context, id string) (int64, error) {
	result, err := sam.DeleteEntity(ctx, id)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (sam *ShelterAuthManager) ReadEntity(ctx context.Context, id int64) (*sql.Rows, error) {
	return sam.ds.ExecuteReadQuery(ctx, getSingleAuthQuery, []interface{}{id})
}

func (sam *ShelterAuthManager) ReadEntities(ctx context.Context) (*sql.Rows, error) {
	return nil, fmt.Errorf("cannot fetch multiple authorizations at once")
}

func (sam *ShelterAuthManager) WriteEntity(ctx context.Context, values []interface{}) (sql.Result, error) {
	result, err := sam.ds.ExecuteWriteQuery(ctx, createAuthQuery, values)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (sam *ShelterAuthManager) DeleteEntity(ctx context.Context, id string) (sql.Result, error) {
	result, err := sam.ds.ExecuteWriteQuery(ctx, deleteAuthQuery, []interface{}{id})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (sam *ShelterAuthManager) buildAuth(result *sql.Rows) ([]*ShelterAuth, error) {
	response := make([]*ShelterAuth, 0)
	for result.Next() {
		var accessID string
		var createdAt int64
		var expiresAt int64
		var shelterID int64
		var username string
		if err := result.Scan(&accessID, &createdAt, &expiresAt, &shelterID, &username); err != nil {
			return nil, err
		}
		auth := ShelterAuth{AccessID: accessID, CreatedAt: createdAt, ExpiresAt: expiresAt, ShelterID: shelterID, Username: username}
		response = append(response, &auth)
	}
	return response, nil
}
