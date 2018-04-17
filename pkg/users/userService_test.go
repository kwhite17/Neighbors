package users

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/kwhite17/Neighbors/pkg/test"
)

var service = UserServiceHandler{Database: test.TestConnection}

func createServer() *httptest.Server {
	testMux := http.NewServeMux()
	testMux.Handle("/users/", service)
	return httptest.NewServer(testMux)
}

func TestRenderNewUserForm(t *testing.T) {
	ts := createServer()
	client := ts.Client()
	defer ts.CloseClientConnections()
	defer ts.Close()

	response, err := client.Get(ts.URL + "/users/new")
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("NewUser Failure - Expected: %d, Actual: %d\n", http.StatusOK, response.StatusCode)
		t.FailNow()
	}
	htmlBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	htmlString := string(htmlBytes)
	if !strings.Contains(htmlString, "/users/") || !strings.Contains(htmlString, "POST") {
		t.Errorf("RenderNewUserForm Failure - Expected html to contain '/users/' and 'POST', Actual: %s\n", htmlString)
	}
}

func TestGetAllUsers(t *testing.T) {
	test.PopulateUsersTable()
	ts := createServer()
	client := ts.Client()
	defer ts.CloseClientConnections()
	defer ts.Close()
	defer test.CleanUsersTable()

	response, err := client.Get(ts.URL + "/users/")
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("GetAllUsers Failure - Expected: %d, Actual: %d\n", http.StatusOK, response.StatusCode)
		t.FailNow()
	}
	htmlBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "table") || !strings.Contains(htmlStr, "testUser") {
		t.Errorf("GetAllUsers Failure - Expected html to contain 'table' or 'testUser'")
	}
}

func TestCreateUser(t *testing.T) {
	ts := createServer()
	client := ts.Client()
	defer ts.CloseClientConnections()
	defer ts.Close()
	defer test.CleanUsersTable()

	jsonBytes, _ := json.Marshal(buildTestUser())
	response, err := client.Post(ts.URL+"/users/", "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("CreateUser Failure - Expected: %d, Actual: %d\n", http.StatusOK, response.StatusCode)
		t.FailNow()
	}
	htmlBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "strong") || !strings.Contains(htmlStr, "Tokyo") {
		t.Errorf("CreateNeighbor Failure - Expected html to contain 'strong' and 'Tokyo', Actual: %s\n", htmlStr)
	}
}

func TestDeleteUser(t *testing.T) {
	ts := createServer()
	client := ts.Client()
	defer ts.CloseClientConnections()
	defer ts.Close()
	defer test.CleanUsersTable()
	ids, err := test.PopulateUsersTable()
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest("DELETE", ts.URL+"/users/"+strconv.FormatInt(ids[0], 10), nil)
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("DeleteUser Failure - Expected: %d, Actual: %d\n", http.StatusOK, response.StatusCode)
		t.FailNow()
	}
}

func TestUpdateItem(t *testing.T) {
	ts := createServer()
	client := ts.Client()
	defer ts.CloseClientConnections()
	defer ts.Close()
	defer test.CleanUsersTable()
	ids, err := test.PopulateUsersTable()
	if err != nil {
		t.Fatal(err)
	}

	jsonBytes, err := json.Marshal(map[string]interface{}{"Location": "Tokyo", "ID": strconv.Itoa(int(ids[0]))})
	req, _ := http.NewRequest("PUT", ts.URL+"/users/"+strconv.FormatInt(ids[0], 10), bytes.NewBuffer(jsonBytes))
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("UpdateUser Failure - Expected: %d, Actual: %d\n", http.StatusOK, response.StatusCode)
		t.FailNow()
	}
	htmlBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "strong") || !strings.Contains(htmlStr, "Tokyo") {
		t.Errorf("UpdateUser Failure - Expected: html to contain 'strong', Actual: %s\n", htmlStr)
	}
}

func TestGetSingleUser(t *testing.T) {
	ts := createServer()
	client := ts.Client()
	defer ts.CloseClientConnections()
	defer ts.Close()
	defer test.CleanUsersTable()
	ids, err := test.PopulateUsersTable()
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.Get(ts.URL + "/users/" + strconv.FormatInt(ids[1], 10))
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("GetSingleUser Failure - Expected: %d, Actual: %d\n", http.StatusOK, response.StatusCode)
		t.FailNow()
	}
	htmlBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "strong") || !strings.Contains(htmlStr, strconv.FormatInt(ids[1], 10)) {
		t.Errorf("GetSingleUser Failure - Expected html to contain 'strong' or correct ID, Actual: %s\n", htmlStr)
	}
}

func buildTestUser() map[string]interface{} {
	user := make(map[string]interface{})
	user["Username"] = "testUser"
	user["Password"] = "testUser"
	user["Location"] = "Tokyo"
	user["Role"] = "SAMARITAN"
	return user
}
