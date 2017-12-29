package neighbors

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/kwhite17/Neighbors/database"
)

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handleCreateNeighbor(w, r)
	case "GET":
		handleGetNeighbor(w, r)
	case "DELETE":
		handleDeleteNeighbor(w, r)
	case "PUT":
		handleUpdateNeighbor(w, r)
	default:
		w.Write([]byte("Invalid Request\n"))
	}
}

func handleCreateNeighbor(w http.ResponseWriter, r *http.Request) {
	userData := make(map[string]string)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userData)
	if err != nil {
		log.Printf("ERROR - CreateNeighbor - User Data Decode: %v\n", err)
		return
	}
	values := make([]interface{}, 0)
	columns := make([]string, 0)
	for k, v := range userData {
		values = append(values, v)
		columns = append(columns, k)
	}
	createNeighborQuery := buildCreateNeighborQuery(columns)
	log.Printf("DEBUG - CreateNeighbor - Executing query: %s\n", createNeighborQuery)
	database.ExecuteWriteQuery(r.Context(), createNeighborQuery, values)
}

func buildCreateNeighborQuery(columns []string) string {
	columnsString := strings.Join(columns, ",")
	args := make([]string, 0)
	for i := 0; i < len(columns); i++ {
		args = append(args, "?")
	}
	argString := strings.Join(args, ",")
	return "INSERT INTO neighbors (" + columnsString + ") VALUES (" + argString + ")"

}

func handleGetNeighbor(w http.ResponseWriter, r *http.Request) {
	if username := strings.TrimPrefix(r.URL.Path, "/neighbors/"); len(username) > 0 {
		handleGetSingleNeighbor(w, r, username)
	} else {
		handleGetAllNeighbors(w, r)
	}
}

var getSingleNeighborQuery = "SELECT Username, Email, Phone, Location from neighbors where Username=?"

func handleGetSingleNeighbor(w http.ResponseWriter, r *http.Request, username string) {
	log.Println("Fetching user: " + username)
	result := database.ExecuteReadQuery(r.Context(), getSingleNeighborQuery, []interface{}{username})
	defer result.Close()
	response, err := buildJsonResposne(result)
	if err != nil {
		log.Printf("ERROR - GetNeighbor - ResponseBuilding: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

var getAllNeighborsQuery = "SELECT Username, Email, Phone, Location from neighbors"

func handleGetAllNeighbors(w http.ResponseWriter, r *http.Request) {
	result := database.ExecuteReadQuery(r.Context(), getAllNeighborsQuery, nil)
	defer result.Close()
	response, err := buildJsonResposne(result)
	if err != nil {
		log.Printf("ERROR - GetNeighbor - ResponseBuilding: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func handleUpdateNeighbor(w http.ResponseWriter, r *http.Request) {
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
	log.Printf("DEBUG - UpdateNeighbor - Executing query: %s\n", updateNeighborQuery)
	database.ExecuteWriteQuery(r.Context(), updateNeighborQuery, values)
}

func buildUpdateNeighborQuery(columns []string, username string) string {
	args := make([]string, 0)
	for i := 0; i < len(columns); i++ {
		args = append(args, columns[i]+"=?")
	}
	argString := strings.Join(args, ",")
	return "UPDATE neighbors SET " + argString + " WHERE Username='" + username + "'"

}

var deleteNeighorQuery = "DELETE FROM neighbors WHERE Username=?"

func handleDeleteNeighbor(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimPrefix(r.URL.Path, "/neighbors/")
	w.Write([]byte("Deleting user data for " + username + "\n"))
	database.ExecuteWriteQuery(r.Context(), deleteNeighorQuery, []interface{}{username})

}

func buildJsonResposne(result *sql.Rows) ([]byte, error) {
	response := make([]map[string]interface{}, 0)
	for result.Next() {
		var username string
		var email interface{}
		var phone interface{}
		var location string
		responseItem := make(map[string]interface{})
		if err := result.Scan(&username, &email, &phone, &location); err != nil {
			return nil, err
		}
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
