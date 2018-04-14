package users

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/kwhite17/Neighbors/pkg/test"
)

var service = UserServiceHandler{Database: test.TestConnection}

func TestRenderNewUserForm(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:8080/users/new", nil)
	response := test.RecordServiceRequest(service, req)
	htmlBytes, _ := ioutil.ReadAll(response.Body)
	htmlString := string(htmlBytes)
	if !strings.Contains(htmlString, "/users/") || !strings.Contains(htmlString, "POST") {
		t.Errorf("RenderNewUserForm Failure - Expected html to contain '/users/' and 'POST', Actual: %s\n", htmlString)
	}
}

func TestGetAllUsers(t *testing.T) {
	test.PopulateUsersTable()
	defer test.CleanUsersTable()

	req, _ := http.NewRequest("GET", "http://localhost:8080/users/", nil)
	response := test.RecordServiceRequest(service, req)
	htmlBytes, _ := ioutil.ReadAll(response.Body)
	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "table") || !strings.Contains(htmlStr, "testUser") {
		t.Errorf("GetAllUsers Failure - Expected html to contain 'table' or 'testUser'")
	}
}

func TestCreateUser(t *testing.T) {
	defer test.CleanUsersTable()
	jsonBytes, _ := json.Marshal(buildTestUser())
	req, _ := http.NewRequest("POST", "http://localhost:8080/users/", bytes.NewBuffer(jsonBytes))
	response := test.RecordServiceRequest(service, req)
	if response.StatusCode != http.StatusOK {
		t.Errorf("CreateUser Failure - Expected: %v, Actual: %v\n", http.StatusOK, response.StatusCode)
	}
	htmlBytes, _ := ioutil.ReadAll(response.Body)
	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "strong") || !strings.Contains(htmlStr, "Tokyo") {
		t.Errorf("CreateNeighbor Failure - Expected html to contain 'strong'")
	}
}

func TestDeleteUser(t *testing.T) {
	defer test.CleanUsersTable()
	ids, err := test.PopulateUsersTable()
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("DELETE", "http://localhost:8080/users/"+strconv.Itoa(int(ids[0])), nil)
	response := test.RecordServiceRequest(service, req)
	if response.StatusCode != http.StatusOK {
		t.Errorf("DeleteUser Failure - Expected: %v, Actual: %v\n", http.StatusOK, response.StatusCode)
	}
}

func TestUpdateItem(t *testing.T) {
	defer test.CleanUsersTable()
	ids, err := test.PopulateUsersTable()
	if err != nil {
		t.Fatal(err)
	}
	jsonBytes, err := json.Marshal(map[string]interface{}{"Location": "Tokyo", "UserID": strconv.Itoa(int(ids[0]))})
	req, _ := http.NewRequest("PUT", "http://localhost:8080/users/"+string(ids[0]), bytes.NewBuffer(jsonBytes))
	response := test.RecordServiceRequest(service, req)
	if response.StatusCode != http.StatusOK {
		t.Errorf("UpdateUser Failure - Expected: %v, Actual: %v\n", http.StatusOK, response.StatusCode)
	}
	htmlBytes, _ := ioutil.ReadAll(response.Body)
	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "strong") || !strings.Contains(htmlStr, "Tokyo") {
		t.Errorf("UpdateUser Failure - Expected: html to contain 'strong', Actual: %v\n", htmlStr)
	}
}

func TestGetSingleUser(t *testing.T) {
	defer test.CleanUsersTable()
	ids, err := test.PopulateUsersTable()
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("GET", "http://localhost:8080/users/"+strconv.Itoa(int(ids[1])), nil)
	response := test.RecordServiceRequest(service, req)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("GetSingleUser Failure - Expected: %v, Actual: %v\n", http.StatusOK, response.StatusCode)
	}
	htmlBytes, _ := ioutil.ReadAll(response.Body)
	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "strong") || !strings.Contains(htmlStr, strconv.Itoa(int(ids[1]))) {
		t.Errorf("GetSingleUser Failure - Expected html to contain 'strong' or correct ID")
	}
}

func buildTestUser() map[string]interface{} {
	user := make(map[string]interface{})
	user["Username"] = "testUser"
	user["Password"] = "testUser"
	user["Location"] = "Tokyo"
	return user
}
