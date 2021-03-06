package managers

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/kwhite17/Neighbors/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

var testCity = "testCity"
var testEmail = "test%v@test.com"
var testName = "testName"
var testPostalCode = "testPostalCode"
var testState = "testState"
var testStreet = "testStreet"

var dbToClose *sql.DB

func initUserManager() *UserManager {
	dbToClose = database.InitDatabase(database.SQLITE3)
	return &UserManager{Datasource: database.StandardDatasource{Database: dbToClose}}
}

func TestCanReadItsOwnShelterWrite(t *testing.T) {
	manager := initUserManager()
	defer cleanDatabase()
	testUser := generateUser(0)
	testUser.UserType = SHELTER

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

func TestCanRetrieveUserByEmail(t *testing.T) {
	manager := initUserManager()
	defer cleanDatabase()
	testUser := generateUser(0)
	testUser.UserType = SHELTER

	id, err := manager.WriteUser(context.Background(), testUser, "password")
	if err != nil {
		t.Error(err)
	}
	testUser.ID = id

	actualUser, err := manager.GetUserByEmail(context.Background(), testUser.Email)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(testUser, actualUser) {
		t.Errorf("Expected %v to equal %v", actualUser, testUser)
	}
}
func TestCanReadItsOwnSamaritanWrite(t *testing.T) {
	manager := initUserManager()
	defer cleanDatabase()
	testUser := generateUser(0)
	testUser.UserType = SAMARITAN

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
	testUser := generateUser(0)

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
		testUser := generateUser(i)

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
	testUser := generateUser(0)

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

func TestCanReadItsOwnUserPasswordUpdate(t *testing.T) {
	manager := initUserManager()
	defer cleanDatabase()
	testUser := generateUser(0)

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

	updatedPassword := "newPassword"
	err = manager.UpdatePasswordForUser(context.Background(), createdUser.Email, updatedPassword)
	if err != nil {
		t.Error(err)
	}

	updatedUser, err := manager.GetPasswordForUsername(context.Background(), createdUser.Name)
	if err != nil {
		t.Error(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte(updatedPassword))
	if err != nil {
		t.Error(err)
	}
}

func TestItGetsPasswordForUsername(t *testing.T) {
	manager := initUserManager()
	defer cleanDatabase()
	testUser := generateUser(0)
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

func generateUser(id int) *User {
	contactInfo := &ContactInformation{
		City:       testCity,
		Email:      fmt.Sprintf(testEmail, id),
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
