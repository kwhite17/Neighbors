package resources

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/kwhite17/Neighbors/pkg/managers"
	"github.com/kwhite17/Neighbors/pkg/retrievers"
)

type HomeServiceHandler struct {
	UserSessionManager *managers.UserSessionManager
}

func (hsh HomeServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	isAuthorized, shelterSession := hsh.isAuthorized(r)

	if !isAuthorized {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case "GET":
		tpl, err := retrievers.RetrieveMultiTemplate("home/layout", "home/index")

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		err = tpl.Execute(w, map[string]interface{}{
			"UserSession": shelterSession,
		})

		if err != nil {
			log.Println(err)
		}
	}

	return
}

func (hsh HomeServiceHandler) isAuthorized(r *http.Request) (bool, *managers.UserSession) {
	cookie, _ := r.Cookie("NeighborsAuth")
	pathArray := strings.Split(strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, usersEndpoint), "/"), "/")

	switch r.Method {
	case http.MethodPost:
		if cookie == nil {
			return true, nil
		}

		shelterSession, err := hsh.UserSessionManager.GetUserSession(r.Context(), cookie.Value)

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

		shelterSession, err := hsh.UserSessionManager.GetUserSession(r.Context(), cookie.Value)

		if err != nil {
			log.Println(err)
			return false, shelterSession
		}

		shelterID, err := strconv.ParseInt(pathArray[len(pathArray)-1], 10, strconv.IntSize)

		if err != nil {
			log.Println(err)
			return false, shelterSession
		}

		return shelterSession.UserID == shelterID, shelterSession
	case http.MethodGet:
		var err error
		shelterSession := &managers.UserSession{}

		if cookie != nil {
			shelterSession, err = hsh.UserSessionManager.GetUserSession(r.Context(), cookie.Value)

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

			return shelterSession != nil && shelterSession.UserID == shelterID, shelterSession
		}
		return true, shelterSession
	default:
		return false, nil
	}
}
