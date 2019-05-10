package items

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/kwhite17/Neighbors/pkg/utils"
)

var serviceEndpoint = "/items/"

type ItemServiceHandler struct {
	ItemManager   *ItemManager
	ItemRetriever *ItemRetriever
}

func (handler ItemServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authRole := &utils.AuthRole{}
	// if authRole == nil {
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
		itemID, err := strconv.ParseInt(pathArray[len(pathArray)-2], 10, 64)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		item, err := handler.ItemManager.GetItem(r.Context(), itemID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		t, err := handler.ItemRetriever.RetrieveEditEntityTemplate()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		t.Execute(w, item)
	default:
		handler.requestMethodHandler(w, r, authRole)
	}
}

func (handler ItemServiceHandler) requestMethodHandler(w http.ResponseWriter, r *http.Request, authRole *utils.AuthRole) {
	switch r.Method {
	case http.MethodPost:
		handler.handleCreateItem(w, r, authRole)
	case http.MethodGet:
		handler.handleGetItem(w, r, authRole)
	case http.MethodDelete:
		handler.handleDeleteItem(w, r, authRole)
	case http.MethodPut:
		handler.handleUpdateItem(w, r, authRole)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (handler ItemServiceHandler) handleCreateItem(w http.ResponseWriter, r *http.Request, authRole *utils.AuthRole) {
	if !handler.isAuthorized(authRole, r, nil) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	item := &Item{}
	err := json.NewDecoder(r.Body).Decode(item)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	itemID, _ := handler.ItemManager.WriteItem(r.Context(), item)
	item.ID = itemID

	json.NewEncoder(w).Encode(item)
}

func (handler ItemServiceHandler) handleGetItem(w http.ResponseWriter, r *http.Request, authRole *utils.AuthRole) {
	if !handler.isAuthorized(authRole, r, nil) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if item := strings.TrimPrefix(r.URL.Path, serviceEndpoint); len(item) > 0 {
		handler.handleGetSingleItem(w, r, item)
	} else {
		handler.handleGetAllItems(w, r)
	}
}

func (handler ItemServiceHandler) handleGetSingleItem(w http.ResponseWriter, r *http.Request, itemID string) {
	id, _ := strconv.ParseInt(itemID, 10, 64)
	item, _ := handler.ItemManager.GetItem(r.Context(), id)

	template, _ := handler.ItemRetriever.RetrieveSingleEntityTemplate()
	template.Execute(w, item)
}

func (handler ItemServiceHandler) handleGetAllItems(w http.ResponseWriter, r *http.Request) {
	items, _ := handler.ItemManager.GetItems(r.Context())

	template, _ := handler.ItemRetriever.RetrieveAllEntitiesTemplate()
	template.Execute(w, items)
}

func (handler ItemServiceHandler) handleUpdateItem(w http.ResponseWriter, r *http.Request, authRole *utils.AuthRole) {
	// if !handler.isAuthorized(authRole, r, nil) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	// item := &Item{}
	// err := json.NewDecoder(r.Body).Decode(item)
	// if err != nil {
	// 	log.Println(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// err = handler.ItemManager(r.Context(), shelter)
	// if err != nil {
	// 	log.Println(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// w.WriteHeader(http.StatusNoContent)
}

func (handler ItemServiceHandler) handleDeleteItem(w http.ResponseWriter, r *http.Request, authRole *utils.AuthRole) {
	shelterID := strings.TrimPrefix(r.URL.Path, serviceEndpoint)

	_, err := handler.ItemManager.DeleteItem(r.Context(), shelterID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler ItemServiceHandler) isAuthorized(role *utils.AuthRole, r *http.Request, data map[string]interface{}) bool {
	return true
	// if role == nil {
	// 	return false
	// }
	// switch r.Method {
	// case http.MethodGet:
	// 	pathArray := strings.Split(strings.TrimPrefix(r.URL.Path, serviceEndpoint), "/")
	// 	switch pathArray[len(pathArray)-1] {
	// 	case "edit":
	// 		itemID := pathArray[len(pathArray)-2]
	// 		itemData, err := utils.HandleGetSingleElementRequest(r, ish, getSingleItemQuery, itemID)
	// 		if err != nil {
	// 			log.Println(err)
	// 			return false
	// 		}
	// 		return itemData[0]["Requestor"] == role.ID && role.Role == "NEIGHBOR"
	// 	default:
	// 		return true
	// 	}
	// case http.MethodPost:
	// 	return role.Role == "NEIGHBOR"
	// case http.MethodDelete:
	// 	if data == nil {
	// 		return false
	// 	}
	// 	return data["Requestor"] == role.ID && role.Role == "SAMARITAN"
	// case http.MethodPut:
	// 	if data == nil {
	// 		return false
	// 	}
	// 	orderStatus, ok := data["OrderStatus"]
	// 	if !ok {
	// 		return false
	// 	}
	// 	switch orderStatus {
	// 	case "REQUESTED":
	// 		return role.Role == "SAMARITAN"
	// 	case "ASSIGNED":
	// 		fallthrough
	// 	case "PURCHASED":
	// 		return role.Role == "SAMARITAN" && data["Fulfiller"] == role.ID
	// 	case "DELIVERED":
	// 		return data["Requestor"] == role.ID && role.Role == "NEIGHBOR"
	// 	default:
	// 		return false
	// 	}
	// default:
	// 	return false
	// }
}
