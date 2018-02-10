package samaritans

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/kwhite17/Neighbors/database"
)

type SamaritanServiceHandler struct {
	Database database.Datasource
}

func (ssh SamaritanServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		ssh.handleCreateSamaritan(w, r)
	case "GET":
		ssh.handleGetSamaritan(w, r)
	case "DELETE":
		ssh.handleDeleteSamaritan(w, r)
	case "PUT":
		ssh.handleUpdateSamaritan(w, r)
	default:
		w.Write([]byte("Invalid Request\n"))
	}
}

func (ssh SamaritanServiceHandler) handleCreateSamaritan(w http.ResponseWriter, r *http.Request) {
	userData := make(map[string]string)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userData)
	if err != nil {
		log.Printf("ERROR - CreateSamaritan - User Data Decode: %v\n", err)
		return
	}
	values := make([]interface{}, 0)
	columns := make([]string, 0)
	for k, v := range userData {
		values = append(values, v)
		columns = append(columns, k)
	}
	createSamaritanQuery := buildCreateSamaritanQuery(columns)
	ssh.Database.ExecuteWriteQuery(r.Context(), createSamaritanQuery, values)
}

func buildCreateSamaritanQuery(columns []string) string {
	columnsString := strings.Join(columns, ",")
	args := make([]string, 0)
	for i := 0; i < len(columns); i++ {
		args = append(args, "?")
	}
	argString := strings.Join(args, ",")
	return "INSERT INTO samaritans (" + columnsString + ") VALUES (" + argString + ")"

}

func (ssh SamaritanServiceHandler) handleGetSamaritan(w http.ResponseWriter, r *http.Request) {
	if username := strings.TrimPrefix(r.URL.Path, "/samaritans/"); len(username) > 0 {
		ssh.handleGetSingleSamaritan(w, r, username)
	} else {
		ssh.handleGetAllSamaritans(w, r)
	}
}

var getSingleSamaritanQuery = "SELECT SamaritanID, Username, Email, Phone, Location from samaritans where SamaritanID=?"

func (ssh SamaritanServiceHandler) handleGetSingleSamaritan(w http.ResponseWriter, r *http.Request, username string) {
	log.Println("Fetching user: " + username)
	result, err := ssh.Database.ExecuteReadQuery(r.Context(), getSingleSamaritanQuery, []interface{}{username})
	if err != nil {
		log.Printf("ERROR - GetSingleSamaritan - Database Read: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer result.Close()
	response, err := buildJsonResposne(result)
	if err != nil {
		log.Printf("ERROR - GetSamaritan - ResponseBuilding: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

var getAllSamaritansQuery = "SELECT SamaritanID, Username, Email, Phone, Location from samaritans"

func (ssh SamaritanServiceHandler) handleGetAllSamaritans(w http.ResponseWriter, r *http.Request) {
	result, err := ssh.Database.ExecuteReadQuery(r.Context(), getAllSamaritansQuery, nil)
	if err != nil {
		log.Printf("ERROR - GetAllSamaritans - Database Read: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer result.Close()
	response, err := buildJsonResposne(result)
	if err != nil {
		log.Printf("ERROR - GetSamaritan - ResponseBuilding: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (ssh SamaritanServiceHandler) handleUpdateSamaritan(w http.ResponseWriter, r *http.Request) {
	userData := make(map[string]string)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userData)
	if err != nil {
		log.Printf("ERROR - UpdateSamaritan - User Data Decode: %v\n", err)
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
	updateSamaritanQuery := buildUpdateSamaritanQuery(columns, userData["Username"])
	ssh.Database.ExecuteWriteQuery(r.Context(), updateSamaritanQuery, values)
}

func buildUpdateSamaritanQuery(columns []string, username string) string {
	args := make([]string, 0)
	for i := 0; i < len(columns); i++ {
		args = append(args, columns[i]+"=?")
	}
	argString := strings.Join(args, ",")
	return "UPDATE samaritans SET " + argString + " WHERE Username='" + username + "'"

}

var deleteNeighorQuery = "DELETE FROM samaritans WHERE SamaritanID=?"

func (ssh SamaritanServiceHandler) handleDeleteSamaritan(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimPrefix(r.URL.Path, "/samaritans/")
	w.Write([]byte("Deleting user data for " + username + "\n"))
	ssh.Database.ExecuteWriteQuery(r.Context(), deleteNeighorQuery, []interface{}{username})

}

func buildJsonResposne(result *sql.Rows) ([]byte, error) {
	response := make([]map[string]interface{}, 0)
	for result.Next() {
		var samaritanID interface{}
		var username string
		var email interface{}
		var phone interface{}
		var location string
		responseItem := make(map[string]interface{})
		if err := result.Scan(&samaritanID, &username, &email, &phone, &location); err != nil {
			return nil, err
		}
		responseItem["SamaritanID"] = samaritanID
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
