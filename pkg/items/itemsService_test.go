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

var service = ItemServiceHandler{Database: test.TestConnection}

func createServer() *httptest.Server {
	directory, _ := os.Getwd()
	testMux := http.NewServeMux()
	testMux.Handle("/items/", service)
	testMux.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir(directory+"/templates/"))))
	return httptest.NewServer(testMux)
}
func TestRenderNewItemForm(t *testing.T) {
	ts := createServer()
	client := ts.Client()
	defer ts.CloseClientConnections()
	defer ts.Close()

	response, err := client.Get(ts.URL + "/items/new")
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
	defer ts.CloseClientConnections()
	defer ts.Close()
	defer test.CleanUsersTable()
	defer test.CleanItemsTable()
	_, _, err := test.PopulateItemsTable()
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.Get(ts.URL + "/items/")
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
	defer ts.CloseClientConnections()
	defer ts.Close()
	defer test.CleanUsersTable()
	defer test.CleanItemsTable()
	ids, err := test.PopulateUsersTable()
	if err != nil {
		t.Fatal(err)
	}

	jsonBytes, _ := json.Marshal(buildTestItem(ids[0]))
	response, err := client.Post(ts.URL+"/items/", "application/json", bytes.NewBuffer(jsonBytes))
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
	defer ts.CloseClientConnections()
	defer ts.Close()
	defer test.CleanUsersTable()
	defer test.CleanItemsTable()
	ids, _, err := test.PopulateItemsTable()
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest("DELETE", ts.URL+"/items/"+strconv.FormatInt(ids[0], 10), nil)
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
	defer ts.CloseClientConnections()
	defer ts.Close()
	defer test.CleanUsersTable()
	defer test.CleanItemsTable()
	itemIds, neighborIds, err := test.PopulateItemsTable()
	if err != nil {
		t.Fatal(err)
	}
	jsonBytes, err := json.Marshal(map[string]interface{}{"DropoffLocation": "Shelter", "Fulfiller": neighborIds[0], "ID": strconv.FormatInt(itemIds[0], 10)})
	req, _ := http.NewRequest("PUT", ts.URL+"/items/"+strconv.FormatInt(itemIds[0], 10), bytes.NewBuffer(jsonBytes))
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
	defer ts.CloseClientConnections()
	defer ts.Close()
	defer test.CleanUsersTable()
	defer test.CleanItemsTable()
	itemIds, _, err := test.PopulateItemsTable()
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.Get(ts.URL + "/items/" + strconv.FormatInt(itemIds[1], 10))
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
