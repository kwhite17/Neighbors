package samaritans

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

var service = SamaritanServiceHandler{Database: test.TestConnection}

func TestGetAllSamaritans(t *testing.T) {
	test.PopulateSamaritansTable()
	defer test.CleanSamaritansTable()

	req, _ := http.NewRequest("GET", "http://localhost:8080/samaritans/", nil)
	response := test.RecordServiceRequest(service, req)
	htmlBytes, _ := ioutil.ReadAll(response.Body)
	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "table") || !strings.Contains(htmlStr, "testUser") {
		t.Errorf("GetAllSamaritans Failure - Expected html to contain 'table' or 'testUser'")
	}
}

func TestCreateSamaritan(t *testing.T) {
	defer test.CleanSamaritansTable()
	jsonBytes, _ := json.Marshal(buildTestSamaritan())
	req, _ := http.NewRequest("POST", "http://localhost:8080/samaritans/", bytes.NewBuffer(jsonBytes))
	response := test.RecordServiceRequest(service, req)
	if response.StatusCode != http.StatusOK {
		t.Errorf("CreateSamaritan Failure - Expected: %v, Actual: %v\n", http.StatusOK, response.StatusCode)
	}
}

func TestDeleteSamaritan(t *testing.T) {
	defer test.CleanSamaritansTable()
	ids, err := test.PopulateSamaritansTable()
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("DELETE", "http://localhost:8080/samaritans/"+strconv.Itoa(int(ids[0])), nil)
	response := test.RecordServiceRequest(service, req)
	if response.StatusCode != http.StatusOK {
		t.Errorf("DeleteSamaritan Failure - Expected: %v, Actual: %v\n", http.StatusOK, response.StatusCode)
	}
}

func TestUpdateItem(t *testing.T) {
	defer test.CleanSamaritansTable()
	ids, err := test.PopulateSamaritansTable()
	if err != nil {
		t.Fatal(err)
	}
	jsonBytes, err := json.Marshal(map[string]interface{}{"Email": "kevinwhite1710@gmail.com"})
	req, _ := http.NewRequest("PUT", "http://localhost:8080/samaritans/"+string(ids[0]), bytes.NewBuffer(jsonBytes))
	response := test.RecordServiceRequest(service, req)
	if response.StatusCode != http.StatusOK {
		t.Errorf("UpdateSamaritan Failure - Expected: %v, Actual: %v\n", http.StatusOK, response.StatusCode)
	}
}

func TestGetSingleSamaritan(t *testing.T) {
	defer test.CleanSamaritansTable()
	ids, err := test.PopulateSamaritansTable()
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("GET", "http://localhost:8080/samaritans/"+strconv.Itoa(int(ids[1])), nil)
	response := test.RecordServiceRequest(service, req)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("GetSingleSamaritan Failure - Expected: %v, Actual: %v\n", http.StatusOK, response.StatusCode)
	}
	htmlBytes, _ := ioutil.ReadAll(response.Body)
	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "strong") || !strings.Contains(htmlStr, strconv.Itoa(int(ids[1]))) {
		t.Errorf("GetSingleSamaritan Failure - Expected html to contain 'strong' or correct ID")
	}
}

func buildTestSamaritan() map[string]interface{} {
	samaritan := make(map[string]interface{})
	samaritan["Username"] = "testItem"
	samaritan["Password"] = "testItem"
	samaritan["Location"] = "Somerville"
	return samaritan
}
