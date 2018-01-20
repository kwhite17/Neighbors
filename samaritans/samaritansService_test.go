package samaritans

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/kwhite17/Neighbors/test"
)

var service = SamaritanServiceHandler{Database: test.TestConnection}

func TestGetAllNeighbors(t *testing.T) {
	test.PopulateSamaritansTable()
	defer test.CleanSamaritansTable()

	req, _ := http.NewRequest("GET", "http://localhost:8080/samaritans/", nil)
	response := test.RecordServiceRequest(service, req)
	data := make([]map[string]interface{}, 0)
	err := json.NewDecoder(response.Body).Decode(&data)
	fmt.Println(err)
	if len(data) < 1 {
		t.Errorf("GetAllNeighbors Failure - Expected: %v, Actual: %v\n", 1, len(data))
	}
}
