package neighbors

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	err := json.NewDecoder(response.Body).Decode(&data)
	fmt.Println(err)
	if len(data) < 1 {
		t.Errorf("GetAllNeighbors Failure - Expected: %v, Actual: %v\n", 1, len(data))
	}
}
