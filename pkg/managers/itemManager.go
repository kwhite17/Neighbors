package managers

import (
	"context"
	"database/sql"

	"github.com/kwhite17/Neighbors/pkg/database"
)

var createItemQuery = "INSERT INTO items (Category, Gender, Quantity, ShelterID, Size, Status) VALUES (?, ?, ?, ?, ?, ?)"
var deleteItemQuery = "DELETE FROM items WHERE id=?"
var getSingleItemQuery = "SELECT ID, Category, Gender, Quantity, ShelterID, Size, Status FROM items WHERE ID=?"
var getAllItemsQuery = "SELECT ID, Category, Gender, Quantity, ShelterID, Size, Status from items"
var updateItemQuery = "UPDATE items SET Category = ?, Gender = ?, Quantity = ?, ShelterID = ?, Size = ?, Status = ? WHERE ID = ?"
var getItemsForShelterQuery = "SELECT ID, Category, Gender, Quantity, ShelterID, Size, Status from items WHERE ShelterID = ?"

type ItemManager struct {
	Datasource database.Datasource
	database.DbManager
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
	result, err := im.ReadEntity(ctx, id)
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
	result, err := im.ReadEntities(ctx)
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
	result, err := im.ReadEntitiesWithQuery(ctx, getItemsForShelterQuery, []interface{}{shelterID})
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
	result, err := im.WriteEntity(ctx, values)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func (im *ItemManager) UpdateItem(ctx context.Context, item *Item) error {
	values := []interface{}{item.Category, item.Gender, item.Quantity, item.ShelterID, item.Size, item.Status, item.ID}
	_, err := im.Datasource.ExecuteWriteQuery(ctx, updateItemQuery, values)
	return err
}

func (im *ItemManager) DeleteItem(ctx context.Context, id interface{}) (int64, error) {
	result, err := im.DeleteEntity(ctx, id)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (im *ItemManager) ReadEntity(ctx context.Context, id interface{}) (*sql.Rows, error) {
	return im.Datasource.ExecuteReadQuery(ctx, getSingleItemQuery, []interface{}{id})
}

func (im *ItemManager) ReadEntities(ctx context.Context) (*sql.Rows, error) {
	return im.ReadEntitiesWithQuery(ctx, getAllItemsQuery, nil)
}

func (im *ItemManager) ReadEntitiesWithQuery(ctx context.Context, query string, args []interface{}) (*sql.Rows, error) {
	return im.Datasource.ExecuteReadQuery(ctx, query, args)
}

func (im *ItemManager) WriteEntity(ctx context.Context, values []interface{}) (sql.Result, error) {
	result, err := im.Datasource.ExecuteWriteQuery(ctx, createItemQuery, values)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (im *ItemManager) DeleteEntity(ctx context.Context, id interface{}) (sql.Result, error) {
	result, err := im.Datasource.ExecuteWriteQuery(ctx, deleteItemQuery, []interface{}{id})
	if err != nil {
		return nil, err
	}
	return result, nil
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
