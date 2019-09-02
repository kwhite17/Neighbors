package managers

import (
	"context"
	"database/sql"
	"reflect"
	"strconv"
	"testing"

	"github.com/kwhite17/Neighbors/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

var testCity = "testCity"
var testCountry = "testCountry"
var testName = "testName"
var testPostalCode = "testPostalCode"
var testState = "testState"
var testStreet = "testStreet"

var dbToClose *sql.DB

func initShelterManager() *ShelterManager {
	dbToClose = database.InitDatabase(database.SQLITE3)
	return &ShelterManager{Datasource: &database.NeighborsDatasource{Database: dbToClose, Config: database.SQLITE3}}
}

func TestCanReadItsOwnShelterWrite(t *testing.T) {
	manager := initShelterManager()
	defer cleanDatabase()
	testShelter := generateShelter()

	id, err := manager.WriteShelter(context.Background(), testShelter, "password")
	if err != nil {
		t.Error(err)
	}
	testShelter.ID = id

	actualShelter, err := manager.GetShelter(context.Background(), id)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(testShelter, actualShelter) {
		t.Errorf("Expected %v to equal %v", actualShelter, testShelter)
	}
}

func TestItCanDeleteShelter(t *testing.T) {
	manager := initShelterManager()
	defer cleanDatabase()
	testShelter := generateShelter()

	id, err := manager.WriteShelter(context.Background(), testShelter, "password")
	if err != nil {
		t.Error(err)
	}

	rowsDeleted, err := manager.DeleteShelter(context.Background(), strconv.FormatInt(id, 10))
	if err != nil {
		t.Error(err)
	}

	if rowsDeleted != 1 {
		t.Error("Expected row to be deleted")
	}
}

func TestItCanGetAllShelters(t *testing.T) {
	manager := initShelterManager()
	defer cleanDatabase()
	testShelters := make([]*Shelter, 0)
	for i := 0; i < 5; i++ {
		testShelter := generateShelter()

		id, err := manager.WriteShelter(context.Background(), testShelter, "password")
		if err != nil {
			t.Error(err)
		}
		testShelter.ID = id
		testShelters = append(testShelters, testShelter)
	}

	allShelters, err := manager.GetShelters(context.Background())
	if err != nil {
		t.Error(err)
	}

	for _, shelter := range allShelters {
		if !containsShelter(shelter, testShelters) {
			t.Errorf("Expected %v to be in %v \n", shelter, testShelters)
		}
	}
}

func TestCanReadItsOwnShelterUpdate(t *testing.T) {
	manager := initShelterManager()
	defer cleanDatabase()
	testShelter := generateShelter()

	id, err := manager.WriteShelter(context.Background(), testShelter, "password")
	if err != nil {
		t.Error(err)
	}
	testShelter.ID = id

	createdShelter, err := manager.GetShelter(context.Background(), id)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(testShelter, createdShelter) {
		t.Errorf("Expected %v to equal %v", createdShelter, testShelter)
	}

	createdShelter.PostalCode = "02139"
	err = manager.UpdateShelter(context.Background(), createdShelter)
	if err != nil {
		t.Error(err)
	}

	finalShelter, err := manager.GetShelter(context.Background(), id)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(finalShelter, createdShelter) {
		t.Errorf("Expected %v to equal %v", finalShelter, createdShelter)
	}
}

func TestItGetsPasswordForUsername(t *testing.T) {
	manager := initShelterManager()
	defer cleanDatabase()
	testShelter := generateShelter()
	unhashedPassword := "password"

	id, err := manager.WriteShelter(context.Background(), testShelter, unhashedPassword)
	if err != nil {
		t.Error(err)
	}
	testShelter.ID = id

	createdShelter, err := manager.GetPasswordForUsername(context.Background(), testShelter.Name)
	if err != nil {
		t.Error(err)
	}

	if bcrypt.CompareHashAndPassword([]byte(createdShelter.Password), []byte(unhashedPassword)) != nil {
		t.Errorf("Expected %v to equal %v", []byte(createdShelter.Password), []byte(unhashedPassword))
	}
}

func generateShelter() *Shelter {
	contactInfo := &ContactInformation{
		City:       testCity,
		Country:    testCountry,
		Name:       testName,
		PostalCode: testPostalCode,
		State:      testState,
		Street:     testStreet,
	}

	return &Shelter{ContactInformation: contactInfo}
}

func containsShelter(candidateShelter *Shelter, expectedShelters []*Shelter) bool {
	for _, shelter := range expectedShelters {
		if reflect.DeepEqual(candidateShelter, shelter) {
			return true
		}
	}
	return false
}
