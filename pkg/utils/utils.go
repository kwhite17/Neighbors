package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/kwhite17/Neighbors/pkg/database"
)

type ServiceHandler interface {
	BuildGenericResponse(result *sql.Rows) ([]map[string]interface{}, error)
	GetDatasource() database.Datasource
}

func BuildJsonResponse(result *sql.Rows, sh ServiceHandler) ([]byte, error) {
	data, err := sh.BuildGenericResponse(result)
	if err != nil {
		return nil, err
	}
	jsonResult, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonResult, nil
}

func RenderTemplate(w http.ResponseWriter, data interface{}, templatePath string) error {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("ERROR - Parse Template - Template Creation: %v\n", err)
	}
	err = t.Execute(w, data)
	if err != nil {
		return fmt.Errorf("ERROR - Render Template - Response Sending: %v\n", err)
	}
	return nil
}

func HandleUpdateRequest(r *http.Request, sh ServiceHandler, updateQuery string, updateID string, updateValues []interface{}) (*http.Request, error) {
	_, err := sh.GetDatasource().ExecuteWriteQuery(r.Context(), updateQuery, updateValues)
	if err != nil {
		return nil, fmt.Errorf("ERROR - UpdateElement - Database Insert: %v\n", err)
	}
	req, err := http.NewRequest("GET", r.URL.String()+updateID, nil)
	if err != nil {
		return nil, fmt.Errorf("ERROR - UpdateElement - Redirect Request: %v\n", err)
	}
	return req, nil
}

func HandleGetAllElementsRequest(r *http.Request, sh ServiceHandler, getAllQuery string) ([]map[string]interface{}, error) {
	result, err := sh.GetDatasource().ExecuteReadQuery(r.Context(), getAllQuery, nil)
	if err != nil {
		return nil, fmt.Errorf("ERROR - GetAllElements - Database Read: %v\n", err)
	}
	defer result.Close()
	response, err := sh.BuildGenericResponse(result)
	if err != nil {
		return nil, fmt.Errorf("ERROR - GetAllElements - Response Building: %v\n", err)
	}
	return response, nil
}

func HandleGetSingleElementRequest(r *http.Request, sh ServiceHandler, getSingleElementQuery string, elementId string) ([]map[string]interface{}, error) {
	result, err := sh.GetDatasource().ExecuteReadQuery(r.Context(), getSingleElementQuery, []interface{}{elementId})
	defer result.Close()
	if err != nil {
		return nil, fmt.Errorf("ERROR - GetSingleElement - Database Read: %v\n", err)
	}
	response, err := sh.BuildGenericResponse(result)
	if err != nil {
		return nil, fmt.Errorf("ERROR - GetSingleElement - Response Building: %v\n", err)
	}
	return response, nil
}

func HandleCreateElementRequest(r *http.Request, sh ServiceHandler, buildCreateQuery func(columns []string) string) (*http.Request, error) {
	data := make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("ERROR - CreateElement - Data Decode: %v\n", err)
	}
	values := make([]interface{}, 0)
	columns := make([]string, 0)
	for k, v := range data {
		values = append(values, v)
		columns = append(columns, k)
	}
	createElementQuery := buildCreateQuery(columns)
	result, err := sh.GetDatasource().ExecuteWriteQuery(r.Context(), createElementQuery, values)
	if err != nil {
		return nil, fmt.Errorf("ERROR - CreateElement - Database Insert: %v\n", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("ERROR - CreateElement - Database Result Parsing: %v\n", err)
	}
	req, err := http.NewRequest("GET", r.URL.String()+strconv.FormatInt(id, 10), nil)
	if err != nil {
		return nil, fmt.Errorf("ERROR - CreateElement - Redirect Request: %v\n", err)
	}
	return req, nil
}
