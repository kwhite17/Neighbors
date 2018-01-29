package items

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"testing"

	"github.com/kwhite17/Neighbors/test"
)

var service = ItemServiceHandler{Database: test.TestConnection}

func TestGetAllItems(t *testing.T) {
	defer test.CleanNeighborsTable()
	defer test.CleanItemsTable()
	defer test.CleanSamaritansTable()
	itemIds, err := test.PopulateItemsTable()
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("GET", "http://localhost:8080/items/", nil)
	response := test.RecordServiceRequest(service, req)
	data := make([]map[string]interface{}, 0)
	json.NewDecoder(response.Body).Decode(&data)
	if len(data) != len(itemIds) {
		t.Errorf("GetAllItems Failure - Expected: %v, Actual: %v\n", len(itemIds), len(data))
	}
}

func TestCreateItem(t *testing.T) {
	defer test.CleanNeighborsTable()
	defer test.CleanItemsTable()
	defer test.CleanSamaritansTable()
	ids, err := test.PopulateNeighborsTable()
	if err != nil {
		t.Fatal(err)
	}
	jsonBytes, _ := json.Marshal(buildTestItem(ids[0]))
	req, _ := http.NewRequest("POST", "http://localhost:8080/items/", bytes.NewBuffer(jsonBytes))
	response := test.RecordServiceRequest(service, req)
	if response.StatusCode != http.StatusOK {
		t.Errorf("CreateItem Failure - Expected: %v, Actual: %v\n", http.StatusOK, response.StatusCode)
	}
}

func TestDeleteItem(t *testing.T) {
	defer test.CleanNeighborsTable()
	defer test.CleanItemsTable()
	defer test.CleanSamaritansTable()
	ids, err := test.PopulateItemsTable()
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("DELETE", "http://localhost:8080/items/"+strconv.Itoa(int(ids[0])), nil)
	response := test.RecordServiceRequest(service, req)
	if response.StatusCode != http.StatusOK {
		t.Errorf("DeleteItem Failure - Expected: %v, Actual: %v\n", http.StatusOK, response.StatusCode)
	}
}

func TestUpdateItem(t *testing.T) {
	defer test.CleanNeighborsTable()
	defer test.CleanSamaritansTable()
	defer test.CleanItemsTable()
	itemIds, err := test.PopulateItemsTable()
	if err != nil {
		t.Fatal(err)
	}
	samaritanIds, err := test.PopulateSamaritansTable()
	if err != nil {
		t.Fatal(err)
	}
	jsonBytes, err := json.Marshal(map[string]interface{}{"Fulfiller": samaritanIds[0]})
	req, _ := http.NewRequest("PUT", "http://localhost:8080/items/"+string(itemIds[0]), bytes.NewBuffer(jsonBytes))
	response := test.RecordServiceRequest(service, req)
	if response.StatusCode != http.StatusOK {
		t.Errorf("UpdateItem Failure - Expected: %v, Actual: %v\n", http.StatusOK, response.StatusCode)
	}
}

func TestGetSingleItem(t *testing.T) {
	defer test.CleanNeighborsTable()
	defer test.CleanSamaritansTable()
	defer test.CleanItemsTable()
	itemIds, err := test.PopulateItemsTable()
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("GET", "http://localhost:8080/items/"+strconv.Itoa(int(itemIds[1])), nil)
	response := test.RecordServiceRequest(service, req)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("GetSingleItem Failure - Expected: %v, Actual: %v\n", http.StatusOK, response.StatusCode)
	}
	data := make([]map[string]interface{}, 0)
	json.NewDecoder(response.Body).Decode(&data)
	log.Println(data)
	if len(data) != 1 {
		t.Fatalf("GetSingleItem Failure - Expected: %v, Actual: %v\n", 1, len(data))
	}
	if int64(data[0]["ItemId"].(float64)) != itemIds[1] {
		t.Errorf("GetSingleItem Failure - Expected: %v, Actual: %v\n", int(itemIds[1]), data[0]["ItemId"])
	}
}

func buildTestItem(requestorId int64) map[string]interface{} {
	item := make(map[string]interface{})
	item["Category"] = "testItem"
	item["Size"] = "M"
	item["Requestor"] = requestorId
	item["DropoffLocation"] = "Home"
	return item
}
