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

func TestCanViewItemsWithNoCookie(t *testing.T) {
	ish := &ItemServiceHandler{}

	req := httptest.NewRequest(http.MethodGet, "/items/1", nil)
	isAuthorized, userSession := ish.isAuthorized(req)
	if !isAuthorized {
		t.Error("Expected users to always be authorized to view items")
	}

	if userSession != nil {
		t.Error("Expected no user session to be returned for GET TestWithoutCookie")
	}
}

func TestCanViewItemsWithCookieAndNoSession(t *testing.T) {
	ish := &ItemServiceHandler{}
	ctrl := gomock.NewController(t)
	ish.UserSessionManager = getMockSessionManager(ctrl, testKey, managers.SHELTER, 1, sql.ErrNoRows)
	defer ctrl.Finish()

	req := httptest.NewRequest(http.MethodGet, "/items/1", nil)
	req.AddCookie(&http.Cookie{Name: "NeighborsAuth", Value: testKey})
	isAuthorized, userSession := ish.isAuthorized(req)
	if !isAuthorized {
		t.Error("Expected users to always be authorized to view items")
	}

	if userSession != nil {
		t.Error("Expected no user session to be returned for GET TestWithCookieAndNoSession")
	}
}

func TestCanViewItemsWithSession(t *testing.T) {
	ish := &ItemServiceHandler{}
	ctrl := gomock.NewController(t)
	ish.UserSessionManager = getMockSessionManager(ctrl, testKey, managers.SHELTER, 1, nil)
	defer ctrl.Finish()

	req := httptest.NewRequest(http.MethodGet, "/items/1", nil)
	req.AddCookie(&http.Cookie{Name: "NeighborsAuth", Value: testKey})
	isAuthorized, userSession := ish.isAuthorized(req)
	if !isAuthorized {
		t.Error("Expected users to always be authorized to view items")
	}

	if userSession == nil {
		t.Error("Expected user session to be returned for GET TestWithCookieAndSession")
	}
}

func TestCannotCreateItemsWithNoCookie(t *testing.T) {
	ish := &ItemServiceHandler{}

	req := httptest.NewRequest(http.MethodPost, "/items/", nil)
	isAuthorized, userSession := ish.isAuthorized(req)
	if isAuthorized {
		t.Error("Expected users with no cookie to be unable to create items")
	}

	if userSession != nil {
		t.Error("Expected no user session to be returned for POST TestWithoutCookie")
	}
}

func TestCannotCreateItemsWithCookieAndNoSession(t *testing.T) {
	ish := &ItemServiceHandler{}
	ctrl := gomock.NewController(t)
	ish.UserSessionManager = getMockSessionManager(ctrl, testKey, managers.SHELTER, 1, sql.ErrNoRows)
	defer ctrl.Finish()

	req := httptest.NewRequest(http.MethodPost, "/items/", nil)
	req.AddCookie(&http.Cookie{Name: "NeighborsAuth", Value: testKey})
	isAuthorized, userSession := ish.isAuthorized(req)
	if isAuthorized {
		t.Error("Expected users with no session to be unable to create items")
	}

	if userSession != nil {
		t.Error("Expected no user session to be returned for POST TestWithCookieAndNoSession")
	}
}
func TestCannotEditItemsWithNoCookie(t *testing.T) {
	ish := &ItemServiceHandler{}

	req := httptest.NewRequest(http.MethodGet, "/items/1/edit", nil)
	isAuthorized, userSession := ish.isAuthorized(req)
	if isAuthorized {
		t.Error("Expected users with no cookie to be unable to create items")
	}

	if userSession != nil {
		t.Error("Expected no user session to be returned for GET EDIT TestWithoutCookie")
	}
}

func TestCannotEditItemsWithCookieAndNoSession(t *testing.T) {
	ish := &ItemServiceHandler{}
	ctrl := gomock.NewController(t)
	ish.UserSessionManager = getMockSessionManager(ctrl, testKey, managers.SHELTER, 1, sql.ErrNoRows)
	defer ctrl.Finish()

	req := httptest.NewRequest(http.MethodGet, "/items/1/edit", nil)
	req.AddCookie(&http.Cookie{Name: "NeighborsAuth", Value: testKey})
	isAuthorized, userSession := ish.isAuthorized(req)
	if isAuthorized {
		t.Error("Expected users with no session to be unable to create items")
	}

	if userSession != nil {
		t.Error("Expected no user session to be returned for GET EDIT TestWithCookieAndNoSession")
	}
}

