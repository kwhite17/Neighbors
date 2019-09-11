package managers

import (
	"context"
	"testing"

	"github.com/kwhite17/Neighbors/pkg/database"
)

var testUsername = "testName"

func initShelterSessionManager() *ShelterSessionManager {
	dbToClose = database.InitDatabase(database.SQLITE3)
	return &ShelterSessionManager{Datasource: database.StandardDatasource{Database: dbToClose}}
}

func TestCanReadItsOwnShelterSessionWrite(t *testing.T) {
	manager := initShelterSessionManager()
	defer cleanDatabase()

	sessionKey, err := manager.WriteShelterSession(context.Background(), testShelterID, testUsername)
	if err != nil {
		t.Error(err)
	}

	createdSession, err := manager.GetShelterSession(context.Background(), sessionKey)
	if err != nil {
		t.Error(err)
	}

	if createdSession.SessionKey != sessionKey {
		t.Errorf("Expected %v to equal %v", createdSession, sessionKey)
	}

	if createdSession.ShelterID != testShelterID {
		t.Errorf("Expected %v to equal %v", createdSession.ShelterID, testShelterID)
	}
}

func TestItCanDeleteShelterSession(t *testing.T) {
	manager := initShelterSessionManager()
	defer cleanDatabase()

	sessionKey, err := manager.WriteShelterSession(context.Background(), testShelterID, testUsername)
	if err != nil {
		t.Error(err)
	}

	rowsDeleted, err := manager.DeleteShelterSession(context.Background(), sessionKey)
	if err != nil {
		t.Error(err)
	}

	if rowsDeleted != 1 {
		t.Error("Expected row to be deleted")
	}
}

func TestCanReadItsOwnSessionUpdate(t *testing.T) {
	manager := initShelterSessionManager()
	defer cleanDatabase()

	sessionKey, err := manager.WriteShelterSession(context.Background(), testShelterID, testUsername)
	if err != nil {
		t.Error(err)
	}

	createdSession, err := manager.GetShelterSession(context.Background(), sessionKey)
	if err != nil {
		t.Error(err)
	}

	if createdSession.SessionKey != sessionKey {
		t.Errorf("Expected %v to equal %v", createdSession, sessionKey)
	}

	if createdSession.ShelterID != testShelterID {
		t.Errorf("Expected %v to equal %v", createdSession.ShelterID, testShelterID)
	}

	updatedLogin := int64(1)
	updatedSeenTime := int64(2)
	err = manager.UpdateShelterSession(context.Background(), testShelterID, updatedLogin, updatedSeenTime)
	if err != nil {
		t.Error(err)
	}

	finalSession, err := manager.GetShelterSession(context.Background(), sessionKey)
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
