package items

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kwhite17/Neighbors/test"
)

var service = ItemServiceHandler{Database: test.TestConnection}

func TestGetAllItems(t *testing.T) {
	defer test.CleanNeighborsTable()
	defer test.CleanItemsTable()
	err := test.PopulateItemsTable()
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("GET", "http://localhost:8080/items/", nil)
	response := test.RecordServiceRequest(service, req)
	data := make([]map[string]interface{}, 0)
	json.NewDecoder(response.Body).Decode(&data)
	if len(data) < 1 {
		t.Errorf("GetAllItems Failure - Expected: %v, Actual: %v\n", 1, len(data))
	}
}
