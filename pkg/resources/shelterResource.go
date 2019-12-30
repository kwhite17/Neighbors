package resources

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/kwhite17/Neighbors/pkg/managers"
	"github.com/kwhite17/Neighbors/pkg/retrievers"
)

var usersEndpoint = "/shelters/"

type UserServiceHandler struct {
	UserManager        *managers.UserManager
	ItemManager        *managers.ItemManager
	UserSessionManager *managers.UserSessionManager
	UserRetriever      *retrievers.ShelterRetriever
}

func (handler UserServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	isAuthorized, userSession := handler.isAuthorized(r)
	if !isAuthorized {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tplMap := map[string]interface{}{
		"UserSession": userSession,
	}

	pathArray := strings.Split(strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, usersEndpoint), "/"), "/")
	switch pathArray[len(pathArray)-1] {
	case "new":
		t, err := handler.UserRetriever.RetrieveCreateEntityTemplate()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = t.Execute(w, tplMap)

		if err != nil {
			log.Println(err)
		}
	case "edit":
		userID, err := strconv.ParseInt(pathArray[len(pathArray)-2], 10, 64)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := handler.UserManager.GetUser(r.Context(), userID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		t, err := handler.UserRetriever.RetrieveEditEntityTemplate()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		tplMap["User"] = user
		t.Execute(w, tplMap)
	default:
		handler.requestMethodHandler(w, r, userSession)
	}
}

func (handler UserServiceHandler) requestMethodHandler(w http.ResponseWriter, r *http.Request, userSession *managers.UserSession) {
	switch r.Method {
	case http.MethodPost:
		handler.handleCreateUser(w, r)
	case http.MethodGet:
		handler.handleGetUser(w, r, userSession)
	case http.MethodDelete:
		handler.handleDeleteUser(w, r)
	case http.MethodPut:
		handler.handleUpdateUser(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (handler UserServiceHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	createData := make(map[string]interface{}, 0)
	err := json.NewDecoder(r.Body).Decode(&createData)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := &managers.User{ContactInformation: handler.buildContactInformation(createData)}
	userID, err := handler.UserManager.WriteUser(r.Context(), user, createData["Password"].(string))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.ID = userID
	userType, err := strconv.Atoi(reflect.ValueOf(createData["UserType"]).String())
	cookieID, err := handler.UserSessionManager.WriteUserSession(r.Context(), userID, managers.UserType(userType))
	if err != nil {
		handler.UserManager.DeleteUser(r.Context(), userID)
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{Name: "NeighborsAuth", Value: cookieID, HttpOnly: false, MaxAge: 24 * 3600 * 7, Secure: false, Path: "/"}
	http.SetCookie(w, &cookie)
	json.NewEncoder(w).Encode(user)
}

func (handler UserServiceHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	user := &managers.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = handler.UserManager.UpdateUser(r.Context(), user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler UserServiceHandler) handleGetUser(w http.ResponseWriter, r *http.Request, userSession *managers.UserSession) {
	if user := strings.TrimPrefix(r.URL.Path, usersEndpoint); len(user) > 0 {
		handler.handleGetSingleUser(w, r, user, userSession)
	} else {
		handler.handleGetAllUsers(w, r, userSession)
	}
}

func (handler UserServiceHandler) handleGetSingleUser(w http.ResponseWriter, r *http.Request, userID string, userSession *managers.UserSession) {
	id, _ := strconv.ParseInt(userID, 10, 64)
	user, _ := handler.UserManager.GetUser(r.Context(), id)

	items, _ := handler.ItemManager.GetItemsForShelter(r.Context(), id)
	template, _ := handler.UserRetriever.RetrieveSingleEntityTemplate()

	responseObject := make(map[string]interface{}, 0)
	responseObject["User"] = user
	responseObject["Items"] = items
	responseObject["UserSession"] = userSession
	template.Execute(w, responseObject)
}

func (handler UserServiceHandler) handleGetAllUsers(w http.ResponseWriter, r *http.Request, userSession *managers.UserSession) {
	users, _ := handler.UserManager.GetUsers(r.Context())

	template, _ := handler.UserRetriever.RetrieveAllEntitiesTemplate()
	responseObject := make(map[string]interface{}, 0)
	responseObject["Users"] = users
	responseObject["UserSession"] = userSession
	template.Execute(w, responseObject)
}

func (handler UserServiceHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimPrefix(r.URL.Path, usersEndpoint)

	_, err := handler.UserManager.DeleteUser(r.Context(), userID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler UserServiceHandler) buildContactInformation(createData map[string]interface{}) *managers.ContactInformation {
	return &managers.ContactInformation{
		City:       createData["City"].(string),
		Country:    createData["Country"].(string),
		Name:       createData["Name"].(string),
		PostalCode: createData["PostalCode"].(string),
		State:      createData["State"].(string),
		Street:     createData["Street"].(string),
	}
}

func (handler UserServiceHandler) isAuthorized(r *http.Request) (bool, *managers.UserSession) {
	cookie, _ := r.Cookie("NeighborsAuth")
	pathArray := strings.Split(strings.TrimPrefix(r.URL.Path, usersEndpoint), "/")
	switch r.Method {
	case http.MethodPost:
		if cookie == nil {
			return true, nil
		}
		userSession, err := handler.UserSessionManager.GetUserSession(r.Context(), cookie.Value)
		if err != nil {
			log.Println(err)
			return err == sql.ErrNoRows, userSession
		}

		return userSession == nil, userSession
	case http.MethodPut:
		fallthrough
	case http.MethodDelete:
		if cookie == nil {
			return false, nil
		}
		userSession, err := handler.UserSessionManager.GetUserSession(r.Context(), cookie.Value)
		if err != nil {
			log.Println(err)
			return false, userSession
		}

		userID, err := strconv.ParseInt(pathArray[len(pathArray)-1], 10, strconv.IntSize)
		if err != nil {
			log.Println(err)
			return false, userSession
		}

		return userSession.UserID == userID, userSession
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
			userID, err := strconv.ParseInt(pathArray[len(pathArray)-2], 10, strconv.IntSize)
			if err != nil {
				log.Println(err)
				return false, userSession
			}

			return userSession != nil && userSession.UserID == userID, userSession
		}
		return true, userSession
	default:
		return false, nil
	}
}
