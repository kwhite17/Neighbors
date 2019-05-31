package login

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var serviceEndpoint = "/login/"

type LoginServiceHandler struct {
	ShelterSessionManager *ShelterSessionManager
	LoginRetriever        *LoginRetriever
}

func (lsh LoginServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		t, err := lsh.LoginRetriever.RetrieveSingleEntityTemplate()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		t.Execute(w, nil)
	case "DELETE":
		pathArray := strings.Split(strings.TrimPrefix(r.URL.Path, serviceEndpoint), "/")
		shelterID := pathArray[len(pathArray)-1]

		lsh.ShelterSessionManager.DeleteShelterSession(r.Context(), shelterID)
	case "POST":
		loginData := make(map[string]string, 0)
		err := json.NewDecoder(r.Body).Decode(&loginData)

		shelterSession, err := lsh.ShelterSessionManager.GetPasswordForUsername(r.Context(), loginData["Name"])
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(shelterSession.Password), []byte(loginData["Password"]))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// currentTime := time.Now().Unix()
		// err = lsh.ShelterSessionManager.UpdateShelterSession(r.Context(), shelterSession.ShelterID, currentTime, currentTime)
		// if err != nil {
		// 	log.Println(err)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }

		cookie := http.Cookie{Name: "NeighborsAuth", Value: shelterSession.SessionKey, HttpOnly: false, MaxAge: 24 * 3600 * 7, Secure: false}
		shelterSession.Password = ""
		http.SetCookie(w, &cookie)
		json.NewEncoder(w).Encode(shelterSession)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
