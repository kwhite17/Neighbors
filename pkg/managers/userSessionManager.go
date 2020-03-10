package managers

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/kwhite17/Neighbors/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

var createUserSessionQuery = "INSERT INTO userSessions (SessionKey, UserID, UserType, LoginTime, LastSeenTime) VALUES ($1, $2, $3, $4, $5)"
var deleteUserSessionQuery = "DELETE FROM userSessions WHERE SessionKey=$1"
var getUserSessionQuery = "SELECT SessionKey, UserID, UserType, LoginTime, LastSeenTime FROM userSessions WHERE SessionKey=$1"
var updateUserSessionQuery = "UPDATE userSessions SET LoginTime = $1, LastSeenTime = $2 WHERE UserID = $3"

type SessionManger interface {
	GetUserSession(ctx context.Context, sessionKey interface{}) (*UserSession, error)
	WriteUserSession(ctx context.Context, userID int64, userType UserType) (string, error)
	UpdateUserSession(ctx context.Context, userID int64, loginTime int64, lastSeenTime int64) error
	DeleteUserSession(ctx context.Context, sessionKey interface{}) (int64, error)
}

type UserSessionManager struct {
	Datasource database.Datasource
}

type UserSession struct {
	SessionKey   string
	UserID       int64
	UserType     UserType
	LoginTime    int64
	LastSeenTime int64
}

func (sm *UserSessionManager) GetUserSession(ctx context.Context, sessionKey interface{}) (*UserSession, error) {
	row := sm.Datasource.ExecuteSingleReadQuery(ctx, getUserSessionQuery, []interface{}{sessionKey})
	var key string
	var userID int64
	var userType UserType
	var loginTime int64
	var lastSeenTime int64
	if err := row.Scan(&key, &userID, &userType, &loginTime, &lastSeenTime); err != nil {
		return nil, err
	}
	return &UserSession{SessionKey: key, UserID: userID, UserType: userType, LoginTime: loginTime, LastSeenTime: lastSeenTime}, nil
}

func (sm *UserSessionManager) WriteUserSession(ctx context.Context, userID int64, userType UserType) (string, error) {
	cookieID := strconv.FormatInt(userID, 10) + "-" + uuid.New().String()
	currentTime := time.Now().Unix()
	values := []interface{}{cookieID, userID, userType, currentTime, currentTime}
	_, err := sm.Datasource.ExecuteWriteQuery(ctx, createUserSessionQuery, values, false)
	if err != nil {
		return "", err
	}
	return cookieID, nil
}

func (sm *UserSessionManager) UpdateUserSession(ctx context.Context, userID int64, loginTime int64, lastSeenTime int64) error {
	values := []interface{}{loginTime, lastSeenTime, userID}
	_, err := sm.Datasource.ExecuteWriteQuery(ctx, updateUserSessionQuery, values, false)
	return err
}

func (sm *UserSessionManager) DeleteUserSession(ctx context.Context, sessionKey interface{}) (int64, error) {
	result, err := sm.Datasource.ExecuteWriteQuery(ctx, deleteUserSessionQuery, []interface{}{sessionKey}, false)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (sm *UserSessionManager) encryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
