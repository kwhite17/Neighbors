package managers

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/kwhite17/Neighbors/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

var createShelterSessionQuery = "INSERT INTO shelterSessions (SessionKey, ShelterID, ShelterName, LoginTime, LastSeenTime) VALUES ($1, $2, $3, $4, $5)"
var deleteShelterSessionQuery = "DELETE FROM shelterSessions WHERE SessionKey=$1"
var getShelterSessionQuery = "SELECT SessionKey, ShelterID, ShelterName, LoginTime, LastSeenTime FROM shelterSessions WHERE SessionKey=$1"
var updateShelterSessionQuery = "UPDATE shelterSessions SET LoginTime = $1, LastSeenTime = $2 WHERE ShelterID = $3"

type ShelterSessionManager struct {
	Datasource database.Datasource
}

type ShelterSession struct {
	SessionKey   string
	ShelterID    int64
	ShelterName  string
	LoginTime    int64
	LastSeenTime int64
}

func (sm *ShelterSessionManager) GetShelterSession(ctx context.Context, sessionKey interface{}) (*ShelterSession, error) {
	row := sm.Datasource.ExecuteSingleReadQuery(ctx, getShelterSessionQuery, []interface{}{sessionKey})
	var key string
	var shelterID int64
	var shelterName string
	var loginTime int64
	var lastSeenTime int64
	if err := row.Scan(&key, &shelterID, &shelterName, &loginTime, &lastSeenTime); err != nil {
		return nil, err
	}
	return &ShelterSession{SessionKey: key, ShelterID: shelterID, ShelterName: shelterName, LoginTime: loginTime, LastSeenTime: lastSeenTime}, nil
}

func (sm *ShelterSessionManager) WriteShelterSession(ctx context.Context, shelterID int64, shelterName string) (string, error) {
	cookieID := strconv.FormatInt(shelterID, 10) + "-" + uuid.New().String()
	currentTime := time.Now().Unix()
	values := []interface{}{cookieID, shelterID, shelterName, currentTime, currentTime}
	_, err := sm.Datasource.ExecuteWriteQuery(ctx, createShelterSessionQuery, values, false)
	if err != nil {
		return "", err
	}
	return cookieID, nil
}

func (sm *ShelterSessionManager) UpdateShelterSession(ctx context.Context, shelterID int64, loginTime int64, lastSeenTime int64) error {
	values := []interface{}{loginTime, lastSeenTime, shelterID}
	_, err := sm.Datasource.ExecuteWriteQuery(ctx, updateShelterSessionQuery, values, false)
	return err
}

func (sm *ShelterSessionManager) DeleteShelterSession(ctx context.Context, sessionKey interface{}) (int64, error) {
	result, err := sm.Datasource.ExecuteWriteQuery(ctx, deleteShelterSessionQuery, []interface{}{sessionKey}, false)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
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

func (sm *ShelterSessionManager) encryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
