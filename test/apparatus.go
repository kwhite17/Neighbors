package test

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/kwhite17/Neighbors/database"
)

var TestConnection = database.NeighborsDatabase
var createTestNeighborsQuery = "INSERT INTO neighbors (Username, Password, Location) VALUES (?, ?, ?)"
var createTestSamaritansQuery = "INSERT INTO samaritans (Username, Password, Location) VALUES (?, ?, ?)"
var createTestItemsQuery = "INSERT INTO items (Category, Size, Quantity, DropoffLocation, Requestor) VALUES (?, ?, ?, ?, ?)"
var deleteTestNeighborsQuery = "DELETE FROM neighbors WHERE Username=?"
var deleteTestSamaritansQuery = "DELETE FROM samaritans WHERE Username=?"
var deleteTestItemsQuery = "DELETE FROM items WHERE Category='testItem'"
var testNeighbors = buildTestNeighbors()
var testItems = buildTestItems()
var testSamaritans = buildTestSamaritans()

func PopulateNeighborsTable() ([]int64, error) {
	ids := make([]int64, 0)
	for _, v := range testNeighbors {
		output := TestConnection.ExecuteWriteQuery(context.Background(), createTestNeighborsQuery, v)
		id, err := output.LastInsertId()
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func PopulateSamaritansTable() ([]int64, error) {
	ids := make([]int64, 0)
	for _, v := range testSamaritans {
		output := TestConnection.ExecuteWriteQuery(context.Background(), createTestSamaritansQuery, v)
		id, err := output.LastInsertId()
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func PopulateItemsTable() error {
	neighborIds, err := PopulateNeighborsTable()
	if err != nil {
		return err
	}
	for i := 0; i < len(neighborIds); i++ {
		TestConnection.ExecuteWriteQuery(context.Background(), createTestItemsQuery, append(testItems["testItem"], neighborIds[i]))
	}
	return nil
}

func RecordServiceRequest(service http.Handler, req *http.Request) *http.Response {
	recorder := httptest.NewRecorder()
	service.ServeHTTP(recorder, req)
	response := recorder.Result()
	return response
}

func buildTestSamaritans() map[string][]interface{} {
	samaritans := make(map[string][]interface{})
	samaritans["testUser"] = []interface{}{"testUser", "testUser", "Somerville"}
	return samaritans
}
func buildTestNeighbors() map[string][]interface{} {
	neighbors := make(map[string][]interface{})
	neighbors["testUser"] = []interface{}{"testUser", "testUser", "Somerville"}
	return neighbors
}

func buildTestItems() map[string][]interface{} {
	items := make(map[string][]interface{})
	items["testItem"] = []interface{}{"testItem", "M", 1, "Home"}
	return items
}

func CleanNeighborsTable() {
	for k := range testNeighbors {
		TestConnection.ExecuteWriteQuery(context.Background(), deleteTestNeighborsQuery, []interface{}{k})
	}
}

func CleanSamaritansTable() {
	for k := range testSamaritans {
		TestConnection.ExecuteWriteQuery(context.Background(), deleteTestSamaritansQuery, []interface{}{k})
	}
}

func CleanItemsTable() {
	TestConnection.ExecuteWriteQuery(context.Background(), deleteTestItemsQuery, []interface{}{})
}
