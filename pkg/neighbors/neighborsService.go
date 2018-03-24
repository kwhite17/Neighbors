package neighbors

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/kwhite17/Neighbors/pkg/database"
	"github.com/kwhite17/Neighbors/pkg/utils"
)

var templateDirectory = "../../templates/neighbors/"

type NeighborServiceHandler struct {
	Database database.Datasource
}

func (nsh NeighborServiceHandler) GetDatasource() database.Datasource {
	return nsh.Database
}

func (nsh NeighborServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathArray := strings.Split(strings.TrimPrefix(r.URL.Path, "/neighbors/"), "/")
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
	redirectReq, err := utils.HandleCreateElementRequest(r, nsh, nsh.buildCreateNeighborQuery)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	nsh.handleGetSingleNeighbor(w, redirectReq, strings.TrimPrefix(redirectReq.URL.Path, "/neighbors/"))
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
	response, err := utils.HandleGetSingleElementRequest(r, nsh, getSingleNeighborQuery, username)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = utils.RenderTemplate(w, response[0], templateDirectory+"neighbor.html")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

var getAllNeighborsQuery = "SELECT NeighborID, Username, Email, Phone, Location from neighbors"

func (nsh NeighborServiceHandler) handleGetAllNeighbors(w http.ResponseWriter, r *http.Request) {
	response, err := utils.HandleGetAllElementsRequest(r, nsh, getAllNeighborsQuery)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = utils.RenderTemplate(w, response, templateDirectory+"neighbors.html")
	if err != nil {
		log.Println(err)
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
	redirectReq, err := utils.HandleUpdateRequest(r, nsh, updateNeighborQuery, userData["NeighborID"], values)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	nsh.handleGetSingleNeighbor(w, redirectReq, userData["NeighborID"])
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

func (nsh NeighborServiceHandler) BuildGenericResponse(result *sql.Rows) ([]map[string]interface{}, error) {
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
