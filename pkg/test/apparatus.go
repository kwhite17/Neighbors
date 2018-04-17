package test

import (
	"context"

	"github.com/kwhite17/Neighbors/pkg/database"
)

var TestConnection = database.NeighborsDatabase
var createTestUsersQuery = "INSERT INTO users (Username, Password, Location, Role) VALUES (?, ?, ?, ?)"
var createTestItemsQuery = "INSERT INTO items (Category, Size, Quantity, DropoffLocation, Requestor) VALUES (?, ?, ?, ?, ?)"
var deleteTestUsersQuery = "DELETE FROM users WHERE Username LIKE 'testUser%'"
var deleteTestItemsQuery = "DELETE FROM items WHERE Category='TESTITEM'"
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
	TestConnection.ExecuteWriteQuery(context.Background(), deleteTestUsersQuery, []interface{}{})
}

func CleanItemsTable() {
	TestConnection.ExecuteWriteQuery(context.Background(), deleteTestItemsQuery, []interface{}{})
}
