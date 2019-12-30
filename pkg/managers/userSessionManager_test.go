package managers

import (
	"context"
	"testing"

	"github.com/kwhite17/Neighbors/pkg/database"
)

var testUsername = "testName"

func initUserSessionManager() *UserSessionManager {
	dbToClose = database.InitDatabase(database.SQLITE3)
	return &UserSessionManager{Datasource: database.StandardDatasource{Database: dbToClose}}
}

func TestCanReadItsOwnSamaritanSessionWrite(t *testing.T) {
	manager := initUserSessionManager()
	defer cleanDatabase()

	sessionKey, err := manager.WriteUserSession(context.Background(), testShelterID, SAMARITAN)
	if err != nil {
		t.Error(err)
	}

	createdSession, err := manager.GetUserSession(context.Background(), sessionKey)
	if err != nil {
		t.Error(err)
	}

	if createdSession.SessionKey != sessionKey {
		t.Errorf("Expected %v to equal %v", createdSession, sessionKey)
	}

	if createdSession.UserID != testShelterID {
		t.Errorf("Expected %v to equal %v", createdSession.UserID, testShelterID)
	}

	if createdSession.UserType != SAMARITAN {
		t.Errorf("Expected %v to equal %v", createdSession.UserType, SAMARITAN)
	}
}
func TestCanReadItsOwnShelterSessionWrite(t *testing.T) {
	manager := initUserSessionManager()
	defer cleanDatabase()

	sessionKey, err := manager.WriteUserSession(context.Background(), testShelterID, SHELTER)
	if err != nil {
		t.Error(err)
	}

	createdSession, err := manager.GetUserSession(context.Background(), sessionKey)
	if err != nil {
		t.Error(err)
	}

	if createdSession.SessionKey != sessionKey {
		t.Errorf("Expected %v to equal %v", createdSession, sessionKey)
	}

	if createdSession.UserID != testShelterID {
		t.Errorf("Expected %v to equal %v", createdSession.UserID, testShelterID)
	}

	if createdSession.UserType != SHELTER {
		t.Errorf("Expected %v to equal %v", createdSession.UserType, SHELTER)
	}
}

func TestItCanDeleteUserSession(t *testing.T) {
	manager := initUserSessionManager()
	defer cleanDatabase()

	sessionKey, err := manager.WriteUserSession(context.Background(), testShelterID, SHELTER)
	if err != nil {
		t.Error(err)
	}

	rowsDeleted, err := manager.DeleteUserSession(context.Background(), sessionKey)
	if err != nil {
		t.Error(err)
	}

	if rowsDeleted != 1 {
		t.Error("Expected row to be deleted")
	}
}

func TestCanReadItsOwnSessionUpdate(t *testing.T) {
	manager := initUserSessionManager()
	defer cleanDatabase()

	sessionKey, err := manager.WriteUserSession(context.Background(), testShelterID, SHELTER)
	if err != nil {
		t.Error(err)
	}

	createdSession, err := manager.GetUserSession(context.Background(), sessionKey)
	if err != nil {
		t.Error(err)
	}

	if createdSession.SessionKey != sessionKey {
		t.Errorf("Expected %v to equal %v", createdSession, sessionKey)
	}

	if createdSession.UserID != testShelterID {
		t.Errorf("Expected %v to equal %v", createdSession.UserID, testShelterID)
	}

	updatedLogin := int64(1)
	updatedSeenTime := int64(2)
	err = manager.UpdateUserSession(context.Background(), testShelterID, updatedLogin, updatedSeenTime)
	if err != nil {
		t.Error(err)
	}

	finalSession, err := manager.GetUserSession(context.Background(), sessionKey)
	if err != nil {
		t.Error(err)
	}

	if finalSession.LoginTime != updatedLogin {
		t.Errorf("Expected %v to equal %v", finalSession, updatedLogin)
	}

	if finalSession.LastSeenTime != updatedSeenTime {
		t.Errorf("Expected %v to equal %v", finalSession.LastSeenTime, updatedSeenTime)
	}
}
