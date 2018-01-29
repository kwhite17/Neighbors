package items

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/kwhite17/Neighbors/database"
)

type ItemServiceHandler struct {
	Database database.Datasource
}

func (ish ItemServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathArray := strings.Split(strings.TrimPrefix(r.URL.Path, "/items/"), "/")
	switch pathArray[len(pathArray)-1] {
	case "new":
		log.Panic("NOT YET IMPLEMENTED")
	case "edit":
		log.Panic("NOT YET IMPLEMENTED")
	default:
		ish.requestMethodHandler(w, r)
	}
}

func (ish ItemServiceHandler) requestMethodHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		ish.handleCreateItem(w, r)
	case "GET":
		ish.handleGetItem(w, r)
	case "DELETE":
		ish.handleDeleteItem(w, r)
	case "PUT":
		ish.handleUpdateItem(w, r)
	default:
		w.Write([]byte("Invalid Request\n"))
	}
}

func (ish ItemServiceHandler) handleCreateItem(w http.ResponseWriter, r *http.Request) {
	itemData := make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&itemData)
	if err != nil {
		log.Printf("ERROR - CreateItem - Item Data Decode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	values := make([]interface{}, 0)
	columns := make([]string, 0)
	for k, v := range itemData {
		values = append(values, v)
		columns = append(columns, k)
	}
	createItemQuery := ish.buildCreateItemQuery(columns)
	log.Printf("DEBUG - CreateItem - Executing query: %s\n", createItemQuery)
	_, err = ish.Database.ExecuteWriteQuery(r.Context(), createItemQuery, values)
	if err != nil {
		log.Printf("ERROR - CreateItem - Database Insert: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	// rowID, err := queryResult.LastInsertId()
	// if err != nil {

	// }
	// handleGetSingleItem(w, r, string(rowID))
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
	if itemID := strings.TrimPrefix(r.URL.Path, "/items/"); len(itemID) > 0 {
		ish.handleGetSingleItem(w, r, itemID)
	} else {
		ish.handleGetAllItems(w, r)
	}
}

var getSingleItemQuery = "SELECT ItemID, Category, Gender, Size, Quantity, DropoffLocation from items where ItemID=?"

func (ish ItemServiceHandler) handleGetSingleItem(w http.ResponseWriter, r *http.Request, itemID string) {
	result, err := ish.Database.ExecuteReadQuery(r.Context(), getSingleItemQuery, []interface{}{itemID})
	defer result.Close()
	if err != nil {
		log.Printf("ERROR - GetItem - Database Read: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response, err := buildJsonResponse(result)
	if err != nil {
		log.Printf("ERROR - GetItem - ResponseBuilding: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(response))
	// tmpl, err := template.ParseFiles("templates/items.html")
	// if err != nil {
	// 	log.Printf("ERROR - GetItem - Template Rendering: %v\n", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// // tmpl.Execute(w, response)
}

var getAllItemsQuery = "SELECT ItemID, Category, Gender, Size, Quantity, DropoffLocation from items"

func (ish ItemServiceHandler) handleGetAllItems(w http.ResponseWriter, r *http.Request) {
	result, err := ish.Database.ExecuteReadQuery(r.Context(), getAllItemsQuery, nil)
	defer result.Close()
	if err != nil {
		log.Printf("ERROR - GetAllItems - Database Read: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response, err := buildJsonResponse(result)
	if err != nil {
		log.Printf("ERROR - GetItem - ResponseBuilding: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(response))
	// tmpl, err := template.ParseFiles("templates/items.html")
	// if err != nil {
	// 	log.Printf("ERROR - GetItem - Template Rendering: %v\n", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// tmpl.Execute(w, response)
}

func (ish ItemServiceHandler) handleUpdateItem(w http.ResponseWriter, r *http.Request) {
	itemData := make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&itemData)
	if err != nil {
		log.Printf("ERROR - UpdateItem - User Data Decode: %v\n", err)
		return
	}
	itemID := strings.TrimPrefix(r.URL.Path, "/items/")
	values := make([]interface{}, 0)
	columns := make([]string, 0)
	for k, v := range itemData {
		values = append(values, v)
		columns = append(columns, k)
	}
	updateItemQuery := ish.buildUpdateItemQuery(columns, itemID)
	log.Printf("DEBUG - UpdateItem - Executing query: %s\n", updateItemQuery)
	_, err = ish.Database.ExecuteWriteQuery(r.Context(), updateItemQuery, values)
	if err != nil {
		log.Printf("ERROR - UpdateItem - Database Update: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// rowID, err := queryResult.LastInsertId()
	// if err != nil {

	// }
	// handleGetSingleItem(w, r, string(rowID))
}

func (ish ItemServiceHandler) buildUpdateItemQuery(columns []string, itemID string) string {
	args := make([]string, 0)
	for i := 0; i < len(columns); i++ {
		args = append(args, columns[i]+"=?")
	}
	argString := strings.Join(args, ",")
	return "UPDATE items SET " + argString + " WHERE ItemID='" + itemID + "'"

}

var deleteNeighorQuery = "DELETE FROM items WHERE ItemID=?"

func (ish ItemServiceHandler) handleDeleteItem(w http.ResponseWriter, r *http.Request) {
	itemID := strings.TrimPrefix(r.URL.Path, "/items/")
	_, err := ish.Database.ExecuteWriteQuery(r.Context(), deleteNeighorQuery, []interface{}{itemID})
	if err != nil {
		log.Printf("ERROR - DeleteItem - Database Delete: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// handleGetAllItems(w, r)
}

func buildGenericResposne(result *sql.Rows) ([]map[string]interface{}, error) {
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
		responseItem["ItemId"] = id
		responseItem["Category"] = category
		responseItem["Gender"] = gender
		responseItem["Size"] = size
		responseItem["Quantity"] = quantity
		responseItem["DropoffLocation"] = dropoffLocation
		response = append(response, responseItem)
	}
	return response, nil
}

func buildJsonResponse(result *sql.Rows) ([]byte, error) {
	data, err := buildGenericResposne(result)
	if err != nil {
		return nil, err
	}
	jsonResult, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonResult, nil

}
