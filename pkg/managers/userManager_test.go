package managers

import (
	"context"
	"database/sql"
	"reflect"
	"strconv"
	"testing"

	"github.com/kwhite17/Neighbors/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

var testCity = "testCity"
var testEmail = "test@test.com"
var testName = "testName"
var testPostalCode = "testPostalCode"
var testState = "testState"
var testStreet = "testStreet"

var dbToClose *sql.DB

func initUserManager() *UserManager {
	dbToClose = database.InitDatabase(database.SQLITE3)
	return &UserManager{Datasource: database.StandardDatasource{Database: dbToClose}}
}

func TestCanReadItsOwnUserWrite(t *testing.T) {
	manager := initUserManager()
	defer cleanDatabase()
	testUser := generateUser()

	id, err := manager.WriteUser(context.Background(), testUser, "password")
	if err != nil {
		t.Error(err)
	}
	testUser.ID = id

	actualUser, err := manager.GetUser(context.Background(), id)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(testUser, actualUser) {
		t.Errorf("Expected %v to equal %v", actualUser, testUser)
	}
}

func TestItCanDeleteUser(t *testing.T) {
	manager := initUserManager()
	defer cleanDatabase()
	testUser := generateUser()

	id, err := manager.WriteUser(context.Background(), testUser, "password")
	if err != nil {
		t.Error(err)
	}

	rowsDeleted, err := manager.DeleteUser(context.Background(), strconv.FormatInt(id, 10))
	if err != nil {
		t.Error(err)
	}

	if rowsDeleted != 1 {
		t.Error("Expected row to be deleted")
	}
}

func TestItCanGetAllUsers(t *testing.T) {
	manager := initUserManager()
	defer cleanDatabase()
	testUsers := make([]*User, 0)
	for i := 0; i < 5; i++ {
		testUser := generateUser()

		id, err := manager.WriteUser(context.Background(), testUser, "password")
		if err != nil {
			t.Error(err)
		}
		testUser.ID = id
		testUsers = append(testUsers, testUser)
	}

	allUsers, err := manager.GetUsers(context.Background())
	if err != nil {
		t.Error(err)
	}

	for _, user := range allUsers {
		if !containsUser(user, testUsers) {
			t.Errorf("Expected %v to be in %v \n", user, testUsers)
		}
	}
}

func TestCanReadItsOwnUserUpdate(t *testing.T) {
	manager := initUserManager()
	defer cleanDatabase()
	testUser := generateUser()

	id, err := manager.WriteUser(context.Background(), testUser, "password")
	if err != nil {
		t.Error(err)
	}
	testUser.ID = id

	createdUser, err := manager.GetUser(context.Background(), id)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(testUser, createdUser) {
		t.Errorf("Expected %v to equal %v", createdUser, testUser)
	}

	createdUser.PostalCode = "02139"
	err = manager.UpdateUser(context.Background(), createdUser)
	if err != nil {
		t.Error(err)
	}

	finalUser, err := manager.GetUser(context.Background(), id)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(finalUser, createdUser) {
		t.Errorf("Expected %v to equal %v", finalUser, createdUser)
	}
}

func TestItGetsPasswordForUsername(t *testing.T) {
	manager := initUserManager()
	defer cleanDatabase()
	testUser := generateUser()
	unhashedPassword := "password"

	id, err := manager.WriteUser(context.Background(), testUser, unhashedPassword)
	if err != nil {
		t.Error(err)
	}
	testUser.ID = id

	createdUser, err := manager.GetPasswordForUsername(context.Background(), testUser.Name)
	if err != nil {
		t.Error(err)
	}

	if bcrypt.CompareHashAndPassword([]byte(createdUser.Password), []byte(unhashedPassword)) != nil {
		t.Errorf("Expected %v to equal %v", []byte(createdUser.Password), []byte(unhashedPassword))
	}
}

func generateUser() *User {
	contactInfo := &ContactInformation{
		City:       testCity,
		Email:      testEmail,
		Name:       testName,
		PostalCode: testPostalCode,
		State:      testState,
		Street:     testStreet,
	}

	return &User{ContactInformation: contactInfo, UserType: SHELTER}
}

func containsUser(candidateUser *User, expectedUsers []*User) bool {
	for _, user := range expectedUsers {
		if reflect.DeepEqual(candidateUser, user) {
			return true
		}
	}
	return false
}
