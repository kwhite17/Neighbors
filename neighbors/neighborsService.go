package neighbors

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/kwhite17/Neighbors/database"
)

type NeighborServiceHandler struct {
	Database database.Datasource
}

func (nsh NeighborServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	_, err = nsh.Database.ExecuteWriteQuery(r.Context(), createNeighborQuery, values)
	if err != nil {
		log.Printf("ERROR - CreateNeighbor - Database Insert: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
	response, err := buildJsonResposne(result)
	if err != nil {
		log.Printf("ERROR - GetNeighbor - ResponseBuilding: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)
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
	response, err := buildJsonResposne(result)
	if err != nil {
		log.Printf("ERROR - GetNeighbor - ResponseBuilding: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)
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
		if k != "Username" {
			values = append(values, v)
			columns = append(columns, k)
		}
	}
	updateNeighborQuery := buildUpdateNeighborQuery(columns, userData["Username"])
	_, err = nsh.Database.ExecuteWriteQuery(r.Context(), updateNeighborQuery, values)
	if err != nil {
		log.Printf("ERROR - UpdateNeighbor - Database Update: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func buildUpdateNeighborQuery(columns []string, username string) string {
	args := make([]string, 0)
	for i := 0; i < len(columns); i++ {
		args = append(args, columns[i]+"=?")
	}
	argString := strings.Join(args, ",")
	return "UPDATE neighbors SET " + argString + " WHERE Username='" + username + "'"

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

func buildJsonResposne(result *sql.Rows) ([]byte, error) {
	response := make([]map[string]interface{}, 0)
	for result.Next() {
		var neighborId interface{}
		var username string
		var email interface{}
		var phone interface{}
		var location string
		responseItem := make(map[string]interface{})
		if err := result.Scan(&neighborId, &username, &email, &phone, &location); err != nil {
			return nil, err
		}
		responseItem["NeighborID"] = neighborId
		responseItem["Username"] = username
		responseItem["Email"] = email
		responseItem["Phone"] = phone
		responseItem["Location"] = location
		response = append(response, responseItem)
	}
	jsonResult, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return jsonResult, nil
}
