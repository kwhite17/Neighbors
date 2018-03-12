package items

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/kwhite17/Neighbors/pkg/test"
)

var service = ItemServiceHandler{Database: test.TestConnection}

func TestRenderNewItemForm(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:8080/items/new", nil)
	response := test.RecordServiceRequest(service, req)
	htmlBytes, _ := ioutil.ReadAll(response.Body)
	htmlString := string(htmlBytes)
	if !strings.Contains(htmlString, "/items/") || !strings.Contains(htmlString, "POST") {
		t.Errorf("RenderNewItemsForm Failure - Expected html to contain '/items/' and 'POST', Actual: %s\n", htmlString)
	}
}
func TestGetAllItems(t *testing.T) {
	defer test.CleanNeighborsTable()
	defer test.CleanItemsTable()
	defer test.CleanSamaritansTable()
	_, err := test.PopulateItemsTable()
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("GET", "http://localhost:8080/items/", nil)
	response := test.RecordServiceRequest(service, req)
	htmlBytes, _ := ioutil.ReadAll(response.Body)
	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "table") || !strings.Contains(htmlStr, "testItem") {
		t.Errorf("GetAllNeighbors Failure - Expected html to contain 'table' or 'testItem'")
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
	htmlBytes, _ := ioutil.ReadAll(response.Body)
	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "strong") || !strings.Contains(htmlStr, "Shelter") {
		t.Errorf("CreateItem Failure - Expected: html to contain 'strong', Actual: %v\n", htmlStr)
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
	jsonBytes, err := json.Marshal(map[string]interface{}{"DropoffLocation": "Shelter", "Fulfiller": samaritanIds[0], "ItemID": strconv.FormatInt(itemIds[0], 10)})
	req, _ := http.NewRequest("PUT", "http://localhost:8080/items/"+strconv.FormatInt(itemIds[0], 10), bytes.NewBuffer(jsonBytes))
	response := test.RecordServiceRequest(service, req)
	if response.StatusCode != http.StatusOK {
		t.Errorf("UpdateItem Failure - Expected: %v, Actual: %v\n", http.StatusOK, response.StatusCode)
	}
	htmlBytes, _ := ioutil.ReadAll(response.Body)
	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "strong") || !strings.Contains(htmlStr, "Shelter") {
		t.Errorf("UpdateItem Failure - Expected html to contain 'strong' or 'testItem'")
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
	htmlBytes, _ := ioutil.ReadAll(response.Body)
	htmlStr := string(htmlBytes)
	fmt.Println(htmlStr)
	if !strings.Contains(htmlStr, "Item Request") || !strings.Contains(htmlStr, strconv.Itoa(int(itemIds[1]))) {
		t.Errorf("GetAllNeighbors Failure - Expected html to contain 'ItemReques' or correct ID")
	}
}

func buildTestItem(requestorId int64) map[string]interface{} {
	item := make(map[string]interface{})
	item["Category"] = "testItem"
	item["Size"] = "M"
	item["Requestor"] = requestorId
	item["DropoffLocation"] = "Shelter"
	return item
}
