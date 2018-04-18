package items

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

var serviceEndpoint = "/items/"

type ItemServiceHandler struct {
	Database database.Datasource
}

func (ish ItemServiceHandler) GetDatasource() database.Datasource {
	return ish.Database
}

func (ish ItemServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authenticated, err := utils.IsAuthenticated(ish, w, r)
	if !authenticated {
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
		ish.requestMethodHandler(w, r)
	}
}

func (ish ItemServiceHandler) requestMethodHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		ish.handleCreateItem(w, r)
	case http.MethodGet:
		ish.handleGetItem(w, r)
	case http.MethodDelete:
		ish.handleDeleteItem(w, r)
	case http.MethodPut:
		ish.handleUpdateItem(w, r)
	default:
		w.Write([]byte("Invalid Request\n"))
	}
}

func (ish ItemServiceHandler) handleCreateItem(w http.ResponseWriter, r *http.Request) {
	redirectReq, err := utils.HandleCreateElementRequest(r, ish, ish.buildCreateItemQuery)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ish.handleGetSingleItem(w, redirectReq, strings.TrimPrefix(redirectReq.URL.Path, serviceEndpoint))
}

func (ish ItemServiceHandler) buildCreateItemQuery(columns []string) string {
	columnsString := strings.Join(columns, ",")
	args := make([]string, 0)
	for i := 0; i < len(columns); i++ {
		args = append(args, "?")
	}
	argString := strings.Join(args, ",")
	return "INSERT INTO items (" + columnsString + ") VALUES (" + argString + ")"

}

func (ish ItemServiceHandler) handleGetItem(w http.ResponseWriter, r *http.Request) {
	if itemID := strings.TrimPrefix(r.URL.Path, serviceEndpoint); len(itemID) > 0 {
		ish.handleGetSingleItem(w, r, itemID)
	} else {
		ish.handleGetAllItems(w, r)
	}
}

var getSingleItemQuery = "SELECT ID, Category, Gender, Size, Quantity, DropoffLocation from items where ID=?"

func (ish ItemServiceHandler) handleGetSingleItem(w http.ResponseWriter, r *http.Request, itemID string) {
	response, err := utils.HandleGetSingleElementRequest(r, ish, getSingleItemQuery, itemID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = utils.RenderTemplate(w, response[0], serviceEndpoint+"item.html")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

var getAllItemsQuery = "SELECT ID, Category, Gender, Size, Quantity, DropoffLocation from items"

func (ish ItemServiceHandler) handleGetAllItems(w http.ResponseWriter, r *http.Request) {
	response, err := utils.HandleGetAllElementsRequest(r, ish, getAllItemsQuery)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = utils.RenderTemplate(w, response, serviceEndpoint+"items.html")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ish ItemServiceHandler) handleUpdateItem(w http.ResponseWriter, r *http.Request) {
	itemData := make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&itemData)
	if err != nil {
		log.Printf("ERROR - UpdateItem - Item Data Decode: %v\n", err)
		return
	}
	values := make([]interface{}, 0)
	columns := make([]string, 0)
	for k, v := range itemData {
		if k != "ID" {
			values = append(values, v)
			columns = append(columns, k)
		}
	}
	itemID := itemData["ID"].(string)
	updateItemQuery := ish.buildUpdateItemQuery(columns, itemID)
	redirectReq, err := utils.HandleUpdateRequest(r, ish, updateItemQuery, itemID, values)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ish.handleGetSingleItem(w, redirectReq, itemID)

}

func (ish ItemServiceHandler) buildUpdateItemQuery(columns []string, itemID string) string {
	args := make([]string, 0)
	for i := 0; i < len(columns); i++ {
		args = append(args, columns[i]+"=?")
	}
	argString := strings.Join(args, ",")
	return "UPDATE items SET " + argString + " WHERE ID='" + itemID + "'"

}

var deleteNeighorQuery = "DELETE FROM items WHERE ID=?"

func (ish ItemServiceHandler) handleDeleteItem(w http.ResponseWriter, r *http.Request) {
	itemID := strings.TrimPrefix(r.URL.Path, serviceEndpoint)
	_, err := ish.Database.ExecuteWriteQuery(r.Context(), deleteNeighorQuery, []interface{}{itemID})
	if err != nil {
		log.Printf("ERROR - DeleteItem - Database Delete: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// handleGetAllItems(w, r)
}

func (ish ItemServiceHandler) BuildGenericResponse(result *sql.Rows) ([]map[string]interface{}, error) {
	response := make([]map[string]interface{}, 0)
	for result.Next() {
		var id int
		var category string
		var gender interface{}
		var size string
		var quantity int
		var dropoffLocation string
		responseItem := make(map[string]interface{})
		if err := result.Scan(&id, &category, &gender, &size, &quantity,
			&dropoffLocation); err != nil {
			return nil, err
		}
		responseItem["ID"] = id
		responseItem["Category"] = category
		responseItem["Gender"] = gender
		responseItem["Size"] = size
		responseItem["Quantity"] = quantity
		responseItem["DropoffLocation"] = dropoffLocation
		response = append(response, responseItem)
	}
	return response, nil
}
