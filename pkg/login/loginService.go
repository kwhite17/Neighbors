package login

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"

	"github.com/kwhite17/Neighbors/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

type LoginServiceHandler struct {
	Database database.Datasource
}

func (lsh LoginServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		values := r.Form
		username := values.Get("username")
		hash, err := lsh.GetPasswordForComparison(username)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(values.Get("password")))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}
		cookieID := make([]byte, sha256.BlockSize)
		sha256.New().Write(cookieID)
		cookie := http.Cookie{Name: "NeighborsAuth", Value: username + "-" + string(cookieID), HttpOnly: true, MaxAge: 24 * 3600 * 7, Secure: true}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (lsh LoginServiceHandler) GetPasswordForComparison(username string) (string, error) {
	result, err := lsh.Database.ExecuteReadQuery(nil, "SELECT Password FROM Users WHERE username=?", []interface{}{username})
	if err != nil {
		return "", fmt.Errorf("ERROR - LoginService - Database Read: %v\n", err)
	}
	for result.Next() {
		var password string
		if err := result.Scan(&password); err != nil {
			return "", fmt.Errorf("ERROR - LoginService - Result Parse: %v\n", err)
		}
		return password, nil
	}
	return "", fmt.Errorf("ERROR - LoginService - No Such User")
}
