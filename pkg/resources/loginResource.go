package resources

import (
	"encoding/json"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/kwhite17/Neighbors/pkg/email"
	"github.com/kwhite17/Neighbors/pkg/managers"
	"github.com/kwhite17/Neighbors/pkg/retrievers"
)

const RESET_PASSWORD_LENGTH = 12

var serviceEndpoint = "/session/"
var alphaNumericRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func GenerateResetPassword(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = alphaNumericRunes[rand.Intn(len(alphaNumericRunes))]
	}
	return string(b)
}

type LoginServiceHandler struct {
	UserSessionManager managers.SessionManger
	UserManager        *managers.UserManager
	LoginRetriever     *retrievers.LoginRetriever
	EmailSender        email.EmailSender
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
		pathArray := strings.Split(strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, usersEndpoint), "/"), "/")
		path := pathArray[len(pathArray)-1]
		var t *template.Template
		var err error
		if path == "reset" {
			t, err = lsh.LoginRetriever.RetrieveEditEntityTemplate()
		} else {
			t, err = lsh.LoginRetriever.RetrieveSingleEntityTemplate()
		}

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
	case "PUT":
		resetData := make(map[string]string, 0)
		err := json.NewDecoder(r.Body).Decode(&resetData)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		emailAddress := resetData["Email"]
		unencryptedPassword := GenerateResetPassword(RESET_PASSWORD_LENGTH)
		err = lsh.UserManager.UpdatePasswordForUser(r.Context(), emailAddress, unencryptedPassword)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := lsh.UserManager.GetUserByEmail(r.Context(), emailAddress)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = lsh.EmailSender.DeliverPasswordResetEmail(r.Context(), user, unencryptedPassword)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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

	if r.Method == http.MethodPut {
		return true, nil
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
