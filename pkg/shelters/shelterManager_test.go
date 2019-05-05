package shelters

import (
	"context"
	"database/sql/driver"
	"reflect"
	"strconv"
	"testing"

	"github.com/kwhite17/Neighbors/pkg/database"

	"github.com/DATA-DOG/go-sqlmock"
)

var testCity = "testCity"
var testCountry = "testCountry"
var testName = "testName"
var testPostalCode = "testPostalCode"
var testState = "testState"
var testStreet = "testStreet"

func TestCanReadItsOwnWrite(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	testShelter := generateShelter()

	mock.ExpectExec(createShelterQuery).WithArgs(shelterToRow(testShelter)...).WillReturnResult(sqlmock.NewResult(1, 1))
	manager := &ShelterManager{Datasource: &database.NeighborsDatasource{Database: db}}
	id, err := manager.WriteShelter(context.Background(), testShelter)
	if err != nil {
		t.Error(err)
	}
	testShelter.ID = id

	expectedRow := []driver.Value{id}
	expectedRow = append(expectedRow, shelterToRow(testShelter)...)
	mock.ExpectQuery(getSingleShelterQuery).WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"ID", "City", "Country", "Name", "PostalCode", "State", "Street"}).AddRow(expectedRow...))
	actualShelter, err := manager.GetShelter(context.Background(), id)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(testShelter, actualShelter) {
		t.Errorf("Expected %v to equal %v", actualShelter, testShelter)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestItCanDeleteShelter(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	testShelter := generateShelter()

	mock.ExpectExec(createShelterQuery).WithArgs(shelterToRow(testShelter)...).WillReturnResult(sqlmock.NewResult(1, 1))
	manager := &ShelterManager{Datasource: &database.NeighborsDatasource{Database: db}}
	id, err := manager.WriteShelter(context.Background(), testShelter)
	if err != nil {
		t.Error(err)
	}

	mock.ExpectExec(deleteShelterQuery).WithArgs(strconv.FormatInt(id, 10)).WillReturnResult(sqlmock.NewResult(-1, 1))
	rowsDeleted, err := manager.DeleteShelter(context.Background(), strconv.FormatInt(id, 10))
	if err != nil {
		t.Error(err)
	}

	if rowsDeleted != 1 {
		t.Error("Expected row to be deleted")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestItCanGetAllShelters(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	testShelters := make([]*Shelter, 0)
	expectedRows := make([][]driver.Value, 0)
	manager := &ShelterManager{Datasource: &database.NeighborsDatasource{Database: db}}
	for i := 0; i < 5; i++ {
		testShelter := generateShelter()

		mock.ExpectExec(createShelterQuery).WithArgs(shelterToRow(testShelter)...).WillReturnResult(sqlmock.NewResult(1, 1))
		id, err := manager.WriteShelter(context.Background(), testShelter)
		if err != nil {
			t.Error(err)
		}
		testShelter.ID = id
		testShelters = append(testShelters, testShelter)
		expectedRow := []driver.Value{id}
		expectedRow = append(expectedRow, shelterToRow(testShelter)...)
		expectedRows = append(expectedRows, expectedRow)
	}

	rowResult := sqlmock.NewRows([]string{"ID", "City", "Country", "Name", "PostalCode", "State", "Street"})
	for _, expectedRow := range expectedRows {
		rowResult = rowResult.AddRow(expectedRow...)
	}
	mock.ExpectQuery(getAllSheltersQuery).WillReturnRows(rowResult)

	allShelters, err := manager.GetShelters(context.Background())
	if err != nil {
		t.Error(err)
	}

	for _, shelter := range allShelters {
		if !contains(shelter, testShelters) {
			t.Errorf("Expected %v to be in %v \n", shelter, testShelters)
		}
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func shelterToRow(shelter *Shelter) []driver.Value {
	return []driver.Value{shelter.City, shelter.Country, shelter.Name, shelter.PostalCode, shelter.State, shelter.Street}
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

func contains(candidateShelter *Shelter, expectedShelters []*Shelter) bool {
	for _, shelter := range expectedShelters {
		if reflect.DeepEqual(candidateShelter, shelter) {
			return true
		}
	}
	return false
}
