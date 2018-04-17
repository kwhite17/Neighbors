package users

import (
	"database/sql"
	"encoding/json"
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
	// authenticated, err := utils.IsAuthenticated(ush, w, r)
	// if !authenticated {
	// 	if err != nil {
	// 		log.Println(err)
	// 		err = nil
	// 	}
	// 	response, err := http.Get("/login/")
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 	}
	// 	page, err := ioutil.ReadAll(response.Body)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}
	// 	w.Write(page)
	// 	return
	// }
	pathArray := strings.Split(strings.TrimPrefix(r.URL.Path, serviceEndpoint), "/")
	switch pathArray[len(pathArray)-1] {
	case "new":
		err := utils.RenderTemplate(w, nil, serviceEndpoint+"new.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "edit":
		err := utils.RenderTemplate(w, nil, serviceEndpoint+"edit.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		ush.requestMethodHandler(w, r)
	}
}

func (ush UserServiceHandler) requestMethodHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		ush.handleCreateUser(w, r)
	case "GET":
		ush.handleGetUser(w, r)
	case "DELETE":
		ush.handleDeleteUser(w, r)
	case "PUT":
		ush.handleUpdateUser(w, r)
	default:
		w.Write([]byte("Invalid Request\n"))
	}
}

func (ush UserServiceHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
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

func (ush UserServiceHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
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

func (ush UserServiceHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	userData := make(map[string]string)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userData)
	if err != nil {
		log.Printf("ERROR - UpdateUser - User Data Decode: %v\n", err)
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
	updateNeighborQuery := buildUpdateUserQuery(columns, userData["ID"])
	redirectReq, err := utils.HandleUpdateRequest(r, ush, updateNeighborQuery, userData["ID"], values)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ush.handleGetSingleUser(w, redirectReq, userData["ID"])
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

func (ush UserServiceHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimPrefix(r.URL.Path, serviceEndpoint)
	w.Write([]byte("Deleting user data for " + username + "\n"))
	_, err := ush.Database.ExecuteWriteQuery(r.Context(), deleteUserQuery, []interface{}{username})
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
