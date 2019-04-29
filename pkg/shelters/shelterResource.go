package shelters

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/kwhite17/Neighbors/pkg/utils"
)

var serviceEndpoint = "/shelters/"

type ShelterServiceHandler struct {
	ShelterManager   *ShelterManager
	ShelterRetriever *ShelterRetriever
}

func (handler ShelterServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authRole := &utils.AuthRole{}
	// authRole, err := utils.IsAuthenticated(handler, w, r)
	// if authRole == nil && r.Method != http.MethodPost {
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
		handler.requestMethodHandler(w, r, authRole)
	}
}

func (handler ShelterServiceHandler) requestMethodHandler(w http.ResponseWriter, r *http.Request, authRole *utils.AuthRole) {
	switch r.Method {
	case http.MethodPost:
		handler.handleCreateShelter(w, r, authRole)
	case http.MethodGet:
		handler.handleGetShelter(w, r, authRole)
	case http.MethodDelete:
		handler.handleDeleteShelter(w, r, authRole)
	case http.MethodPut:
		handler.handleUpdateShelter(w, r, authRole)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (handler ShelterServiceHandler) handleCreateShelter(w http.ResponseWriter, r *http.Request, authRole *utils.AuthRole) {
	if !handler.isAuthorized(authRole, r, nil) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	contactInfo := &ContactInformation{}
	err := json.NewDecoder(r.Body).Decode(contactInfo)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	shelter := &Shelter{ContactInformation: contactInfo}
	shelterID, _ := handler.ShelterManager.WriteShelter(r.Context(), shelter)
	shelter.ID = shelterID

	json.NewEncoder(w).Encode(shelter)
}

func (handler ShelterServiceHandler) handleUpdateShelter(w http.ResponseWriter, r *http.Request, authRole *utils.AuthRole) {
	if !handler.isAuthorized(authRole, r, nil) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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

func (handler ShelterServiceHandler) handleGetShelter(w http.ResponseWriter, r *http.Request, authRole *utils.AuthRole) {
	if !handler.isAuthorized(authRole, r, nil) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if username := strings.TrimPrefix(r.URL.Path, serviceEndpoint); len(username) > 0 {
		handler.handleGetSingleShelter(w, r, username)
	} else {
		handler.handleGetAllShelters(w, r)
	}
}

func (handler ShelterServiceHandler) handleGetSingleShelter(w http.ResponseWriter, r *http.Request, username string) {
	id, _ := strconv.ParseInt(username, 10, 64)
	shelter, _ := handler.ShelterManager.GetShelter(r.Context(), id)

	template, _ := handler.ShelterRetriever.RetrieveSingleEntityTemplate()
	template.Execute(w, shelter)
}

func (handler ShelterServiceHandler) handleGetAllShelters(w http.ResponseWriter, r *http.Request) {
	shelters, _ := handler.ShelterManager.GetShelters(r.Context())

	template, _ := handler.ShelterRetriever.RetrieveAllEntitiesTemplate()
	template.Execute(w, shelters)
}

func (handler ShelterServiceHandler) handleDeleteShelter(w http.ResponseWriter, r *http.Request, authRole *utils.AuthRole) {
	username := strings.TrimPrefix(r.URL.Path, serviceEndpoint)

	_, err := handler.ShelterManager.DeleteShelter(r.Context(), username)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler ShelterServiceHandler) isAuthorized(role *utils.AuthRole, r *http.Request, data map[string]interface{}) bool {
	return true
	// if role == nil {
	// 	return false
	// }
	// switch r.Method {
	// case http.MethodGet:
	// 	return true
	// 	// pathArray := strings.Split(strings.TrimPrefix(r.URL.Path, serviceEndpoint), "/")
	// 	// switch pathArray[len(pathArray)-1] {
	// 	// case "edit":
	// 	// 	userID := pathArray[len(pathArray)-2]
	// 	// 	id, err := strconv.ParseInt(userID, 10, 64)

	// 	// 	userData, err := handler.ShelterManager.GetShelter(r.Context(), id)
	// 	// 	if err != nil {
	// 	// 		log.Println(err)
	// 	// 		return false
	// 	// 	}
	// 	// 	return userData.ID == role.ID
	// 	// default:
	// 	// 	return true
	// 	// }
	// case http.MethodPost:
	// 	return true
	// case http.MethodPut:
	// 	fallthrough
	// case http.MethodDelete:
	// 	return role.ID == data["ID"]
	// default:
	// 	return false
	// }
}
