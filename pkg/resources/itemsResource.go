package resources

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/kwhite17/Neighbors/pkg/email"

	"github.com/kwhite17/Neighbors/pkg/managers"
	"github.com/kwhite17/Neighbors/pkg/retrievers"
)

var itemsEndpoint = "/items/"

type ItemServiceHandler struct {
	ItemManager        *managers.ItemManager
	ItemRetriever      *retrievers.ItemRetriever
	UserSessionManager managers.SessionManger
	EmailSender        email.EmailSender
}

func (handler ItemServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	isAuthorized, userSession := handler.isAuthorized(r)
	if !isAuthorized {
		tpl, _ := retrievers.RetrieveTemplate("home/unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		if tpl != nil {
			tpl.Execute(w, nil)
		}
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
			t, _ = retrievers.RetrieveTemplate("home/error")
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			if t != nil {
				t.Execute(w, nil)
			}
			return
		}
		t.Execute(w, tplMap)
	case "edit":
		itemID, err := strconv.ParseInt(pathArray[len(pathArray)-2], 10, 64)

		if err != nil {
			t, _ := retrievers.RetrieveTemplate("home/error")
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			if t != nil {
				t.Execute(w, nil)
			}
			return
		}

		item, err := handler.ItemManager.GetItem(r.Context(), itemID)

		if err != nil {
			t, _ := retrievers.RetrieveTemplate("home/error")
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			if t != nil {
				t.Execute(w, nil)
			}
			return
		}

		t, err := handler.ItemRetriever.RetrieveEditEntityTemplate()

		if err != nil {
			t, _ := retrievers.RetrieveTemplate("home/error")
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			if t != nil {
				t.Execute(w, nil)
			}
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
		handler.handleUpdateItem(w, r, userSession)
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
	id, err := strconv.ParseInt(itemID, 10, 64)
	if err != nil {
		t, _ := retrievers.RetrieveTemplate("home/error")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		if t != nil {
			t.Execute(w, nil)
		}
		return
	}

	item, err := handler.ItemManager.GetItem(r.Context(), id)
	if err != nil {
		t, _ := retrievers.RetrieveTemplate("home/error")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		if t != nil {
			t.Execute(w, nil)
		}
		return
	}

	template, err := handler.ItemRetriever.RetrieveSingleEntityTemplate()
	if err != nil {
		t, _ := retrievers.RetrieveTemplate("home/error")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		if t != nil {
			t.Execute(w, nil)
		}
		return
	}

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

func (handler ItemServiceHandler) handleUpdateItem(w http.ResponseWriter, r *http.Request, userSession *managers.UserSession) {
	item := &managers.Item{}
	err := json.NewDecoder(r.Body).Decode(item)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	previousItem, err := handler.ItemManager.GetItem(r.Context(), item.ID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if item.Status == managers.CREATED {
		item.SamaritanID = 0
	} else if userSession.UserType == managers.SAMARITAN {
		item.SamaritanID = userSession.UserID
	}

	err = handler.ItemManager.UpdateItem(r.Context(), item)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if shouldSendUpdateNotification(previousItem, item, userSession) {
		err = handler.EmailSender.DeliverEmail(r.Context(), previousItem, item, userSession)
		if err != nil {
			log.Println(err)
		}
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
	var userSession *managers.UserSession
	var userSessionError error
	cookie, cookieError := r.Cookie("NeighborsAuth")
	pathArray := strings.Split(strings.TrimPrefix(r.URL.Path, usersEndpoint), "/")

	if cookieError != nil {
		userSession, userSessionError = nil, cookieError
	} else {
		userSession, userSessionError = handler.UserSessionManager.GetUserSession(r.Context(), cookie.Value)
	}

	if r.Method == http.MethodGet && pathArray[len(pathArray)-1] != "edit" {
		if userSessionError != nil {
			log.Println(userSessionError)
		}

		return true, userSession
	}

	if cookie == nil {
		if userSessionError != nil {
			log.Println(userSessionError)
		}
		return false, userSession
	}

	if userSessionError != nil {
		log.Println(userSessionError)
		return false, userSession
	}

	if r.Method == http.MethodPost {
		return isUserAuthorized(userSession, nil, http.MethodPost), userSession
	}

	itemID, err := strconv.ParseInt(pathArray[getElementIDPathIndex(pathArray, r.Method)], 10, strconv.IntSize)
	if err != nil {
		log.Println(err)
		return false, userSession
	}

	item, err := handler.ItemManager.GetItem(r.Context(), itemID)
	if err != nil {
		log.Println(err)
		return false, userSession
	}

	return isUserAuthorized(userSession, item, r.Method), userSession
}

func isUserAuthorized(userSession *managers.UserSession, item *managers.Item, httpMethod string) bool {
	if userSession == nil {
		return false
	}

	if time.Now().After(time.Unix(userSession.LoginTime+24*7*3600, 0)) {
		return false
	}

	switch httpMethod {
	case http.MethodGet:
		return isShelterAuthorized(userSession, item) ||
			isSamaritanAuthorized(userSession, item) ||
			(userSession.UserType == managers.SAMARITAN && item.SamaritanID <= 0)
	case http.MethodPost:
		return userSession != nil && userSession.UserType == managers.SHELTER
	case http.MethodDelete:
		return isShelterAuthorized(userSession, item)
	case http.MethodPut:
		switch item.Status {
		case managers.CREATED:
			return isShelterAuthorized(userSession, item) || (userSession.UserType == managers.SAMARITAN && item.SamaritanID <= 0)
		case managers.CLAIMED:
			fallthrough
		case managers.DELIVERED:
			return isShelterAuthorized(userSession, item) || isSamaritanAuthorized(userSession, item)
		case managers.RECEIVED:
			fallthrough
		default:
			return false
		}
	default:
		return false
	}
}

func isShelterAuthorized(userSession *managers.UserSession, item *managers.Item) bool {
	return item.ShelterID == userSession.UserID && userSession.UserType == managers.SHELTER
}

func isSamaritanAuthorized(userSession *managers.UserSession, item *managers.Item) bool {
	return item.SamaritanID == userSession.UserID && userSession.UserType == managers.SAMARITAN
}

func getElementIDPathIndex(pathArray []string, method string) int {
	pathArraySize := len(pathArray)
	if method == http.MethodGet {
		return pathArraySize - 2
	}
	return pathArraySize - 1
}

func shouldSendUpdateNotification(previousItem *managers.Item, updatedItem *managers.Item, updater *managers.UserSession) bool {
	if updater.UserType == managers.SAMARITAN {
		return previousItem.Status != updatedItem.Status
	}

	return !reflect.DeepEqual(previousItem, updatedItem)
}
