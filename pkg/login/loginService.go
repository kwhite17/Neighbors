package login

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kwhite17/Neighbors/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

type LoginServiceHandler struct {
	Database database.Datasource
}

func (lsh LoginServiceHandler) GetDatasource() database.Datasource {
	return lsh.Database
}

func (lsh LoginServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		username := r.FormValue("username")
		hash, err := lsh.getPasswordForComparison(r.Context(), username)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(r.FormValue("password")))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}
		cookieID := username + "-" + uuid.New().String()
		cookie := http.Cookie{Name: "NeighborsAuth", Value: cookieID, HttpOnly: true, MaxAge: 24 * 3600 * 7, Secure: true}
		err = lsh.generateUserSession(r.Context(), username, cookieID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (lsh LoginServiceHandler) getPasswordForComparison(ctx context.Context, username string) (string, error) {
	result, err := lsh.Database.ExecuteReadQuery(ctx, "SELECT Password FROM users WHERE Username=?", []interface{}{username})
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

func (lsh LoginServiceHandler) generateUserSession(ctx context.Context, username string, cookieID string) error {
	result, err := lsh.Database.ExecuteReadQuery(ctx, "SELECT ID, Role FROM users WHERE Username=?", []interface{}{username})
	if err != nil {
		return fmt.Errorf("ERROR - LoginService - Database Write: %v\n", err)
	}
	for result.Next() {
		var ID int64
		var role string
		if err := result.Scan(&ID, &role); err != nil {
			return fmt.Errorf("ERROR - LoginService - Result Parse: %v\n", err)
		}
		currentTime := time.Now().Unix()
		sessionResult, err := lsh.Database.ExecuteWriteQuery(ctx,
			"INSERT INTO userSession (SessionKey, UserID, LoginTime, LastSeenTime, Role) VALUES (?, ?, ?, ?, ?)",
			[]interface{}{cookieID, ID, currentTime, currentTime, role})
		if err != nil {
			return fmt.Errorf("ERROR - LoginService - SessionCreation: %v\n", err)
		}
		rowsChanged, err := sessionResult.RowsAffected()
		if err != nil || rowsChanged < int64(1) {
			return fmt.Errorf("ERROR - LoginService - SessionCreation: %v, RowsChanged: %d\n", err, rowsChanged)
		}
		return nil
	}
	return fmt.Errorf("ERROR - LoginService - No Such User")
}
