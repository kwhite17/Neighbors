package items

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

var service = ItemResourceHandler{Database: test.TestConnection}

func createServer() *httptest.Server {
	directory, _ := os.Getwd()
	testMux := http.NewServeMux()
	testMux.Handle("/items/", service)
	testMux.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir(directory+"/templates/"))))
	return httptest.NewServer(testMux)
}

func teardown(ts *httptest.Server) {
	ts.CloseClientConnections()
	ts.Close()
	test.CleanUserSessionTable()
	test.CleanUsersTable()
	test.CleanItemsTable()
}

func TestRenderNewItemForm(t *testing.T) {
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

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/items/new", nil)
	req.AddCookie(&cookie)
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("NewItem Failure - Expected: %d, Actual: %d\n", http.StatusOK, response.StatusCode)
		t.FailNow()
	}
	htmlBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	htmlString := string(htmlBytes)
	if !strings.Contains(htmlString, "/items/") || !strings.Contains(htmlString, "POST") {
		t.Errorf("RenderNewItemsForm Failure - Expected html to contain '/items/' and 'POST', Actual: %s\n", htmlString)
	}
}
func TestGetAllItems(t *testing.T) {
	ts := createServer()
	client := ts.Client()
	defer teardown(ts)
	_, neighborIds, err := test.PopulateItemsTable()
	if err != nil {
		t.Fatal(err)
	}
	cookieID, err := test.BuildUserSession(service, neighborIds[0])
	if err != nil {
		t.Fatal(err)
	}
	cookie := http.Cookie{Name: "NeighborsAuth", Value: cookieID, HttpOnly: true, Secure: true, MaxAge: 24 * 60 * 60}

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/items/", nil)
	req.AddCookie(&cookie)
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("GetAllItems Failure - Expected: %d, Actual: %d\n", http.StatusOK, response.StatusCode)
		t.FailNow()
	}
	htmlBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "table") || !strings.Contains(htmlStr, "TESTITEM") {
		t.Errorf("GetAllItems Failure - Expected html to contain 'table' or 'TESTITEM'")
	}
}

func TestCreateItem(t *testing.T) {
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

	jsonBytes, _ := json.Marshal(buildTestItem(ids[0]))
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/items/", bytes.NewBuffer(jsonBytes))
	req.AddCookie(&cookie)
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("CreateItem Failure - Expected: %d, Actual: %d\n", http.StatusOK, response.StatusCode)
		t.FailNow()
	}
	htmlBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "strong") || !strings.Contains(htmlStr, "Shelter") {
		t.Errorf("CreateItem Failure - Expected: html to contain 'strong', Actual: %s\n", htmlStr)
	}
}

func TestDeleteItem(t *testing.T) {
	ts := createServer()
	client := ts.Client()
	defer teardown(ts)
	ids, neighborIds, err := test.PopulateItemsTable()
	if err != nil {
		t.Fatal(err)
	}
	cookieID, err := test.BuildUserSession(service, neighborIds[0])
	if err != nil {
		t.Fatal(err)
	}
	cookie := http.Cookie{Name: "NeighborsAuth", Value: cookieID, HttpOnly: true, Secure: true, MaxAge: 24 * 60 * 60}

	req, _ := http.NewRequest("DELETE", ts.URL+"/items/"+strconv.FormatInt(ids[0], 10), nil)
	req.AddCookie(&cookie)
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("DeleteItem Failure - Expected: %d, Actual: %d\n", http.StatusOK, response.StatusCode)
	}
}

func TestUpdateItem(t *testing.T) {
	ts := createServer()
	client := ts.Client()
	defer teardown(ts)
	itemIds, neighborIds, err := test.PopulateItemsTable()
	if err != nil {
		t.Fatal(err)
	}
	cookieID, err := test.BuildUserSession(service, neighborIds[0])
	if err != nil {
		t.Fatal(err)
	}
	cookie := http.Cookie{Name: "NeighborsAuth", Value: cookieID, HttpOnly: true, Secure: true, MaxAge: 24 * 60 * 60}

	jsonBytes, err := json.Marshal(map[string]interface{}{"DropoffLocation": "Shelter", "Fulfiller": neighborIds[0], "ID": strconv.FormatInt(itemIds[0], 10)})
	req, _ := http.NewRequest("PUT", ts.URL+"/items/"+strconv.FormatInt(itemIds[0], 10), bytes.NewBuffer(jsonBytes))
	req.AddCookie(&cookie)
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("UpdateItem Failure - Expected: %d, Actual: %d\n", http.StatusOK, response.StatusCode)
	}
	htmlBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "strong") || !strings.Contains(htmlStr, "Shelter") {
		t.Errorf("UpdateItem Failure - Expected html to contain 'strong' or 'testItem'")
	}
}

func TestGetSingleItem(t *testing.T) {
	ts := createServer()
	client := ts.Client()
	defer teardown(ts)
	itemIds, neighborIds, err := test.PopulateItemsTable()
	if err != nil {
		t.Fatal(err)
	}
	cookieID, err := test.BuildUserSession(service, neighborIds[0])
	if err != nil {
		t.Fatal(err)
	}
	cookie := http.Cookie{Name: "NeighborsAuth", Value: cookieID, HttpOnly: true, Secure: true, MaxAge: 24 * 60 * 60}

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/items/"+strconv.FormatInt(itemIds[1], 10), nil)
	req.AddCookie(&cookie)
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("GetSingleItem Failure - Expected: %d, Actual: %d\n", http.StatusOK, response.StatusCode)
		t.FailNow()
	}
	htmlBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "Item Request") || !strings.Contains(htmlStr, strconv.Itoa(int(itemIds[1]))) {
		t.Errorf("GetSingleItem Failure - Expected html to contain 'strong' or correct ID, Actual: %s\n", htmlStr)
	}
}

func buildTestItem(requestorId int64) map[string]interface{} {
	item := make(map[string]interface{})
	item["Category"] = "TESTITEM"
	item["Size"] = "M"
	item["Quantity"] = 1
	item["Requestor"] = requestorId
	item["DropoffLocation"] = "Shelter"
	return item
}
