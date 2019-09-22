package managers

import (
	"context"
	"math/rand"
	"reflect"
	"strconv"
	"testing"

	"github.com/kwhite17/Neighbors/pkg/database"
)

var testCategory = "testCategory"
var testGender = "testGender"
var testQuantity = int8(rand.Int() % 127)
var testShelterID = rand.Int63()
var testSize = "testSize"
var testStatus = "testStatus"

func initItemManager() *ItemManager {
	dbToClose = database.InitDatabase(database.SQLITE3)
	return &ItemManager{Datasource: database.StandardDatasource{Database: dbToClose}}
}

func cleanDatabase() {
	dbToClose.Close()
}

func TestCanReadItsOwnItemWrite(t *testing.T) {
	manager := initItemManager()
	defer cleanDatabase()
	testItem := generateItem()

	id, err := manager.WriteItem(context.Background(), testItem)
	if err != nil {
		t.Error(err)
	}
	testItem.ID = id

	actualItem, err := manager.GetItem(context.Background(), id)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(testItem, actualItem) {
		t.Errorf("Expected %v to equal %v", actualItem, testItem)
	}
}

func TestItCanDeleteItem(t *testing.T) {
	manager := initItemManager()
	defer cleanDatabase()
	testItem := generateItem()

	id, err := manager.WriteItem(context.Background(), testItem)
	if err != nil {
		t.Error(err)
	}

	rowsDeleted, err := manager.DeleteItem(context.Background(), strconv.FormatInt(id, 10))
	if err != nil {
		t.Error(err)
	}

	if rowsDeleted != 1 {
		t.Error("Expected row to be deleted")
	}
}

func TestItCanGetAllItems(t *testing.T) {
	manager := initItemManager()
	defer cleanDatabase()
	testItems := make([]*Item, 0)
	for i := 0; i < 5; i++ {
		testItem := generateItem()

		id, err := manager.WriteItem(context.Background(), testItem)
		if err != nil {
			t.Error(err)
		}
		testItem.ID = id
		testItems = append(testItems, testItem)
	}

	allItems, err := manager.GetItems(context.Background())
	if err != nil {
		t.Error(err)
	}

	for _, item := range allItems {
		if !contains(item, testItems) {
			t.Errorf("Expected %v to be in %v \n", item, testItems)
		}
	}
}

func TestCanReadItsOwnItemUpdate(t *testing.T) {
	manager := initItemManager()
	defer cleanDatabase()
	testItem := generateItem()

	id, err := manager.WriteItem(context.Background(), testItem)
	if err != nil {
		t.Error(err)
	}
	testItem.ID = id

	createdItem, err := manager.GetItem(context.Background(), id)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(testItem, createdItem) {
		t.Errorf("Expected %v to equal %v", createdItem, testItem)
	}

	createdItem.Size = "L"
	err = manager.UpdateItem(context.Background(), createdItem)
	if err != nil {
		t.Error(err)
	}

	finalItem, err := manager.GetItem(context.Background(), id)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(finalItem, createdItem) {
		t.Errorf("Expected %v to equal %v", finalItem, createdItem)
	}
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
