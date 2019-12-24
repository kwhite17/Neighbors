package resources

import "testing"

import "net/http/httptest"

import "net/http"

import "time"

import "github.com/kwhite17/Neighbors/pkg/managers"

import "math/rand"

func TestCanAlwaysViewItems(t *testing.T) {
	ish := &ItemServiceHandler{}
	req := httptest.NewRequest(http.MethodGet, "/items/1", nil)

	isAuthorized, userSession := ish.isAuthorized(req)
	if !isAuthorized {
		t.Error("Expected users to always be authorized to view items")
	}

	if userSession != nil {
		t.Error("Expected no user session to be returned for non-edit GETS")
	}
}

func TestCanNeverAuthorizedUserWithExpiredCookieForEdits(t *testing.T) {
	ish := &ItemServiceHandler{}
	req := httptest.NewRequest(http.MethodGet, "/items/1/edit", nil)
	cookie := &http.Cookie{Expires: time.Unix(0, 0), Name: "NeighborsAuth"}
	req.AddCookie(cookie)

	isAuthorized, userSession := ish.isAuthorized(req)
	if isAuthorized {
		t.Error("User should never be authorized for edits with expired cookie")
	}

	if userSession != nil {
		t.Error("Unauthorized user should have no corresponding session")
	}
}

func TestCanNeverAuthorizedUserWithExpiredCookieForWrites(t *testing.T) {
	ish := &ItemServiceHandler{}
	req := httptest.NewRequest(http.MethodPut, "/items/1", nil)
	cookie := &http.Cookie{Expires: time.Unix(0, 0), Name: "NeighborsAuth"}
	req.AddCookie(cookie)

	isAuthorized, userSession := ish.isAuthorized(req)
	if isAuthorized {
		t.Error("User should never be authorized for edits with expired cookie")
	}

	if userSession != nil {
		t.Error("Unauthorized user should have no corresponding session")
	}
}

func TestCanLoadEditItemPageWhenAuthorizedShelter(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID}
	item := &managers.Item{ShelterID: shelterID}

	if !isUserAuthorized(userSession, item, http.MethodGet) {
		t.Error("Expected shelter to be authorized to load edit page")
	}
}

func TestCanLoadEditItemPageWhenAuthorizedSamaritan(t *testing.T) {
	samaritanID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SAMARITAN, UserID: samaritanID}
	item := &managers.Item{SamaritanID: samaritanID}

	if !isUserAuthorized(userSession, item, http.MethodGet) {
		t.Error("Expected samaritan to be authorized to load edit page")
	}
}

func TestCannotLoadEditItemPageWhenUnauthorizedSamaritan(t *testing.T) {
	samaritanID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SAMARITAN, UserID: samaritanID}
	item := &managers.Item{SamaritanID: samaritanID - 1}

	if isUserAuthorized(userSession, item, http.MethodGet) {
		t.Error("Expected samaritan to be unauthorized to load page")
	}
}

func TestCannotLoadEditItemPageWhenUnauthorizedShelter(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID}
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
	userSession := &managers.UserSession{UserType: managers.SAMARITAN}
	if isUserAuthorized(userSession, nil, http.MethodPost) {
		t.Error("Expected samaritan to be unauthorized to create items")
	}
}

func TestCanCreateItemWhenShelterUserSessionPresent(t *testing.T) {
	userSession := &managers.UserSession{UserType: managers.SHELTER}
	if !isUserAuthorized(userSession, nil, http.MethodPost) {
		t.Error("Expected shelter to be authorized to create items")
	}
}

func TestCanDeleteItemWhenAuthorizedShelter(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID}
	item := &managers.Item{ShelterID: shelterID}

	if !isUserAuthorized(userSession, item, http.MethodDelete) {
		t.Error("Expected shelter to be authorized to delete item")
	}
}

func TestCannotDeleteItemWhenUnauthorizedSamaritan(t *testing.T) {
	samaritanID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SAMARITAN, UserID: samaritanID}
	item := &managers.Item{SamaritanID: samaritanID}

	if isUserAuthorized(userSession, item, http.MethodDelete) {
		t.Error("Expected samaritan to be unauthorized deleteItem")
	}
}

func TestCannotDeleteWhenUnauthorizedShelter(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID}
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
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID}
	item := &managers.Item{ShelterID: shelterID, Status: managers.DELIVERED}

	if !isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected shelter to be authorized to update delivered item")
	}
}

func TestCanUpdateDeliveredWhenAuthorizedSamaritan(t *testing.T) {
	samaritanID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SAMARITAN, UserID: samaritanID}
	item := &managers.Item{SamaritanID: samaritanID, Status: managers.DELIVERED}

	if !isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected samaritan to be authorized to update delivered item")
	}
}

func TestCannotUpdateDelieveredItemWhenUnauthorizedSamaritan(t *testing.T) {
	samaritanID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SAMARITAN, UserID: samaritanID}
	item := &managers.Item{SamaritanID: samaritanID - 1, Status: managers.DELIVERED}

	if isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected samaritan to be unauthorized to update delivered item")
	}
}

func TestCannotUpdateDelieveredItemWhenUnauthorizedShelter(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID}
	item := &managers.Item{ShelterID: shelterID - 1, Status: managers.DELIVERED}

	if isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected samaritan to be unauthorized to update delivered item")
	}
}

func TestCanUpdateClaimedItemWhenAuthorizedShelter(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID}
	item := &managers.Item{ShelterID: shelterID, Status: managers.CLAIMED}

	if !isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected shelter to be authorized to update claimed item")
	}
}

func TestCanUpdateClaimedWhenAuthorizedSamaritan(t *testing.T) {
	samaritanID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SAMARITAN, UserID: samaritanID}
	item := &managers.Item{SamaritanID: samaritanID, Status: managers.CLAIMED}

	if !isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected samaritan to be authorized to update claimed item")
	}
}

func TestCannotUpdateClaimedItemWhenUnauthorizedSamaritan(t *testing.T) {
	samaritanID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SAMARITAN, UserID: samaritanID}
	item := &managers.Item{SamaritanID: samaritanID - 1, Status: managers.CLAIMED}

	if isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected samaritan to be unauthorized to update claimed item")
	}
}

func TestCannotUpdateClaimedItemWhenUnauthorizedShelter(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID}
	item := &managers.Item{ShelterID: shelterID - 1, Status: managers.CLAIMED}

	if isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected samaritan to be unauthorized to update claimed item")
	}
}

func TestCannotUpdateCreatedItemWhenUnauthorizedShelter(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID}
	item := &managers.Item{ShelterID: shelterID - 1, Status: managers.CREATED}

	if isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected samaritan to be unauthorized to update created item")
	}
}

func TestCanUpdateCreatedItemWhenAuthorizedShelter(t *testing.T) {
	shelterID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SHELTER, UserID: shelterID}
	item := &managers.Item{ShelterID: shelterID, Status: managers.CREATED}

	if !isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected shelter to be authorized to update created item")
	}
}

func TestCanUpdateCreatedItemWhenSamaritanSessionAndItemUnclaimed(t *testing.T) {
	samaritanID := rand.Int63()
	userSession := &managers.UserSession{UserType: managers.SAMARITAN, UserID: samaritanID}
	item := &managers.Item{Status: managers.CREATED}

	if !isUserAuthorized(userSession, item, http.MethodPut) {
		t.Error("Expected samaritan to be authorized to update claimed item")
	}
}
