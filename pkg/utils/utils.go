package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

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

func RenderTemplate(w http.ResponseWriter, t *template.Template, data interface{}, function string) error {
	err := t.Execute(w, data)
	if err != nil {
		log.Printf("ERROR - "+function+" - Response Sending: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}

func HandleUpdateRequest(w http.ResponseWriter, r *http.Request, sh ServiceHandler, updateQuery string, updateID string, updateValues []interface{}) (*http.Request, error) {
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
