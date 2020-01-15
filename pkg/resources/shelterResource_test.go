package resources

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kwhite17/Neighbors/pkg/managers"
)

var ush = &UserServiceHandler{}

func TestCanAlwaysViewUsers(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/shelters/1", nil)

	isAuthorized, userSession := ush.isAuthorized(req)
	if !isAuthorized {
		t.Error("Expected users to always be authorized to view items")
	}

	if userSession != nil {
		t.Error("Expected no user session to be returned for non-edit GETS")
	}
}

func TestCannotAuthorizeUserWithExpiredCookie(t *testing.T) {
	userSession := &managers.UserSession{LoginTime: 0}
	if ush.isUserAuthorized(userSession, -1, http.MethodGet) {
		t.Error("User should never be authorized for edits with expired cookie")
	}
}

func TestCanUpdateUserWhenAuthorizedUser(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID, LoginTime: time.Now().Unix()}

	if !ush.isUserAuthorized(userSession, shelterID, http.MethodPut) {
		t.Error("Expected shelter to be authorized to load edit page")
	}
}

func TestCanLoadEditUserPageWhenAuthorizedUser(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID, LoginTime: time.Now().Unix()}

	if !ush.isUserAuthorized(userSession, shelterID, http.MethodGet) {
		t.Error("Expected shelter to be authorized to load edit page")
	}
}

func TestCanDeleteUserWhenAuthorizedUser(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID, LoginTime: time.Now().Unix()}

	if !ush.isUserAuthorized(userSession, shelterID, http.MethodDelete) {
		t.Error("Expected shelter to be authorized to load edit page")
	}
}

func TestCannotCreateUserWhenUserSessionPresent(t *testing.T) {
	userSession := &managers.UserSession{UserType: managers.SAMARITAN, LoginTime: time.Now().Unix()}
	if ush.isUserAuthorized(userSession, -1, http.MethodPost) {
		t.Error("Expected samaritan to be unauthorized to create items")
	}
}

func TestCannotUpdateUserWhenUnauthorizedUser(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID, LoginTime: time.Now().Unix()}

	if ush.isUserAuthorized(userSession, shelterID-1, http.MethodPut) {
		t.Error("Expected shelter to be unauthorized to edit user")
	}
}

func TestCannotLoadEditUserPageWhenUnauthorizedUser(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID, LoginTime: time.Now().Unix()}

	if ush.isUserAuthorized(userSession, shelterID-1, http.MethodGet) {
		t.Error("Expected shelter to be unauthorized to load edit page")
	}
}

func TestCannotDeleteUserWhenUnauthorizedUser(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID, LoginTime: time.Now().Unix()}

	if ush.isUserAuthorized(userSession, shelterID-1, http.MethodDelete) {
		t.Error("Expected shelter to be unauthorized to delete user")
	}
}
