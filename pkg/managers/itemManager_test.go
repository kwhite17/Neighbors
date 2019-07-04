package managers

import (
	"context"
	"database/sql/driver"
	"math/rand"
	"reflect"
	"strconv"
	"testing"

	"github.com/kwhite17/Neighbors/pkg/database"

	"github.com/DATA-DOG/go-sqlmock"
)

var testCategory = "testCategory"
var testGender = "testGender"
var testQuantity = int8(rand.Int() % 127)
var testShelterID = rand.Int63()
var testSize = "testSize"
var testStatus = "testStatus"

func TestCanReadItsOwnItemWrite(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	testItem := generateItem()

	mock.ExpectExec(createItemQuery).WithArgs(itemToRow(testItem)...).WillReturnResult(sqlmock.NewResult(1, 1))
	manager := &ItemManager{Datasource: &database.NeighborsDatasource{Database: db}}
	id, err := manager.WriteItem(context.Background(), testItem)
	if err != nil {
		t.Error(err)
	}
	testItem.ID = id

	expectedRow := []driver.Value{id}
	expectedRow = append(expectedRow, itemToRow(testItem)...)
	mock.ExpectQuery(getSingleItemQuery).WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"ID", "Category", "Gender", "Quantity", "ShelterID", "Size", "Status"}).AddRow(expectedRow...))
	actualItem, err := manager.GetItem(context.Background(), id)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(testItem, actualItem) {
		t.Errorf("Expected %v to equal %v", actualItem, testItem)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestItCanDeleteItem(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	testItem := generateItem()

	mock.ExpectExec(createItemQuery).WithArgs(itemToRow(testItem)...).WillReturnResult(sqlmock.NewResult(1, 1))
	manager := &ItemManager{Datasource: &database.NeighborsDatasource{Database: db}}
	id, err := manager.WriteItem(context.Background(), testItem)
	if err != nil {
		t.Error(err)
	}

	mock.ExpectExec(deleteItemQuery).WithArgs(strconv.FormatInt(id, 10)).WillReturnResult(sqlmock.NewResult(-1, 1))
	rowsDeleted, err := manager.DeleteItem(context.Background(), strconv.FormatInt(id, 10))
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

func TestItCanGetAllItems(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	testItems := make([]*Item, 0)
	expectedRows := make([][]driver.Value, 0)
	manager := &ItemManager{Datasource: &database.NeighborsDatasource{Database: db}}
	for i := 0; i < 5; i++ {
		testItem := generateItem()

		mock.ExpectExec(createItemQuery).WithArgs(itemToRow(testItem)...).WillReturnResult(sqlmock.NewResult(1, 1))
		id, err := manager.WriteItem(context.Background(), testItem)
		if err != nil {
			t.Error(err)
		}
		testItem.ID = id
		testItems = append(testItems, testItem)
		expectedRow := []driver.Value{id}
		expectedRow = append(expectedRow, itemToRow(testItem)...)
		expectedRows = append(expectedRows, expectedRow)
	}

	rowResult := sqlmock.NewRows([]string{"ID", "Category", "Gender", "Quantity", "ShelterID", "Size", "Status"})
	for _, expectedRow := range expectedRows {
		rowResult = rowResult.AddRow(expectedRow...)
	}
	mock.ExpectQuery(getAllItemsQuery).WillReturnRows(rowResult)

	allItems, err := manager.GetItems(context.Background())
	if err != nil {
		t.Error(err)
	}

	for _, item := range allItems {
		if !contains(item, testItems) {
			t.Errorf("Expected %v to be in %v \n", item, testItems)
		}
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func itemToRow(item *Item) []driver.Value {
	return []driver.Value{item.Category, item.Gender, item.Quantity, item.ShelterID, item.Size, item.Status}
}

func generateItem() *Item {
	return &Item{
		Category:  testCategory,
		Gender:    testGender,
		Quantity:  testQuantity,
		ShelterID: testShelterID,
		Size:      testSize,
		Status:    testStatus,
	}
}

func contains(candidateItem *Item, expectedItems []*Item) bool {
	for _, item := range expectedItems {
		if reflect.DeepEqual(candidateItem, item) {
			return true
		}
	}
	return false
}
