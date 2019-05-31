package shelters

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/kwhite17/Neighbors/pkg/login"
)

var serviceEndpoint = "/shelters/"

type ShelterServiceHandler struct {
	ShelterManager        *ShelterManager
	ShelterSessionManager *login.ShelterSessionManager
	ShelterRetriever      *ShelterRetriever
}

func (handler ShelterServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("NeighborsAuth")
	if !handler.isAuthorized(r, cookie) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	pathArray := strings.Split(strings.TrimPrefix(r.URL.Path, serviceEndpoint), "/")
	switch pathArray[len(pathArray)-1] {
	case "new":
		t, err := handler.ShelterRetriever.RetrieveCreateEntityTemplate()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
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
		t.Execute(w, shelter)
	default:
		handler.requestMethodHandler(w, r)
	}
}

func (handler ShelterServiceHandler) requestMethodHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handler.handleCreateShelter(w, r)
	case http.MethodGet:
		handler.handleGetShelter(w, r)
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

	shelter := &Shelter{ContactInformation: handler.buildContactInformation(createData)}
	shelterID, err := handler.ShelterManager.WriteShelter(r.Context(), shelter)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	shelter.ID = shelterID
	cookieID, err := handler.ShelterSessionManager.WriteShelterSession(r.Context(), shelterID, shelter.Name, createData["Password"].(string))
	if err != nil {
		handler.ShelterManager.DeleteShelter(r.Context(), shelterID)
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{Name: "NeighborsAuth", Value: cookieID, HttpOnly: false, MaxAge: 24 * 3600 * 7, Secure: false}
	http.SetCookie(w, &cookie)
	json.NewEncoder(w).Encode(shelter)
}

func (handler ShelterServiceHandler) handleUpdateShelter(w http.ResponseWriter, r *http.Request) {
	shelter := &Shelter{}
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

func (handler ShelterServiceHandler) handleGetShelter(w http.ResponseWriter, r *http.Request) {
	if shelter := strings.TrimPrefix(r.URL.Path, serviceEndpoint); len(shelter) > 0 {
		handler.handleGetSingleShelter(w, r, shelter)
	} else {
		handler.handleGetAllShelters(w, r)
	}
}

func (handler ShelterServiceHandler) handleGetSingleShelter(w http.ResponseWriter, r *http.Request, shelterID string) {
	id, _ := strconv.ParseInt(shelterID, 10, 64)
	shelter, _ := handler.ShelterManager.GetShelter(r.Context(), id)

	template, _ := handler.ShelterRetriever.RetrieveSingleEntityTemplate()
	template.Execute(w, shelter)
}

func (handler ShelterServiceHandler) handleGetAllShelters(w http.ResponseWriter, r *http.Request) {
	shelters, _ := handler.ShelterManager.GetShelters(r.Context())

	template, _ := handler.ShelterRetriever.RetrieveAllEntitiesTemplate()
	template.Execute(w, shelters)
}

func (handler ShelterServiceHandler) handleDeleteShelter(w http.ResponseWriter, r *http.Request) {
	shelterID := strings.TrimPrefix(r.URL.Path, serviceEndpoint)

	_, err := handler.ShelterManager.DeleteShelter(r.Context(), shelterID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler ShelterServiceHandler) buildContactInformation(createData map[string]interface{}) *ContactInformation {
	return &ContactInformation{
		City:       createData["City"].(string),
		Country:    createData["Country"].(string),
		Name:       createData["Name"].(string),
		PostalCode: createData["PostalCode"].(string),
		State:      createData["State"].(string),
		Street:     createData["Street"].(string),
	}
}

func (handler ShelterServiceHandler) isAuthorized(r *http.Request, cookie *http.Cookie) bool {
	pathArray := strings.Split(strings.TrimPrefix(r.URL.Path, serviceEndpoint), "/")
	switch r.Method {
	case http.MethodPost:
		return cookie == nil
	case http.MethodPut:
		fallthrough
	case http.MethodDelete:
		shelterSession, err := handler.ShelterSessionManager.GetShelterSession(r.Context(), cookie.Value)
		if err != nil {
			log.Println(err)
			return false
		}

		shelterID, err := strconv.ParseInt(pathArray[len(pathArray)-1], 10, strconv.IntSize)
		if err != nil {
			log.Println(err)
			return false
		}

		return shelterSession.ShelterID == shelterID
	case http.MethodGet:
		if pathArray[len(pathArray)-1] == "edit" {
			shelterSession, err := handler.ShelterSessionManager.GetShelterSession(r.Context(), cookie.Value)
			if err != nil {
				log.Println(err)
				return false
			}

			shelterID, err := strconv.ParseInt(pathArray[len(pathArray)-2], 10, strconv.IntSize)
			if err != nil {
				log.Println(err)
				return false
			}

			return shelterSession.ShelterID == shelterID
		}
		return true
	default:
		return false
	}
}
