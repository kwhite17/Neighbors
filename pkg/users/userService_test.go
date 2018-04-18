package users

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/kwhite17/Neighbors/pkg/test"
)

var service = UserServiceHandler{Database: test.TestConnection}

func createServer() *httptest.Server {
	directory, _ := os.Getwd()
	testMux := http.NewServeMux()
	testMux.Handle("/users/", service)
	testMux.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir(directory+"/templates/"))))
	return httptest.NewServer(testMux)
}

func teardown(ts *httptest.Server) {
	ts.CloseClientConnections()
	ts.Close()
	test.CleanUserSessionTable()
	test.CleanUsersTable()
}

func TestGetAllUsers(t *testing.T) {
	ts := createServer()
	client := ts.Client()
	defer teardown(ts)
	ids, err := test.PopulateUsersTable()
	if err != nil {
		t.Fatal(err)
	}
	cookieID, err := test.BuildUserSession(service, ids[0])
	if err != nil {
		t.Fatal(err)
	}
	cookie := http.Cookie{Name: "NeighborsAuth", Value: cookieID, HttpOnly: true, Secure: true, MaxAge: 24 * 60 * 60}

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/users/", nil)
	req.AddCookie(&cookie)
	response, err := client.Do(req)
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
	defer teardown(ts)

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
	defer teardown(ts)
	ids, err := test.PopulateUsersTable()
	if err != nil {
		t.Fatal(err)
	}
	cookieID, err := test.BuildUserSession(service, ids[0])
	if err != nil {
		t.Fatal(err)
	}
	cookie := http.Cookie{Name: "NeighborsAuth", Value: cookieID, HttpOnly: true, Secure: true, MaxAge: 24 * 60 * 60}

	req, _ := http.NewRequest("DELETE", ts.URL+"/users/"+strconv.FormatInt(ids[0], 10), nil)
	req.AddCookie(&cookie)
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
	defer teardown(ts)
	ids, err := test.PopulateUsersTable()
	if err != nil {
		t.Fatal(err)
	}
	cookieID, err := test.BuildUserSession(service, ids[0])
	if err != nil {
		t.Fatal(err)
	}
	cookie := http.Cookie{Name: "NeighborsAuth", Value: cookieID, HttpOnly: true, Secure: true, MaxAge: 24 * 60 * 60}

	jsonBytes, err := json.Marshal(map[string]interface{}{"Location": "Tokyo", "ID": strconv.Itoa(int(ids[0]))})
	req, _ := http.NewRequest("PUT", ts.URL+"/users/"+strconv.FormatInt(ids[0], 10), bytes.NewBuffer(jsonBytes))
	req.AddCookie(&cookie)
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
	cookieID, err := test.BuildUserSession(service, ids[0])
	if err != nil {
		t.Fatal(err)
	}
	cookie := http.Cookie{Name: "NeighborsAuth", Value: cookieID, HttpOnly: true, Secure: true, MaxAge: 24 * 60 * 60}

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/users/"+strconv.FormatInt(ids[1], 10), nil)
	req.AddCookie(&cookie)
	response, err := client.Do(req)
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
