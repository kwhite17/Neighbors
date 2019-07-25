package resources

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/kwhite17/Neighbors/pkg/managers"
	"github.com/kwhite17/Neighbors/pkg/retrievers"
)

var serviceEndpoint = "/session/"

type LoginServiceHandler struct {
	ShelterSessionManager *managers.ShelterSessionManager
	ShelterManager        *managers.ShelterManager
	LoginRetriever        *retrievers.LoginRetriever
}

func (lsh LoginServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	isAuthorized, shelterSession := lsh.isAuthorized(r)

	if isAuthorized && shelterSession != nil && shelterSession.ShelterID > 0 {
		http.Redirect(w, r, "/", http.StatusFound)

		return
	}

	tplMap := map[string]interface{}{
		"ShelterSession": shelterSession,
	}

	switch r.Method {
	case "GET":
		t, err := lsh.LoginRetriever.RetrieveSingleEntityTemplate()

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		t.Execute(w, tplMap)
	case "DELETE":
		cookie, err := r.Cookie("NeighborsAuth")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		_, err = lsh.ShelterSessionManager.DeleteShelterSession(r.Context(), cookie.Value)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusNoContent)
	case "POST":
		loginData := make(map[string]string, 0)
		err := json.NewDecoder(r.Body).Decode(&loginData)

		shelter, err := lsh.ShelterManager.GetPasswordForUsername(r.Context(), loginData["Name"])
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(shelter.Password), []byte(loginData["Password"]))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		sessionKey, err := lsh.ShelterSessionManager.WriteShelterSession(r.Context(), shelter.ID, loginData["Name"])
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{Name: "NeighborsAuth", Value: sessionKey, HttpOnly: false, MaxAge: 24 * 3600 * 7, Secure: false, Path: "/"}
		shelter.Password = ""
		http.SetCookie(w, &cookie)
		json.NewEncoder(w).Encode(shelter)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (lsh LoginServiceHandler) isAuthorized(r *http.Request) (bool, *managers.ShelterSession) {
	cookie, _ := r.Cookie("NeighborsAuth")
	pathArray := strings.Split(strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, sheltersEndpoint), "/"), "/")

	switch r.Method {
	case http.MethodPost:
		if cookie == nil {
			return true, nil
		}

		shelterSession, err := lsh.ShelterSessionManager.GetShelterSession(r.Context(), cookie.Value)

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

		shelterSession, err := lsh.ShelterSessionManager.GetShelterSession(r.Context(), cookie.Value)

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
			shelterSession, err = lsh.ShelterSessionManager.GetShelterSession(r.Context(), cookie.Value)

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
