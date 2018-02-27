package neighbors

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/kwhite17/Neighbors/test"
)

var service = NeighborServiceHandler{Database: test.TestConnection}

func TestGetAllNeighbors(t *testing.T) {
	test.PopulateNeighborsTable()
	defer test.CleanNeighborsTable()

	req, _ := http.NewRequest("GET", "http://localhost:8080/neighbors/", nil)
	response := test.RecordServiceRequest(service, req)
	htmlBytes, _ := ioutil.ReadAll(response.Body)
	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "table") || !strings.Contains(htmlStr, "testUser") {
		t.Errorf("GetAllNeighbors Failure - Expected html to contain 'table' or 'testUser'")
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
	htmlBytes, _ := ioutil.ReadAll(response.Body)
	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "strong") || !strings.Contains(htmlStr, "Tokyo") {
		t.Errorf("CreateNeighbor Failure - Expected html to contain 'strong'")
	}
}

func TestDeleteNeighbor(t *testing.T) {
	defer test.CleanNeighborsTable()
	ids, err := test.PopulateNeighborsTable()
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("DELETE", "http://localhost:8080/neighbors/"+strconv.FormatInt(ids[0], 10), nil)
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
	jsonBytes, err := json.Marshal(map[string]interface{}{"Location": "Tokyo", "NeighborID": strconv.Itoa(int(ids[0]))})
	req, _ := http.NewRequest("PUT", "http://localhost:8080/neighbors/", bytes.NewBuffer(jsonBytes))
	response := test.RecordServiceRequest(service, req)
	if response.StatusCode != http.StatusOK {
		t.Errorf("UpdateNeighbor Failure - Expected: %v, Actual: %v\n", http.StatusOK, response.StatusCode)
	}
	htmlBytes, _ := ioutil.ReadAll(response.Body)
	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "strong") || !strings.Contains(htmlStr, "Tokyo") {
		t.Errorf("UpdateNeighbor Failure - Expected: html to contain 'strong', Actual: %v\n", htmlStr)
	}
}

func TestGetSingleNeighbor(t *testing.T) {
	defer test.CleanNeighborsTable()
	ids, err := test.PopulateNeighborsTable()
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("GET", "http://localhost:8080/neighbors/"+strconv.FormatInt(ids[1], 10), nil)
	response := test.RecordServiceRequest(service, req)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("GetSingleNeighbor Failure - Expected: %v, Actual: %v\n", http.StatusOK, response.StatusCode)
	}
	htmlBytes, _ := ioutil.ReadAll(response.Body)
	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "strong") || !strings.Contains(htmlStr, strconv.FormatInt(ids[1], 10)) {
		t.Errorf("GetAllNeighbors Failure - Expected html to contain 'strong' or correct ID")
	}
}

func buildTestNeighbor() map[string]interface{} {
	neighbor := make(map[string]interface{})
	neighbor["Username"] = "testItem"
	neighbor["Password"] = "testItem"
	neighbor["Location"] = "Tokyo"
	return neighbor
}
