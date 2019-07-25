package resources

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/kwhite17/Neighbors/pkg/managers"
	"github.com/kwhite17/Neighbors/pkg/retrievers"
)

var sheltersEndpoint = "/shelters/"

type ShelterServiceHandler struct {
	ShelterManager        *managers.ShelterManager
	ItemManager           *managers.ItemManager
	ShelterSessionManager *managers.ShelterSessionManager
	ShelterRetriever      *retrievers.ShelterRetriever
}

func (handler ShelterServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	isAuthorized, shelterSession := handler.isAuthorized(r)
	if !isAuthorized {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tplMap := map[string]interface{}{
		"ShelterSession": shelterSession,
	}

	pathArray := strings.Split(strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, sheltersEndpoint), "/"), "/")
	switch pathArray[len(pathArray)-1] {
	case "new":
		t, err := handler.ShelterRetriever.RetrieveCreateEntityTemplate()
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
		shelterID, err := strconv.ParseInt(pathArray[len(pathArray)-2], 10, 64)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		shelter, err := handler.ShelterManager.GetShelter(r.Context(), shelterID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		t, err := handler.ShelterRetriever.RetrieveEditEntityTemplate()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		tplMap["Shelter"] = shelter
		t.Execute(w, tplMap)
	default:
		handler.requestMethodHandler(w, r, shelterSession)
	}
}

func (handler ShelterServiceHandler) requestMethodHandler(w http.ResponseWriter, r *http.Request, shelterSession *managers.ShelterSession) {
	switch r.Method {
	case http.MethodPost:
		handler.handleCreateShelter(w, r)
	case http.MethodGet:
		handler.handleGetShelter(w, r, shelterSession)
	case http.MethodDelete:
		handler.handleDeleteShelter(w, r)
	case http.MethodPut:
		handler.handleUpdateShelter(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (handler ShelterServiceHandler) handleCreateShelter(w http.ResponseWriter, r *http.Request) {
	createData := make(map[string]interface{}, 0)
	err := json.NewDecoder(r.Body).Decode(&createData)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	shelter := &managers.Shelter{ContactInformation: handler.buildContactInformation(createData)}
	shelterID, err := handler.ShelterManager.WriteShelter(r.Context(), shelter, createData["Password"].(string))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	shelter.ID = shelterID
	cookieID, err := handler.ShelterSessionManager.WriteShelterSession(r.Context(), shelterID, shelter.Name)
	if err != nil {
		handler.ShelterManager.DeleteShelter(r.Context(), shelterID)
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{Name: "NeighborsAuth", Value: cookieID, HttpOnly: false, MaxAge: 24 * 3600 * 7, Secure: false, Path: "/"}
	http.SetCookie(w, &cookie)
	json.NewEncoder(w).Encode(shelter)
}

func (handler ShelterServiceHandler) handleUpdateShelter(w http.ResponseWriter, r *http.Request) {
	shelter := &managers.Shelter{}
	err := json.NewDecoder(r.Body).Decode(shelter)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = handler.ShelterManager.UpdateShelter(r.Context(), shelter)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler ShelterServiceHandler) handleGetShelter(w http.ResponseWriter, r *http.Request, shelterSession *managers.ShelterSession) {
	if shelter := strings.TrimPrefix(r.URL.Path, sheltersEndpoint); len(shelter) > 0 {
		handler.handleGetSingleShelter(w, r, shelter, shelterSession)
	} else {
		handler.handleGetAllShelters(w, r, shelterSession)
	}
}

func (handler ShelterServiceHandler) handleGetSingleShelter(w http.ResponseWriter, r *http.Request, shelterID string, shelterSession *managers.ShelterSession) {
	id, _ := strconv.ParseInt(shelterID, 10, 64)
	shelter, _ := handler.ShelterManager.GetShelter(r.Context(), id)

	items, _ := handler.ItemManager.GetItemsForShelter(r.Context(), id)
	template, _ := handler.ShelterRetriever.RetrieveSingleEntityTemplate()

	responseObject := make(map[string]interface{}, 0)
	responseObject["Shelter"] = shelter
	responseObject["Items"] = items
	responseObject["ShelterSession"] = shelterSession
	template.Execute(w, responseObject)
}

func (handler ShelterServiceHandler) handleGetAllShelters(w http.ResponseWriter, r *http.Request, shelterSession *managers.ShelterSession) {
	shelters, _ := handler.ShelterManager.GetShelters(r.Context())

	template, _ := handler.ShelterRetriever.RetrieveAllEntitiesTemplate()
	responseObject := make(map[string]interface{}, 0)
	responseObject["Shelters"] = shelters
	responseObject["ShelterSession"] = shelterSession
	template.Execute(w, responseObject)
}

func (handler ShelterServiceHandler) handleDeleteShelter(w http.ResponseWriter, r *http.Request) {
	shelterID := strings.TrimPrefix(r.URL.Path, sheltersEndpoint)

	_, err := handler.ShelterManager.DeleteShelter(r.Context(), shelterID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler ShelterServiceHandler) buildContactInformation(createData map[string]interface{}) *managers.ContactInformation {
	return &managers.ContactInformation{
		City:       createData["City"].(string),
		Country:    createData["Country"].(string),
		Name:       createData["Name"].(string),
		PostalCode: createData["PostalCode"].(string),
		State:      createData["State"].(string),
		Street:     createData["Street"].(string),
	}
}

func (handler ShelterServiceHandler) isAuthorized(r *http.Request) (bool, *managers.ShelterSession) {
	cookie, _ := r.Cookie("NeighborsAuth")
	pathArray := strings.Split(strings.TrimPrefix(r.URL.Path, sheltersEndpoint), "/")
	switch r.Method {
	case http.MethodPost:
		if cookie == nil {
			return true, nil
		}
		shelterSession, err := handler.ShelterSessionManager.GetShelterSession(r.Context(), cookie.Value)
		if err != nil {
			log.Println(err)
			return err == sql.ErrNoRows, shelterSession
		}

		return shelterSession == nil, shelterSession
	case http.MethodPut:
		fallthrough
	case http.MethodDelete:
		if cookie == nil {
			return false, nil
		}
		shelterSession, err := handler.ShelterSessionManager.GetShelterSession(r.Context(), cookie.Value)
		if err != nil {
			log.Println(err)
			return false, shelterSession
		}

		shelterID, err := strconv.ParseInt(pathArray[len(pathArray)-1], 10, strconv.IntSize)
		if err != nil {
			log.Println(err)
			return false, shelterSession
		}

		return shelterSession.ShelterID == shelterID, shelterSession
	case http.MethodGet:
		var err error
		shelterSession := &managers.ShelterSession{}
		if cookie != nil {
			shelterSession, err = handler.ShelterSessionManager.GetShelterSession(r.Context(), cookie.Value)
			if err != nil && err != sql.ErrNoRows {
				log.Println(err)
				return false, shelterSession
			}
		}

		if pathArray[len(pathArray)-1] == "edit" {
			shelterID, err := strconv.ParseInt(pathArray[len(pathArray)-2], 10, strconv.IntSize)
			if err != nil {
				log.Println(err)
				return false, shelterSession
			}

			return shelterSession != nil && shelterSession.ShelterID == shelterID, shelterSession
		}
		return true, shelterSession
	default:
		return false, nil
	}
}
