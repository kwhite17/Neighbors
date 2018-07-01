package test

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/kwhite17/Neighbors/pkg/utils"

	"github.com/kwhite17/Neighbors/pkg/database"
)

var TestConnection = database.NeighborsDatabase
var createTestUsersQuery = "INSERT INTO users (Username, Password, Location, Role) VALUES (?, ?, ?, ?)"
var createTestItemsQuery = "INSERT INTO items (Category, Size, Quantity, DropoffLocation, Requestor, OrderStatus) VALUES (?, ?, ?, ?, ?, 'REQUESTED')"
var deleteTestUsersQuery = "DELETE FROM users WHERE Username LIKE 'testUser%'"
var deleteTestItemsQuery = "DELETE FROM items WHERE Category='TESTITEM'"
var deleteTestSessionQuery = "DELETE FROM userSession WHERE SessionKey LIKE 'testKey-%'"
var testUsers = buildTestUsers()
var testItems = buildTestItems()

func PopulateUsersTable() ([]int64, error) {
	ids := make([]int64, 0)
	for _, v := range testUsers {
		output, err := TestConnection.ExecuteWriteQuery(context.Background(), createTestUsersQuery, v)
		id, err := output.LastInsertId()
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func PopulateItemsTable() ([]int64, []int64, error) {
	neighborIds, err := PopulateUsersTable()
	if err != nil {
		return nil, nil, err
	}
	ids := make([]int64, 0)
	for i := 0; i < len(neighborIds); i++ {
		for _, v := range testItems {
			output, err := TestConnection.ExecuteWriteQuery(context.Background(), createTestItemsQuery, append(v, neighborIds[i]))
			id, err := output.LastInsertId()
			if err != nil {
				return nil, nil, err
			}
			ids = append(ids, id)
		}
	}
	return ids, neighborIds, nil
}

func buildTestUsers() map[string][]interface{} {
	neighbors := make(map[string][]interface{})
	neighbors["testUserSomerville"] = []interface{}{"testUser1", "testUser", "Somerville", "SAMARITAN"}
	neighbors["testUserDetroit"] = []interface{}{"testUser2", "testUser", "Detroit", "NEIGHBOR"}
	return neighbors
}

func buildTestItems() map[string][]interface{} {
	items := make(map[string][]interface{})
	items["testItemMedium"] = []interface{}{"TESTITEM", "M", 1, "Home"}
	items["testItemLarge"] = []interface{}{"TESTITEM", "L", 1, "Home"}
	return items
}

func CleanUsersTable() {
	TestConnection.ExecuteWriteQuery(context.Background(), deleteTestUsersQuery, nil)
}

func CleanItemsTable() {
	TestConnection.ExecuteWriteQuery(context.Background(), deleteTestItemsQuery, nil)
}

func CleanUserSessionTable() {
	TestConnection.ExecuteWriteQuery(context.Background(), deleteTestSessionQuery, nil)
}

func BuildUserSession(sh utils.ServiceHandler, userID int64) (string, error) {
	return BuildUserSessionWithRole(sh, userID, "NEIGHBOR")
}

func BuildUserSessionWithRole(sh utils.ServiceHandler, userID int64, userRole string) (string, error) {
	userSessionQuery := "INSERT INTO userSession (SessionKey, UserID, LoginTime, LastSeenTime, Role) VALUES (?, ?, ?, ?, ?)"
	curTime := time.Now().UnixNano()
	testKey := "testKey-" + strconv.FormatInt(curTime, 10)
	_, err := sh.GetDatasource().ExecuteWriteQuery(context.Background(), userSessionQuery, []interface{}{testKey, userID, curTime, curTime, userRole})
	if err != nil {
		return "", fmt.Errorf("ERROR - BuildUserSession: %v", err)
	}
	return testKey, nil
}
