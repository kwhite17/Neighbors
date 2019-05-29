package shelters

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/kwhite17/Neighbors/pkg/database"
)

var createShelterSessionQuery = "INSERT INTO shelterSessions (SessionKey, ShelterID, LoginTime, LastSeenTime) VALUES (?, ?, ?, ?)"
var deleteShelterSessionQuery = "DELETE FROM shelterSessions WHERE ShelterID=?"
var getShelterSessionQuery = "SELECT SessionKey, ShelterID, LoginTime, LastSeenTime FROM shelterSessions WHERE ShelterID=?"
var updateShelterSessionQuery = "UPDATE shelterSessions SET LoginTime = ?, LastSeenTime = ? WHERE ShelterID = ?"

type ShelterSessionManager struct {
	Datasource database.Datasource
	database.DbManager
}

type ShelterSession struct {
	SessionKey   string
	ShelterID    int64
	LoginTime    int64
	LastSeenTime int64
}

func (sm *ShelterSessionManager) GetShelterSession(ctx context.Context, shelterID interface{}) (*ShelterSession, error) {
	result, err := sm.ReadEntity(ctx, shelterID)
	if err != nil {
		return nil, err
	}
	shelter, err := sm.buildShelterSession(result)
	if err != nil {
		return nil, err
	}
	return shelter[0], nil
}

func (sm *ShelterSessionManager) WriteShelterSession(ctx context.Context, shelterID int64) (string, error) {
	cookieID := strconv.FormatInt(shelterID, 10) + "-" + uuid.New().String()
	currentTime := time.Now().Unix()
	values := []interface{}{cookieID, shelterID, currentTime, currentTime}
	_, err := sm.WriteEntity(ctx, values)
	if err != nil {
		return "", err
	}
	return cookieID, nil
}

func (sm *ShelterSessionManager) UpdateShelterSession(ctx context.Context, shelterID int64, loginTime int64, lastSeenTime int64) error {
	values := []interface{}{loginTime, lastSeenTime, shelterID}
	_, err := sm.Datasource.ExecuteWriteQuery(ctx, updateShelterSessionQuery, values)
	return err
}

func (sm *ShelterSessionManager) DeleteShelterSession(ctx context.Context, shelterID interface{}) (int64, error) {
	result, err := sm.DeleteEntity(ctx, shelterID)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (sm *ShelterSessionManager) ReadEntity(ctx context.Context, id interface{}) (*sql.Rows, error) {
	return sm.Datasource.ExecuteReadQuery(ctx, getShelterSessionQuery, []interface{}{id})
}

func (sm *ShelterSessionManager) ReadEntities(ctx context.Context) (*sql.Rows, error) {
	return nil, fmt.Errorf("ShelterSessionManager does not implement method")
}

func (sm *ShelterSessionManager) WriteEntity(ctx context.Context, values []interface{}) (sql.Result, error) {
	result, err := sm.Datasource.ExecuteWriteQuery(ctx, createShelterSessionQuery, values)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (sm *ShelterSessionManager) DeleteEntity(ctx context.Context, id interface{}) (sql.Result, error) {
	result, err := sm.Datasource.ExecuteWriteQuery(ctx, deleteShelterSessionQuery, []interface{}{id})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (sm *ShelterSessionManager) buildShelterSession(result *sql.Rows) ([]*ShelterSession, error) {
	response := make([]*ShelterSession, 0)
	for result.Next() {
		var sessionKey string
		var shelterID int64
		var loginTime int64
		var lastSeenTime int64
		if err := result.Scan(&sessionKey, &shelterID, &loginTime, &lastSeenTime); err != nil {
			return nil, err
		}
		shelterSession := &ShelterSession{SessionKey: sessionKey, ShelterID: shelterID, LoginTime: loginTime, LastSeenTime: lastSeenTime}
		response = append(response, shelterSession)
	}
	return response, nil
}
