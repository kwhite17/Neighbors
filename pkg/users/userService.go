package users

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/kwhite17/Neighbors/pkg/database"
	"github.com/kwhite17/Neighbors/pkg/utils"
)

var serviceEndpoint = "/users/"

type UserServiceHandler struct {
	Database database.Datasource
}

func (ush UserServiceHandler) GetDatasource() database.Datasource {
	return ush.Database
}

func (ush UserServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authRole, err := utils.IsAuthenticated(ush, w, r)
	if authRole == nil && r.Method != http.MethodPost {
		if err != nil {
			log.Println(err)
			err = nil
		}
		response, err := http.Get("/login/")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		page, err := ioutil.ReadAll(response.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(page)
		return
	}
	pathArray := strings.Split(strings.TrimPrefix(r.URL.Path, serviceEndpoint), "/")
	switch pathArray[len(pathArray)-1] {
	case "edit":
		err := utils.RenderTemplate(w, nil, serviceEndpoint+"edit.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		ush.requestMethodHandler(w, r, authRole)
	}
}

func (ush UserServiceHandler) requestMethodHandler(w http.ResponseWriter, r *http.Request, authRole *utils.AuthRole) {
	switch r.Method {
	case http.MethodPost:
		ush.handleCreateUser(w, r, authRole)
	case http.MethodGet:
		ush.handleGetUser(w, r, authRole)
	case http.MethodDelete:
		ush.handleDeleteUser(w, r, authRole)
	case http.MethodPut:
		ush.handleUpdateUser(w, r, authRole)
	default:
		w.Write([]byte("Invalid Request\n"))
	}
}

func (ush UserServiceHandler) handleCreateUser(w http.ResponseWriter, r *http.Request, authRole *utils.AuthRole) {
	if ush.isAuthorized(authRole, r, nil) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	redirectReq, err := utils.HandleCreateElementRequest(r, ush, ush.buildCreateUserQuery)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ush.handleGetSingleUser(w, redirectReq, strings.TrimPrefix(redirectReq.URL.Path, serviceEndpoint))
}

func (ush UserServiceHandler) buildCreateUserQuery(columns []string) string {
	columnsString := strings.Join(columns, ",")
	args := make([]string, 0)
	for i := 0; i < len(columns); i++ {
		args = append(args, "?")
	}
	argString := strings.Join(args, ",")
	return "INSERT INTO users (" + columnsString + ") VALUES (" + argString + ")"
}

func (ush UserServiceHandler) handleGetUser(w http.ResponseWriter, r *http.Request, authRole *utils.AuthRole) {
	if !ush.isAuthorized(authRole, r, nil) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if username := strings.TrimPrefix(r.URL.Path, serviceEndpoint); len(username) > 0 {
		ush.handleGetSingleUser(w, r, username)
	} else {
		ush.handleGetAllUsers(w, r)
	}
}

var getSingleNeighborQuery = "SELECT ID, Username, Email, Phone, Location, Role from users where ID=?"

func (ush UserServiceHandler) handleGetSingleUser(w http.ResponseWriter, r *http.Request, username string) {
	response, err := utils.HandleGetSingleElementRequest(r, ush, getSingleNeighborQuery, username)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = utils.RenderTemplate(w, response[0], serviceEndpoint+"user.html")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

var getAllNeighborsQuery = "SELECT ID, Username, Email, Phone, Location, Role from users"

func (ush UserServiceHandler) handleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	response, err := utils.HandleGetAllElementsRequest(r, ush, getAllNeighborsQuery)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = utils.RenderTemplate(w, response, serviceEndpoint+"users.html")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ush UserServiceHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request, authRole *utils.AuthRole) {
	userData := make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userData)
	if err != nil {
		log.Printf("ERROR - UpdateUser - User Data Decode: %v\n", err)
		return
	}
	if !ush.isAuthorized(authRole, r, userData) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	values := make([]interface{}, 0)
	columns := make([]string, 0)
	for k, v := range userData {
		if k != "ID" {
			values = append(values, v)
			columns = append(columns, k)
		}
	}
	updateNeighborQuery := buildUpdateUserQuery(columns, userData["ID"].(string))
	redirectReq, err := utils.HandleUpdateRequest(r, ush, updateNeighborQuery, userData["ID"].(string), values)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ush.handleGetSingleUser(w, redirectReq, userData["ID"].(string))
}

func buildUpdateUserQuery(columns []string, username string) string {
	args := make([]string, 0)
	for i := 0; i < len(columns); i++ {
		args = append(args, columns[i]+"=?")
	}
	argString := strings.Join(args, ",")
	return "UPDATE users SET " + argString + " WHERE ID='" + username + "'"
}

var deleteUserQuery = "DELETE FROM users WHERE ID=?"

func (ush UserServiceHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request, authRole *utils.AuthRole) {
	username := strings.TrimPrefix(r.URL.Path, serviceEndpoint)
	response, err := utils.HandleGetSingleElementRequest(r, ush, getSingleNeighborQuery, username)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ush.isAuthorized(authRole, r, response[0]) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write([]byte("Deleting user data for " + username + "\n"))
	_, err = ush.Database.ExecuteWriteQuery(r.Context(), deleteUserQuery, []interface{}{username})
	if err != nil {
		log.Printf("ERROR - DeleteUser - Database Delete: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ush UserServiceHandler) BuildGenericResponse(result *sql.Rows) ([]map[string]interface{}, error) {
	response := make([]map[string]interface{}, 0)
	for result.Next() {
		var ID interface{}
		var username string
		var email interface{}
		var phone interface{}
		var location string
		var role string
		responseItem := make(map[string]interface{})
		if err := result.Scan(&ID, &username, &email, &phone, &location, &role); err != nil {
			return nil, err
		}
		responseItem["ID"] = ID
		responseItem["Username"] = username
		responseItem["Email"] = email
		responseItem["Phone"] = phone
		responseItem["Location"] = location
		responseItem["Role"] = role
		response = append(response, responseItem)
	}
	return response, nil
}

func (ush UserServiceHandler) isAuthorized(role *utils.AuthRole, r *http.Request, data map[string]interface{}) bool {
	if role == nil {
		return false
	}
	switch r.Method {
	case http.MethodPost:
		fallthrough
	case http.MethodGet:
		return true
	case http.MethodPut:
		fallthrough
	case http.MethodDelete:
		return role.ID == data["ID"]
	default:
		return false
	}
}
