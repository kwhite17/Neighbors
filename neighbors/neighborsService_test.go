package neighbors

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"testing"

	"github.com/kwhite17/Neighbors/test"
)

var service = NeighborServiceHandler{Database: test.TestConnection}

func TestGetAllNeighbors(t *testing.T) {
	test.PopulateNeighborsTable()
	defer test.CleanNeighborsTable()

	req, _ := http.NewRequest("GET", "http://localhost:8080/neighbors/", nil)
	response := test.RecordServiceRequest(service, req)
	data := make([]map[string]interface{}, 0)
	json.NewDecoder(response.Body).Decode(&data)
	if len(data) < 1 {
		t.Errorf("GetAllNeighbors Failure - Expected: %v, Actual: %v\n", 1, len(data))
	}
}

func TestCreateNeighbor(t *testing.T) {
	defer test.CleanNeighborsTable()
	jsonBytes, _ := json.Marshal(buildTestNeighbor())
	req, _ := http.NewRequest("POST", "http://localhost:8080/neighbors/", bytes.NewBuffer(jsonBytes))
	response := test.RecordServiceRequest(service, req)
	if response.StatusCode != http.StatusOK {
		t.Errorf("CreateNeighbor Failure - Expected: %v, Actual: %v\n", http.StatusOK, response.StatusCode)
	}
}

func TestDeleteNeighbor(t *testing.T) {
	defer test.CleanNeighborsTable()
	ids, err := test.PopulateNeighborsTable()
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("DELETE", "http://localhost:8080/neighbors/"+strconv.Itoa(int(ids[0])), nil)
	response := test.RecordServiceRequest(service, req)
	if response.StatusCode != http.StatusOK {
		t.Errorf("DeleteNeighbor Failure - Expected: %v, Actual: %v\n", http.StatusOK, response.StatusCode)
	}
}

func TestUpdateItem(t *testing.T) {
	defer test.CleanNeighborsTable()
	ids, err := test.PopulateNeighborsTable()
	if err != nil {
		t.Fatal(err)
	}
	jsonBytes, err := json.Marshal(map[string]interface{}{"Email": "kevinwhite1710@gmail.com"})
	req, _ := http.NewRequest("PUT", "http://localhost:8080/neighbors/"+string(ids[0]), bytes.NewBuffer(jsonBytes))
	response := test.RecordServiceRequest(service, req)
	if response.StatusCode != http.StatusOK {
		t.Errorf("UpdateNeighbor Failure - Expected: %v, Actual: %v\n", http.StatusOK, response.StatusCode)
	}
}

func TestGetSingleNeighbor(t *testing.T) {
	defer test.CleanNeighborsTable()
	ids, err := test.PopulateNeighborsTable()
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("GET", "http://localhost:8080/neighbors/"+strconv.Itoa(int(ids[1])), nil)
	response := test.RecordServiceRequest(service, req)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("GetSingleNeighbor Failure - Expected: %v, Actual: %v\n", http.StatusOK, response.StatusCode)
	}
	data := make([]map[string]interface{}, 0)
	json.NewDecoder(response.Body).Decode(&data)
	log.Println(data)
	if len(data) != 1 {
		t.Fatalf("GetSingleNeighbor Failure - Expected: %v, Actual: %v\n", 1, len(data))
	}
	if int64(data[0]["NeighborID"].(float64)) != ids[1] {
		t.Errorf("GetSingleNeighbor Failure - Expected: %v, Actual: %v\n", int(ids[1]), data[0]["NeighborID"])
	}
}

func buildTestNeighbor() map[string]interface{} {
	neighbor := make(map[string]interface{})
	neighbor["Username"] = "testItem"
	neighbor["Password"] = "testItem"
	neighbor["Location"] = "Somerville"
	return neighbor
}
