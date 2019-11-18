package managers

import (
	"context"
	"database/sql"

	"github.com/kwhite17/Neighbors/pkg/database"
)

var createItemQuery = "INSERT INTO items (Category, Gender, Quantity, ShelterID, Size, Status) VALUES ($1, $2, $3, $4, $5, $6)"
var deleteItemQuery = "DELETE FROM items WHERE id=$1"
var getSingleItemQuery = "SELECT ID, Category, Gender, Quantity, ShelterID, Size, Status FROM items WHERE ID=$1"
var getAllItemsQuery = "SELECT ID, Category, Gender, Quantity, ShelterID, Size, Status from items"
var updateItemQuery = "UPDATE items SET Category = $1, Gender = $2, Quantity = $3, ShelterID = $4, Size = $5, Status = $6 WHERE ID = $7"
var getItemsForShelterQuery = "SELECT ID, Category, Gender, Quantity, ShelterID, Size, Status from items WHERE ShelterID = $1"

type ItemManager struct {
	Datasource database.Datasource
}

type Item struct {
	ID        int64
	Category  string
	Gender    string
	Quantity  int8
	ShelterID int64
	Size      string
	Status    string
}

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
	values := []interface{}{item.Category, item.Gender, item.Quantity, item.ShelterID, item.Size, item.Status, item.ID}
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
		var size string
		var status string
		if err := result.Scan(&id, &category, &gender, &quantity, &shelterID, &size, &status); err != nil {
			return nil, err
		}
		item := Item{ID: id, Category: category, Gender: gender, Quantity: quantity, ShelterID: shelterID, Size: size, Status: status}
		response = append(response, &item)
	}
	return response, nil
}