func TestCanNeverAuthorizeUserWithExpiredCookie(t *testing.T) {
	userSession := &managers.UserSession{LoginTime: 0}
	if isUserAuthorized(userSession, nil, http.MethodGet) {
		t.Error("User should never be authorized for edits with expired cookie")
	}
}

func TestCanLoadEditItemPageWhenAuthorizedShelter(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID, LoginTime: time.Now().Unix()}
	item := &managers.Item{ShelterID: shelterID}

	if !isUserAuthorized(userSession, item, http.MethodGet) {
		t.Error("Expected shelter to be authorized to load edit page")
	}
}

func TestCanLoadEditItemPageWhenAuthorizedSamaritan(t *testing.T) {
	samaritanID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SAMARITAN, UserID: samaritanID, LoginTime: time.Now().Unix()}
	item := &managers.Item{SamaritanID: samaritanID}

	if !isUserAuthorized(userSession, item, http.MethodGet) {
		t.Error("Expected samaritan to be authorized to load edit page")
	}
}

func TestCannotLoadEditItemPageWhenUnauthorizedSamaritan(t *testing.T) {
	samaritanID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SAMARITAN, UserID: samaritanID, LoginTime: time.Now().Unix()}
	item := &managers.Item{SamaritanID: samaritanID - 1}

	if isUserAuthorized(userSession, item, http.MethodGet) {
		t.Error("Expected samaritan to be unauthorized to load page")
	}
}

func TestCannotLoadEditItemPageWhenUnauthorizedShelter(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID, LoginTime: time.Now().Unix()}
	item := &managers.Item{ShelterID: shelterID - 1}

	if isUserAuthorized(userSession, item, http.MethodGet) {
		t.Error("Expected samaritan to be unauthorized to load page")
	}
}

func TestCannotCreateItemWhenNoUserSessionPresent(t *testing.T) {
	if isUserAuthorized(nil, nil, http.MethodPost) {
		t.Error("Expected samaritan to be unauthorized to create items")
	}
}

func TestCannotCreateItemWhenSamaritanUserSessionPresent(t *testing.T) {
	userSession := &managers.UserSession{UserType: managers.SAMARITAN, LoginTime: time.Now().Unix()}
	if isUserAuthorized(userSession, nil, http.MethodPost) {
		t.Error("Expected samaritan to be unauthorized to create items")
	}
}

func TestCanCreateItemWhenShelterUserSessionPresent(t *testing.T) {
	userSession := &managers.UserSession{UserType: managers.SHELTER, LoginTime: time.Now().Unix()}
	if !isUserAuthorized(userSession, nil, http.MethodPost) {
		t.Error("Expected shelter to be authorized to create items")
	}
}

func TestCanDeleteItemWhenAuthorizedShelter(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID, LoginTime: time.Now().Unix()}
	item := &managers.Item{ShelterID: shelterID}

	if !isUserAuthorized(userSession, item, http.MethodDelete) {
		t.Error("Expected shelter to be authorized to delete item")
	}
}

func TestCannotDeleteItemWhenUnauthorizedSamaritan(t *testing.T) {
	samaritanID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SAMARITAN, UserID: samaritanID, LoginTime: time.Now().Unix()}
	item := &managers.Item{SamaritanID: samaritanID}

	if isUserAuthorized(userSession, item, http.MethodDelete) {
		t.Error("Expected samaritan to be unauthorized deleteItem")
	}
}

func TestCannotDeleteWhenUnauthorizedShelter(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID, LoginTime: time.Now().Unix()}
	item := &managers.Item{ShelterID: shelterID - 1}

	if isUserAuthorized(userSession, item, http.MethodDelete) {
		t.Error("Expected samaritan to be unauthorized to load page")
	}
}

func TestCannotPerformOtherRequestMethodsWithNoUserSession(t *testing.T) {
	if isUserAuthorized(nil, nil, http.MethodPatch) {
		t.Error("Expected to be unauthorized to perform non-[GET|POST|PUT|DELETE] request")
	}
}

