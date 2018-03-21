package samaritans

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/kwhite17/Neighbors/pkg/utils"

	"github.com/kwhite17/Neighbors/pkg/database"
)

var templateDirectory = "../../templates/samaritans/"

type SamaritanServiceHandler struct {
	Database database.Datasource
}

func (ssh SamaritanServiceHandler) GetDatasource() database.Datasource {
	return ssh.Database
}

func (ssh SamaritanServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathArray := strings.Split(strings.TrimPrefix(r.URL.Path, "/samaritans/"), "/")
	switch pathArray[len(pathArray)-1] {
	case "new":
		err := utils.RenderTemplate(w, nil, templateDirectory+"new.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "edit":
		err := utils.RenderTemplate(w, nil, templateDirectory+"edit.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		ssh.requestMethodHandler(w, r)
	}
}

func (ssh SamaritanServiceHandler) requestMethodHandler(w http.ResponseWriter, r *http.Request) {
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
	result, err := ssh.Database.ExecuteWriteQuery(r.Context(), createSamaritanQuery, values)
	if err != nil {
		log.Printf("ERROR - CreateSamaritan - Database Insert: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("ERROR - CreateSamaritan - Database Result Parsing: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req, err := http.NewRequest("GET", r.URL.String()+strconv.FormatInt(id, 10), nil)
	if err != nil {
		log.Printf("ERROR - CreateSamaritan - Redirect Request: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ssh.handleGetSingleSamaritan(w, req, strconv.FormatInt(id, 10))
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
	response, err := utils.HandleGetSingleElementRequest(r, ssh, getSingleSamaritanQuery, username)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = utils.RenderTemplate(w, response[0], templateDirectory+"samaritan.html")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

var getAllSamaritansQuery = "SELECT SamaritanID, Username, Email, Phone, Location from samaritans"

func (ssh SamaritanServiceHandler) handleGetAllSamaritans(w http.ResponseWriter, r *http.Request) {
	response, err := utils.HandleGetAllElementsRequest(r, ssh, getAllSamaritansQuery)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = utils.RenderTemplate(w, response, templateDirectory+"samaritans.html")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
		if k != "SamaritanID" {
			values = append(values, v)
			columns = append(columns, k)
		}
	}
	updateSamaritanQuery := buildUpdateSamaritanQuery(columns, userData["SamaritanID"])
	redirectReq, err := utils.HandleUpdateRequest(r, ssh, updateSamaritanQuery, userData["SamaritanID"], values)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ssh.handleGetSingleSamaritan(w, redirectReq, userData["SamaritanID"])
}

func buildUpdateSamaritanQuery(columns []string, username string) string {
	args := make([]string, 0)
	for i := 0; i < len(columns); i++ {
		args = append(args, columns[i]+"=?")
	}
	argString := strings.Join(args, ",")
	return "UPDATE samaritans SET " + argString + " WHERE SamaritanID='" + username + "'"

}

var deleteNeighorQuery = "DELETE FROM samaritans WHERE SamaritanID=?"

func (ssh SamaritanServiceHandler) handleDeleteSamaritan(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimPrefix(r.URL.Path, "/samaritans/")
	w.Write([]byte("Deleting user data for " + username + "\n"))
	ssh.Database.ExecuteWriteQuery(r.Context(), deleteNeighorQuery, []interface{}{username})

}

func (ssh SamaritanServiceHandler) BuildGenericResponse(result *sql.Rows) ([]map[string]interface{}, error) {
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
	return response, nil
}
