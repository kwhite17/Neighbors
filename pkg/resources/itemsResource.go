package resources

import (
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
	ItemManager           *managers.ItemManager
	ItemRetriever         *retrievers.ItemRetriever
	ShelterSessionManager *managers.ShelterSessionManager
}

func (handler ItemServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("NeighborsAuth")
	if !handler.isAuthorized(r, cookie) {
		w.WriteHeader(http.StatusUnauthorized)
		return
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
		t.Execute(w, nil)
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
		handler.requestMethodHandler(w, r)
	}
}

func (handler ItemServiceHandler) requestMethodHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handler.handleCreateItem(w, r)
	case http.MethodGet:
		handler.handleGetItem(w, r)
	case http.MethodDelete:
		handler.handleDeleteItem(w, r)
	case http.MethodPut:
		handler.handleUpdateItem(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (handler ItemServiceHandler) handleCreateItem(w http.ResponseWriter, r *http.Request) {
	item := &managers.Item{}
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

func (handler ItemServiceHandler) handleGetItem(w http.ResponseWriter, r *http.Request) {
	if item := strings.TrimPrefix(r.URL.Path, itemsEndpoint); len(item) > 0 {
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

func (handler ItemServiceHandler) isAuthorized(r *http.Request, cookie *http.Cookie) bool {
	pathArray := strings.Split(strings.TrimPrefix(r.URL.Path, itemsEndpoint), "/")
	switch r.Method {
	case http.MethodPost:
		return cookie != nil && time.Now().After(cookie.Expires)
	case http.MethodPut:
		fallthrough
	case http.MethodDelete:
		shelterSession, err := handler.ShelterSessionManager.GetShelterSession(r.Context(), cookie.Value)
		if err != nil {
			log.Println(err)
			return false
		}

		itemID, err := strconv.ParseInt(pathArray[len(pathArray)-2], 10, strconv.IntSize)
		if err != nil {
			log.Println(err)
			return false
		}

		item, err := handler.ItemManager.GetItem(r.Context(), itemID)
		if err != nil {
			log.Println(err)
			return false
		}

		return item.ShelterID == shelterSession.ShelterID
	case http.MethodGet:
		if pathArray[len(pathArray)-1] == "edit" {
			shelterSession, err := handler.ShelterSessionManager.GetShelterSession(r.Context(), cookie.Value)
			if err != nil {
				log.Println(err)
				return false
			}

			itemID, err := strconv.ParseInt(pathArray[len(pathArray)-2], 10, strconv.IntSize)
			if err != nil {
				log.Println(err)
				return false
			}

			item, err := handler.ItemManager.GetItem(r.Context(), itemID)
			if err != nil {
				log.Println(err)
				return false
			}

			return item.ShelterID == shelterSession.ShelterID
		}
		return true
	default:
		return false
	}
}
