package resources

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kwhite17/Neighbors/pkg/managers"
	"github.com/kwhite17/Neighbors/pkg/retrievers"
)

var itemsEndpoint = "/items/"

type ItemServiceHandler struct {
	ItemManager        *managers.ItemManager
	ItemRetriever      *retrievers.ItemRetriever
	UserSessionManager *managers.UserSessionManager
}

//maps can't be constants smdh
var USER_TYPE_TO_VALID_ITEM_UPDATE = map[managers.UserType][]managers.ItemStatus{
	managers.SHELTER:   []managers.ItemStatus{managers.CREATED, managers.RECEIVED},
	managers.SAMARITAN: []managers.ItemStatus{managers.CLAIMED, managers.DELIVERED},
}

func (handler ItemServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	isAuthorized, userSession := handler.isAuthorized(r)
	if !isAuthorized {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tplMap := map[string]interface{}{
		"UserSession": userSession,
	}

	pathArray := strings.Split(strings.TrimPrefix(r.URL.Path, itemsEndpoint), "/")
	switch pathArray[len(pathArray)-1] {
	case "new":
		t, err := handler.ItemRetriever.RetrieveCreateEntityTemplate()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		t.Execute(w, tplMap)
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

		tplMap["Item"] = item
		t.Execute(w, tplMap)
	default:
		handler.requestMethodHandler(w, r, userSession)
	}
}

func (handler ItemServiceHandler) requestMethodHandler(w http.ResponseWriter, r *http.Request, userSession *managers.UserSession) {
	switch r.Method {
	case http.MethodPost:
		handler.handleCreateItem(w, r, userSession)
	case http.MethodGet:
		handler.handleGetItem(w, r, userSession)
	case http.MethodDelete:
		handler.handleDeleteItem(w, r)
	case http.MethodPut:
		handler.handleUpdateItem(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (handler ItemServiceHandler) handleCreateItem(w http.ResponseWriter, r *http.Request, userSession *managers.UserSession) {
	item := &managers.Item{}
	err := json.NewDecoder(r.Body).Decode(item)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	item.ShelterID = userSession.UserID
	itemID, _ := handler.ItemManager.WriteItem(r.Context(), item)
	item.ID = itemID

	json.NewEncoder(w).Encode(item)
}

func (handler ItemServiceHandler) handleGetItem(w http.ResponseWriter, r *http.Request, userSession *managers.UserSession) {
	if item := strings.TrimPrefix(r.URL.Path, itemsEndpoint); len(item) > 0 {
		handler.handleGetSingleItem(w, r, item, userSession)
	} else {
		handler.handleGetAllItems(w, r, userSession)
	}
}

func (handler ItemServiceHandler) handleGetSingleItem(w http.ResponseWriter, r *http.Request, itemID string, userSession *managers.UserSession) {
	id, _ := strconv.ParseInt(itemID, 10, 64)
	item, _ := handler.ItemManager.GetItem(r.Context(), id)

	template, _ := handler.ItemRetriever.RetrieveSingleEntityTemplate()
	responseObject := make(map[string]interface{}, 0)
	responseObject["Item"] = item
	responseObject["UserSession"] = userSession
	template.Execute(w, responseObject)
}

func (handler ItemServiceHandler) handleGetAllItems(w http.ResponseWriter, r *http.Request, userSession *managers.UserSession) {
	items, _ := handler.ItemManager.GetItems(r.Context())

	template, _ := handler.ItemRetriever.RetrieveAllEntitiesTemplate()
	responseObject := make(map[string]interface{}, 0)
	responseObject["Items"] = items
	responseObject["UserSession"] = userSession
	template.Execute(w, responseObject)
}

func (handler ItemServiceHandler) handleUpdateItem(w http.ResponseWriter, r *http.Request) {
	item := &managers.Item{}
	err := json.NewDecoder(r.Body).Decode(item)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = handler.ItemManager.UpdateItem(r.Context(), item)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler ItemServiceHandler) handleDeleteItem(w http.ResponseWriter, r *http.Request) {
	shelterID := strings.TrimPrefix(r.URL.Path, itemsEndpoint)

	_, err := handler.ItemManager.DeleteItem(r.Context(), shelterID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler ItemServiceHandler) isAuthorized(r *http.Request) (bool, *managers.UserSession) {
	cookie, _ := r.Cookie("NeighborsAuth")
	pathArray := strings.Split(strings.TrimPrefix(r.URL.Path, itemsEndpoint), "/")
	switch r.Method {
	case http.MethodPost:
		if cookie == nil {
			return false, nil
		}

		userSession, err := handler.UserSessionManager.GetUserSession(r.Context(), cookie.Value)
		if err != nil {
			log.Println(err)
			return false, nil
		}

		return time.Now().After(cookie.Expires), userSession
	case http.MethodPut:
		fallthrough
	case http.MethodDelete:
		userSession, err := handler.UserSessionManager.GetUserSession(r.Context(), cookie.Value)
		if err != nil {
			log.Println(err)
			return false, nil
		}

		itemID, err := strconv.ParseInt(pathArray[len(pathArray)-1], 10, strconv.IntSize)
		if err != nil {
			log.Println(err)
			return false, nil
		}

		item, err := handler.ItemManager.GetItem(r.Context(), itemID)
		if err != nil {
			log.Println(err)
			return false, nil
		}

		return item.ShelterID == userSession.UserID, nil
	case http.MethodGet:
		var err error
		userSession := &managers.UserSession{}
		if cookie != nil {
			userSession, err = handler.UserSessionManager.GetUserSession(r.Context(), cookie.Value)
			if err != nil && err != sql.ErrNoRows {
				log.Println(err)
				return false, userSession
			}
		}

		if pathArray[len(pathArray)-1] == "edit" {
			itemID, err := strconv.ParseInt(pathArray[len(pathArray)-2], 10, strconv.IntSize)
			if err != nil {
				log.Println(err)
				return false, userSession
			}

			item, err := handler.ItemManager.GetItem(r.Context(), itemID)
			if err != nil {
				log.Println(err)
				return false, userSession
			}

			return item.ShelterID == userSession.UserID, userSession
		}
		return true, userSession
	default:
		return false, nil
	}
}
