package neighbors

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/kwhite17/Neighbors/database"
)

type NeighborServiceHandler struct {
	Database database.Datasource
}

func (nsh NeighborServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathArray := strings.Split(strings.TrimPrefix(r.URL.Path, "/neighbors/"), "/")
	switch pathArray[len(pathArray)-1] {
	case "new":
		t, err := template.ParseFiles("../templates/neighbors/new.html")
		if err != nil {
			log.Printf("ERROR - NewNeighbor - Template Rendering: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = t.Execute(w, nil)
		if err != nil {
			log.Printf("ERROR - NewNeighbor - Response Sending: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "edit":
		t, err := template.ParseFiles("../templates/neighbors/edit.html")
		if err != nil {
			log.Printf("ERROR - EditNeighbor - Template Rendering: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = t.Execute(w, nil)
		if err != nil {
			log.Printf("ERROR - EditNeighbor - Response Sending: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		nsh.requestMethodHandler(w, r)
	}
}

func (nsh NeighborServiceHandler) requestMethodHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		nsh.handleCreateNeighbor(w, r)
	case "GET":
		nsh.handleGetNeighbor(w, r)
	case "DELETE":
		nsh.handleDeleteNeighbor(w, r)
	case "PUT":
		nsh.handleUpdateNeighbor(w, r)
	default:
		w.Write([]byte("Invalid Request\n"))
	}
}

func (nsh NeighborServiceHandler) handleCreateNeighbor(w http.ResponseWriter, r *http.Request) {
	userData := make(map[string]string)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userData)
	if err != nil {
		log.Printf("ERROR - CreateNeighbor - User Data Decode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	values := make([]interface{}, 0)
	columns := make([]string, 0)
	for k, v := range userData {
		values = append(values, v)
		columns = append(columns, k)
	}
	createNeighborQuery := nsh.buildCreateNeighborQuery(columns)
	result, err := nsh.Database.ExecuteWriteQuery(r.Context(), createNeighborQuery, values)
	if err != nil {
		log.Printf("ERROR - CreateNeighbor - Database Insert: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("ERROR - CreateNeighbor - Database Result Parsing: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req, err := http.NewRequest("GET", r.URL.String()+strconv.FormatInt(id, 10), nil)
	if err != nil {
		log.Printf("ERROR - UpdateNeighbor - Redirect Request: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	nsh.handleGetSingleNeighbor(w, req, strconv.FormatInt(id, 10))
}

func (nsh NeighborServiceHandler) buildCreateNeighborQuery(columns []string) string {
	columnsString := strings.Join(columns, ",")
	args := make([]string, 0)
	for i := 0; i < len(columns); i++ {
		args = append(args, "?")
	}
	argString := strings.Join(args, ",")
	return "INSERT INTO neighbors (" + columnsString + ") VALUES (" + argString + ")"

}

func (nsh NeighborServiceHandler) handleGetNeighbor(w http.ResponseWriter, r *http.Request) {
	if username := strings.TrimPrefix(r.URL.Path, "/neighbors/"); len(username) > 0 {
		nsh.handleGetSingleNeighbor(w, r, username)
	} else {
		nsh.handleGetAllNeighbors(w, r)
	}
}

var getSingleNeighborQuery = "SELECT NeighborID, Username, Email, Phone, Location from neighbors where NeighborID=?"

func (nsh NeighborServiceHandler) handleGetSingleNeighbor(w http.ResponseWriter, r *http.Request, username string) {
	result, err := nsh.Database.ExecuteReadQuery(r.Context(), getSingleNeighborQuery, []interface{}{username})
	if err != nil {
		log.Printf("ERROR - GetSingleNeighbor - Database Read: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer result.Close()
	response, err := buildGenericResponse(result)
	if err != nil {
		log.Printf("ERROR - GetNeighbor - ResponseBuilding: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t, err := template.ParseFiles("../templates/neighbors/neighbor.html")
	if err != nil {
		log.Printf("ERROR - GetNeighbor - Template Creation: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, response[0])
	if err != nil {
		log.Printf("ERROR - GetNeighbor - Template Resolution: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

var getAllNeighborsQuery = "SELECT NeighborID, Username, Email, Phone, Location from neighbors"

func (nsh NeighborServiceHandler) handleGetAllNeighbors(w http.ResponseWriter, r *http.Request) {
	result, err := nsh.Database.ExecuteReadQuery(r.Context(), getAllNeighborsQuery, nil)
	if err != nil {
		log.Printf("ERROR - GetAllNeighbors - Database Read: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer result.Close()
	response, err := buildGenericResponse(result)
	if err != nil {
		log.Printf("ERROR - GetNeighbor - ResponseBuilding: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t, err := template.ParseFiles("../templates/neighbors/neighbors.html")
	if err != nil {
		log.Printf("ERROR - GetNeighbor - Template Creation: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, response)
	if err != nil {
		log.Printf("ERROR - GetNeighbor - Template Resolution: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (nsh NeighborServiceHandler) handleUpdateNeighbor(w http.ResponseWriter, r *http.Request) {
	userData := make(map[string]string)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userData)
	if err != nil {
		log.Printf("ERROR - UpdateNeighbor - User Data Decode: %v\n", err)
		return
	}
	values := make([]interface{}, 0)
	columns := make([]string, 0)
	for k, v := range userData {
		if k != "NeighborID" {
			values = append(values, v)
			columns = append(columns, k)
		}
	}
	updateNeighborQuery := buildUpdateNeighborQuery(columns, userData["NeighborID"])
	_, err = nsh.Database.ExecuteWriteQuery(r.Context(), updateNeighborQuery, values)
	if err != nil {
		log.Printf("ERROR - UpdateNeighbor - Database Update: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req, err := http.NewRequest("GET", r.URL.String()+userData["NeighborID"], nil)
	if err != nil {
		log.Printf("ERROR - UpdateNeighbor - Redirect Request: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	nsh.handleGetSingleNeighbor(w, req, userData["NeighborID"])
}

func buildUpdateNeighborQuery(columns []string, username string) string {
	args := make([]string, 0)
	for i := 0; i < len(columns); i++ {
		args = append(args, columns[i]+"=?")
	}
	argString := strings.Join(args, ",")
	return "UPDATE neighbors SET " + argString + " WHERE NeighborID='" + username + "'"

}

var deleteNeighorQuery = "DELETE FROM neighbors WHERE NeighborID=?"

func (nsh NeighborServiceHandler) handleDeleteNeighbor(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimPrefix(r.URL.Path, "/neighbors/")
	w.Write([]byte("Deleting user data for " + username + "\n"))
	_, err := nsh.Database.ExecuteWriteQuery(r.Context(), deleteNeighorQuery, []interface{}{username})
	if err != nil {
		log.Printf("ERROR - DeleteNeighbor - Database Delete: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func buildGenericResponse(result *sql.Rows) ([]map[string]interface{}, error) {
	response := make([]map[string]interface{}, 0)
	for result.Next() {
		var neighborID interface{}
		var username string
		var email interface{}
		var phone interface{}
		var location string
		responseItem := make(map[string]interface{})
		if err := result.Scan(&neighborID, &username, &email, &phone, &location); err != nil {
			return nil, err
		}
		responseItem["NeighborID"] = neighborID
		responseItem["Username"] = username
		responseItem["Email"] = email
		responseItem["Phone"] = phone
		responseItem["Location"] = location
		response = append(response, responseItem)
	}
	return response, nil
}

func buildJsonResponse(result *sql.Rows) ([]byte, error) {
	data, err := buildGenericResponse(result)
	if err != nil {
		return nil, err
	}
	jsonResult, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonResult, nil

}