func TestCannotPerformOtherRequestMethodsWithUserSession(t *testing.T) {
	if isUserAuthorized(&managers.UserSession{}, nil, http.MethodPatch) {
		t.Error("Expected to be unauthorized to perform non-[GET|POST|PUT|DELETE] request")
	}
}

func TestCannotUpdateItemThatHasBeenReceived(t *testing.T) {
	if isUserAuthorized(nil, &managers.Item{Status: managers.RECEIVED}, http.MethodPut) {
		t.Error("Expected to be unauthorized to perform update on item that has already been received")
	}
}

func TestCanUpdateDeliveredItemWhenAuthorizedShelter(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID, LoginTime: time.Now().Unix()}
	item := &managers.Item{ShelterID: shelterID, Status: managers.DELIVERED}

	if !isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected shelter to be authorized to update delivered item")
	}
}

func TestCanUpdateDeliveredWhenAuthorizedSamaritan(t *testing.T) {
	samaritanID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SAMARITAN, UserID: samaritanID, LoginTime: time.Now().Unix()}
	item := &managers.Item{SamaritanID: samaritanID, Status: managers.DELIVERED}

	if !isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected samaritan to be authorized to update delivered item")
	}
}

func TestCannotUpdateDelieveredItemWhenUnauthorizedSamaritan(t *testing.T) {
	samaritanID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SAMARITAN, UserID: samaritanID, LoginTime: time.Now().Unix()}
	item := &managers.Item{SamaritanID: samaritanID - 1, Status: managers.DELIVERED}

	if isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected samaritan to be unauthorized to update delivered item")
	}
}

func TestCannotUpdateDelieveredItemWhenUnauthorizedShelter(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID, LoginTime: time.Now().Unix()}
	item := &managers.Item{ShelterID: shelterID - 1, Status: managers.DELIVERED}

	if isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected samaritan to be unauthorized to update delivered item")
	}
}

func TestCanUpdateClaimedItemWhenAuthorizedShelter(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID, LoginTime: time.Now().Unix()}
	item := &managers.Item{ShelterID: shelterID, Status: managers.CLAIMED}

	if !isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected shelter to be authorized to update claimed item")
	}
}

func TestCanUpdateClaimedWhenAuthorizedSamaritan(t *testing.T) {
	samaritanID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SAMARITAN, UserID: samaritanID, LoginTime: time.Now().Unix()}
	item := &managers.Item{SamaritanID: samaritanID, Status: managers.CLAIMED}

	if !isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected samaritan to be authorized to update claimed item")
	}
}

func TestCannotUpdateClaimedItemWhenUnauthorizedSamaritan(t *testing.T) {
	samaritanID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SAMARITAN, UserID: samaritanID, LoginTime: time.Now().Unix()}
	item := &managers.Item{SamaritanID: samaritanID - 1, Status: managers.CLAIMED}

	if isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected samaritan to be unauthorized to update claimed item")
	}
}

func TestCannotUpdateClaimedItemWhenUnauthorizedShelter(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID, LoginTime: time.Now().Unix()}
	item := &managers.Item{ShelterID: shelterID - 1, Status: managers.CLAIMED}

	if isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected samaritan to be unauthorized to update claimed item")
	}
}

func TestCannotUpdateCreatedItemWhenUnauthorizedShelter(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID, LoginTime: time.Now().Unix()}
	item := &managers.Item{ShelterID: shelterID - 1, Status: managers.CREATED}

	if isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected samaritan to be unauthorized to update created item")
	}
}

func TestCanUpdateCreatedItemWhenAuthorizedShelter(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID, LoginTime: time.Now().Unix()}
	item := &managers.Item{ShelterID: shelterID, Status: managers.CREATED}

	if !isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected shelter to be authorized to update created item")
	}
}

func TestCanUpdateCreatedItemWhenSamaritanSessionAndItemUnclaimed(t *testing.T) {
	samaritanID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SAMARITAN, UserID: samaritanID, LoginTime: time.Now().Unix()}
	item := &managers.Item{Status: managers.CREATED}

	if !isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected samaritan to be authorized to update claimed item")
	}
}
