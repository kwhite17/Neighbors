package resources

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/kwhite17/Neighbors/pkg/managers"
	"github.com/kwhite17/Neighbors/pkg/retrievers"
)

var serviceEndpoint = "/session/"

type LoginServiceHandler struct {
	UserSessionManager managers.SessionManger
	UserManager        *managers.UserManager
	LoginRetriever     *retrievers.LoginRetriever
}

func (lsh LoginServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	isAuthorized, userSession := lsh.isAuthorized(r)

	if !isAuthorized {
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	tplMap := map[string]interface{}{
		"UserSession": userSession,
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
		_, err := lsh.UserSessionManager.DeleteUserSession(r.Context(), userSession.SessionKey)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusNoContent)
	case "POST":
		loginData := make(map[string]string, 0)
		err := json.NewDecoder(r.Body).Decode(&loginData)

		shelter, err := lsh.UserManager.GetPasswordForUsername(r.Context(), loginData["Name"])
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

		sessionKey, err := lsh.UserSessionManager.WriteUserSession(r.Context(), shelter.ID, shelter.UserType)
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

func (lsh LoginServiceHandler) isAuthorized(r *http.Request) (bool, *managers.UserSession) {
	var userSession *managers.UserSession
	var userSessionError error
	cookie, cookieError := r.Cookie("NeighborsAuth")

	if cookieError != nil {
		userSession, userSessionError = nil, cookieError
	} else {
		userSession, userSessionError = lsh.UserSessionManager.GetUserSession(r.Context(), cookie.Value)
	}

	if r.Method == http.MethodGet || r.Method == http.MethodPost {
		if userSessionError != nil {
			log.Println(userSessionError)
		}
		return true, userSession
	}

	if cookie == nil || r.Method != http.MethodDelete {
		if userSessionError != nil {
			log.Println(userSessionError)
		}
		return false, nil
	}

	if userSessionError != nil || userSession == nil {
		log.Println(userSessionError)
		return false, nil
	}

	if time.Now().After(time.Unix(userSession.LoginTime+24*7*3600, 0)) {
		return false, nil
	}

	return userSession.SessionKey == cookie.Value, userSession
}
