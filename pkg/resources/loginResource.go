package resources

import (
	"encoding/json"
	"log"
	"net/http"

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
