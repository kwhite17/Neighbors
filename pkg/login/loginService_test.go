package login

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/kwhite17/Neighbors/pkg/test"
)

var service = LoginServiceHandler{Database: test.TestConnection}

func createServer() *httptest.Server {
	testMux := http.NewServeMux()
	testMux.Handle("/login/", service)
	return httptest.NewServer(testMux)
}

func TestAcceptUserLogin(t *testing.T) {
	ts := createServer()
	client := ts.Client()
	password, err := bcrypt.GenerateFromPassword([]byte("testPassword"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("ERROR - AcceptUserLogin - PasswordGeneration: %v\n", err)
	}
	createUserQuery := "INSERT INTO users (Username, Password, Location, Role) VALUES ('testUser', ?, 'Detroit', 'NEIGHBOR')"
	defer test.CleanUsersTable()
	defer ts.CloseClientConnections()
	defer ts.Close()
	_, err = service.Database.ExecuteWriteQuery(context.Background(), createUserQuery, []interface{}{string(password)})
	if err != nil {
		t.Fatalf("ERROR - AcceptUserLogin - DatabaseSetup: %v\n", err)
	}

	response, err := client.PostForm(ts.URL+"/login/", url.Values{"username": {"testUser"}, "password": {"testPassword"}})
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("AcceptUserLogin Failure - Expected: %d, Actual: %d\n", http.StatusOK, response.StatusCode)
		t.FailNow()
	}
	cookies := response.Cookies()
	if cookies[0].Value == "" || cookies[0].Name != "NeighborsAuth" {
		t.Errorf("AcceptUserLogin Failure - Expected Name: %s, Actual: %s\n", "NeighborsAuth", cookies[0].Name)
	}
}

func TestRejectUserLogin(t *testing.T) {
	ts := createServer()
	client := ts.Client()
	password, err := bcrypt.GenerateFromPassword([]byte("testPassword"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("ERROR - AcceptUserLogin - PasswordGeneration: %v\n", err)
	}
	createUserQuery := "INSERT INTO users (Username, Password, Location, Role) VALUES ('testUser', ?, 'Detroit', 'NEIGHBOR')"
	defer test.CleanUsersTable()
	defer ts.CloseClientConnections()
	defer ts.Close()
	_, err = service.Database.ExecuteWriteQuery(context.Background(), createUserQuery, []interface{}{string(password)})
	if err != nil {
		t.Fatalf("ERROR - AcceptUserLogin - DatabaseSetup: %v\n", err)
	}

	response, err := client.PostForm(ts.URL+"/login/", url.Values{"username": {"testUser"}, "password": {"badPassword"}})
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if response.StatusCode != http.StatusUnauthorized {
		t.Errorf("RejectUserLogin Failure - Expected: %d, Actual: %d\n", http.StatusUnauthorized, response.StatusCode)
		t.FailNow()
	}
}
