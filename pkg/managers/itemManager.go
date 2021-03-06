package managers

import (
	"context"
	"database/sql"
	"reflect"

	"github.com/kwhite17/Neighbors/pkg/database"
)

var createItemQuery = "INSERT INTO items (Category, Gender, Quantity, ShelterID, Size, Status) VALUES ($1, $2, $3, $4, $5, $6)"
var deleteItemQuery = "DELETE FROM items WHERE id=$1"
var getSingleItemQuery = "SELECT ID, Category, Gender, Quantity, ShelterID, SamaritanID, Size, Status FROM items WHERE ID=$1"
var getAllItemsQuery = "SELECT ID, Category, Gender, Quantity, ShelterID, SamaritanID, Size, Status from items"
var updateItemQuery = "UPDATE items SET Category = $1, Gender = $2, Quantity = $3, ShelterID = $4, SamaritanID = $5, Size = $6, Status = $7 WHERE ID = $8"
var getItemsForShelterQuery = "SELECT ID, Category, Gender, Quantity, ShelterID, SamaritanID, Size, Status from items WHERE ShelterID = $1"

type ItemManager struct {
	Datasource database.Datasource
}

type Item struct {
	ID          int64
	Category    string
	Gender      string
	Quantity    int8
	ShelterID   int64
	SamaritanID int64
	Size        string
	Status      ItemStatus
}

type ItemStatus int

const (
	CREATED   ItemStatus = 1
	CLAIMED   ItemStatus = 2
	DELIVERED ItemStatus = 3
	RECEIVED  ItemStatus = 4
)

func (im *ItemManager) GetItem(ctx context.Context, id interface{}) (*Item, error) {
	result, err := im.Datasource.ExecuteBatchReadQuery(ctx, getSingleItemQuery, []interface{}{id})
	if err != nil {
		return nil, err
	}
	item, err := im.buildItems(result)
	if err != nil {
		return nil, err
	}
	return item[0], nil
}

func (im *ItemManager) GetItems(ctx context.Context) ([]*Item, error) {
	result, err := im.Datasource.ExecuteBatchReadQuery(ctx, getAllItemsQuery, nil)
	if err != nil {
		return nil, err
	}
	items, err := im.buildItems(result)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (im *ItemManager) GetItemsForShelter(ctx context.Context, shelterID int64) ([]*Item, error) {
	result, err := im.Datasource.ExecuteBatchReadQuery(ctx, getItemsForShelterQuery, []interface{}{shelterID})
	if err != nil {
		return nil, err
	}
	items, err := im.buildItems(result)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (im *ItemManager) WriteItem(ctx context.Context, item *Item) (int64, error) {
	values := []interface{}{item.Category, item.Gender, item.Quantity, item.ShelterID, item.Size, item.Status}
	result, err := im.Datasource.ExecuteWriteQuery(ctx, createItemQuery, values, true)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func (im *ItemManager) UpdateItem(ctx context.Context, item *Item) error {
	values := []interface{}{item.Category, item.Gender, item.Quantity, item.ShelterID, item.SamaritanID, item.Size, item.Status, item.ID}
	_, err := im.Datasource.ExecuteWriteQuery(ctx, updateItemQuery, values, true)
	return err
}

func (im *ItemManager) DeleteItem(ctx context.Context, id interface{}) (int64, error) {
	result, err := im.Datasource.ExecuteWriteQuery(ctx, deleteItemQuery, []interface{}{id}, true)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (im *ItemManager) buildItems(result *sql.Rows) ([]*Item, error) {
	response := make([]*Item, 0)
	for result.Next() {
		var id int64
		var category string
		var gender string
		var quantity int8
		var shelterID int64
		var samaritan interface{}
		var size string
		var status ItemStatus
		if err := result.Scan(&id, &category, &gender, &quantity, &shelterID, &samaritan, &size, &status); err != nil {
			return nil, err
		}
		item := Item{ID: id, Category: category, Gender: gender, Quantity: quantity, ShelterID: shelterID, Size: size, Status: status}
		if samaritan != nil {
			item.SamaritanID = reflect.ValueOf(samaritan).Int()
		}
		response = append(response, &item)
	}
	return response, nil
}
