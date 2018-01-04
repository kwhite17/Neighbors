package items

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
		handleCreateItem(w, r)
	case "GET":
		handleGetItem(w, r)
	case "DELETE":
		handleDeleteItem(w, r)
	case "PUT":
		handleUpdateItem(w, r)
	default:
		w.Write([]byte("Invalid Request\n"))
	}
}

func handleCreateItem(w http.ResponseWriter, r *http.Request) {
	itemData := make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&itemData)
	if err != nil {
		log.Printf("ERROR - CreateItem - Item Data Decode: %v\n", err)
		return
	}
	values := make([]interface{}, 0)
	columns := make([]string, 0)
	for k, v := range itemData {
		values = append(values, v)
		columns = append(columns, k)
	}
	createItemQuery := buildCreateItemQuery(columns)
	log.Printf("DEBUG - CreateItem - Executing query: %s\n", createItemQuery)
	database.ExecuteWriteQuery(r.Context(), createItemQuery, values)
}

func buildCreateItemQuery(columns []string) string {
	columnsString := strings.Join(columns, ",")
	args := make([]string, 0)
	for i := 0; i < len(columns); i++ {
		args = append(args, "?")
	}
	argString := strings.Join(args, ",")
	return "INSERT INTO items (" + columnsString + ") VALUES (" + argString + ")"

}

func handleGetItem(w http.ResponseWriter, r *http.Request) {
	if itemID := strings.TrimPrefix(r.URL.Path, "/items/"); len(itemID) > 0 {
		handleGetSingleItem(w, r, itemID)
	} else {
		handleGetAllItems(w, r)
	}
}

var getSingleItemQuery = "SELECT ItemID, Category, Gender, Size, Quantity, DropoffLocation from items where ItemID=?"

func handleGetSingleItem(w http.ResponseWriter, r *http.Request, itemID string) {
	result := database.ExecuteReadQuery(r.Context(), getSingleItemQuery, []interface{}{itemID})
	defer result.Close()
	response, err := buildJsonResposne(result)
	if err != nil {
		log.Printf("ERROR - GetItem - ResponseBuilding: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

var getAllItemsQuery = "SELECT ItemID, Category, Gender, Size, Quantity, DropoffLocation from items"

func handleGetAllItems(w http.ResponseWriter, r *http.Request) {
	result := database.ExecuteReadQuery(r.Context(), getAllItemsQuery, nil)
	defer result.Close()
	response, err := buildJsonResposne(result)
	if err != nil {
		log.Printf("ERROR - GetItem - ResponseBuilding: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func handleUpdateItem(w http.ResponseWriter, r *http.Request) {
	itemData := make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&itemData)
	if err != nil {
		log.Printf("ERROR - UpdateItem - User Data Decode: %v\n", err)
		return
	}
	values := make([]interface{}, 0)
	columns := make([]string, 0)
	for k, v := range itemData {
		if k != "ItemID" {
			values = append(values, v)
			columns = append(columns, k)
		}
	}
	updateItemQuery := buildUpdateItemQuery(columns, itemData["ItemID"].(string))
	log.Printf("DEBUG - UpdateItem - Executing query: %s\n", updateItemQuery)
	database.ExecuteWriteQuery(r.Context(), updateItemQuery, values)
}

func buildUpdateItemQuery(columns []string, itemID string) string {
	args := make([]string, 0)
	for i := 0; i < len(columns); i++ {
		args = append(args, columns[i]+"=?")
	}
	argString := strings.Join(args, ",")
	return "UPDATE items SET " + argString + " WHERE ItemID='" + itemID + "'"

}

var deleteNeighorQuery = "DELETE FROM items WHERE ItemID=?"

func handleDeleteItem(w http.ResponseWriter, r *http.Request) {
	itemID := strings.TrimPrefix(r.URL.Path, "/items/")
	database.ExecuteWriteQuery(r.Context(), deleteNeighorQuery, []interface{}{itemID})
}

func buildJsonResposne(result *sql.Rows) ([]byte, error) {
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
	jsonResult, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return jsonResult, nil
}
