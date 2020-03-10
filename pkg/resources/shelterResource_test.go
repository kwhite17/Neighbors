package resources

import (
	"database/sql"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/kwhite17/Neighbors/pkg/managers"
)

var ush = &UserServiceHandler{}
var testKey = "testKey"

func TestCanViewUsersWithoutCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/shelters/1", nil)
	isAuthorized, userSession := ush.isAuthorized(req)
	if !isAuthorized {
		t.Error("Expected users to always be authorized to view items")
	}

	if userSession != nil {
		t.Error("Expected no user session to be returned for GET TestWithoutCookie")
	}
}
func TestCanViewUsersWithCookieAndNoSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	ush.UserSessionManager = getMockSessionManager(ctrl, testKey, managers.SHELTER, 1, sql.ErrNoRows)
	defer ctrl.Finish()

	req := httptest.NewRequest(http.MethodGet, "/shelters/1", nil)
	req.AddCookie(&http.Cookie{Name: "NeighborsAuth", Value: testKey})
	isAuthorized, userSession := ush.isAuthorized(req)
	if !isAuthorized {
		t.Error("Expected users to always be authorized to view items")
	}

	if userSession != nil {
		t.Error("Expected no user session to be returned for GET TestWithCookieAndNoSession")
	}
}
func TestCanViewUsersWithSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	ush.UserSessionManager = getMockSessionManager(ctrl, testKey, managers.SHELTER, 1, nil)
	defer ctrl.Finish()

	req := httptest.NewRequest(http.MethodGet, "/shelters/1", nil)
	req.AddCookie(&http.Cookie{Name: "NeighborsAuth", Value: testKey})
	isAuthorized, userSession := ush.isAuthorized(req)
	if !isAuthorized {
		t.Error("Expected users to always be authorized to view items")
	}

	if userSession == nil {
		t.Error("Expected user session to be returned for GET TestWithCookieAndSession")
	}
}
func TestCanCreateUserWithoutCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/shelters/", nil)
	isAuthorized, userSession := ush.isAuthorized(req)
	if !isAuthorized {
		t.Error("Expected requests with no cookie to be able to create user")
	}

	if userSession != nil {
		t.Error("Expected no user session to be returned for POST TestWithoutCookie")
	}
}
func TestCanCreateUserWithCookieAndNoSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	ush.UserSessionManager = getMockSessionManager(ctrl, testKey, managers.SHELTER, 1, sql.ErrNoRows)
	defer ctrl.Finish()

	req := httptest.NewRequest(http.MethodPost, "/shelters/", nil)
	req.AddCookie(&http.Cookie{Name: "NeighborsAuth", Value: testKey})
	isAuthorized, userSession := ush.isAuthorized(req)
	if !isAuthorized {
		t.Error("Expected requests with no corresponding user session to be able to create user")
	}

	if userSession != nil {
		t.Error("Expected no user session to be returned for POST TestWithCookieAndNoSession")
	}
}
func TestCannotCreateUserWithSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	ush.UserSessionManager = getMockSessionManager(ctrl, testKey, managers.SHELTER, 1, nil)
	defer ctrl.Finish()

	req := httptest.NewRequest(http.MethodPost, "/shelters/", nil)
	req.AddCookie(&http.Cookie{Name: "NeighborsAuth", Value: testKey})
	isAuthorized, userSession := ush.isAuthorized(req)
	if isAuthorized {
		t.Error("Expected requests with corresponding user session to be unable to create new users")
	}

	if userSession == nil {
		t.Error("Expected user session to be returned for POST TestWithCookieAndSession")
	}
}

func TestCannotAuthorizeUserWithExpiredSession(t *testing.T) {
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

func getMockSessionManager(ctrl *gomock.Controller, sessionKey string, userType managers.UserType, userID int64, errorMessage error) managers.SessionManger {
	sessionManager := NewMockSessionManger(ctrl)
	if errorMessage != nil {
		sessionManager.EXPECT().GetUserSession(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, errorMessage)
	} else {
		expectedSession := &managers.UserSession{SessionKey: sessionKey, UserType: userType, UserID: userID}
		sessionManager.EXPECT().GetUserSession(gomock.Any(), gomock.Any()).AnyTimes().Return(expectedSession, nil)
	}
	return sessionManager
}
