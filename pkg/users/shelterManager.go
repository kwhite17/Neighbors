package users

import (
	"context"
	"database/sql"

	"github.com/kwhite17/Neighbors/pkg/database"
)

var createShelterQuery = "INSERT INTO shelters (City, Country, Name, PostalCode, State, Street) VALUES (?, ?, ?, ?, ?, ?)"
var deleteShelterQuery = "DELETE FROM shelters WHERE id=?"
var getSingleShelterQuery = "SELECT ID, City, Country, Name, PostalCode, State, Street from shelters where id=?"
var getAllSheltersQuery = "SELECT ID, City, Country, Name, PostalCode, State, Street from shelters"

type ShelterManager struct {
	ds database.Datasource
	database.DbManager
}

type Shelter struct {
	ID         int64
	City       string
	Country    string
	Name       string
	PostalCode string
	State      string
	Street     string
}

func (sm *ShelterManager) GetShelter(ctx context.Context, id int64) (*Shelter, error) {
	result, err := sm.ReadEntity(ctx, id)
	if err != nil {
		return nil, err
	}
	shelter, err := sm.buildShelters(result)
	if err != nil {
		return nil, err
	}
	return shelter[0], nil
}

func (sm *ShelterManager) GetShelters(ctx context.Context) ([]*Shelter, error) {
	result, err := sm.ReadEntities(ctx)
	if err != nil {
		return nil, err
	}
	shelters, err := sm.buildShelters(result)
	if err != nil {
		return nil, err
	}
	return shelters, nil
}

func (sm *ShelterManager) WriteShelter(ctx context.Context, shelter *Shelter) (int64, error) {
	values := []interface{}{shelter.City, shelter.Country, shelter.Name, shelter.PostalCode, shelter.State, shelter.Street}
	result, err := sm.WriteEntity(ctx, values)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func (sm *ShelterManager) DeleteShelter(ctx context.Context, id string) (int64, error) {
	result, err := sm.DeleteEntity(ctx, id)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (sm *ShelterManager) ReadEntity(ctx context.Context, id int64) (*sql.Rows, error) {
	return sm.ds.ExecuteReadQuery(ctx, getSingleShelterQuery, []interface{}{id})
}

func (sm *ShelterManager) ReadEntities(ctx context.Context) (*sql.Rows, error) {
	return sm.ds.ExecuteReadQuery(ctx, getAllSheltersQuery, nil)
}

func (sm *ShelterManager) WriteEntity(ctx context.Context, values []interface{}) (sql.Result, error) {
	result, err := sm.ds.ExecuteWriteQuery(ctx, createShelterQuery, values)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (sm *ShelterManager) DeleteEntity(ctx context.Context, id string) (sql.Result, error) {
	result, err := sm.ds.ExecuteWriteQuery(ctx, deleteShelterQuery, []interface{}{id})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (sm *ShelterManager) buildShelters(result *sql.Rows) ([]*Shelter, error) {
	response := make([]*Shelter, 0)
	for result.Next() {
		var id int64
		var city string
		var country string
		var name string
		var postalCode string
		var state string
		var street string
		if err := result.Scan(&id, &city, &country, &name, &postalCode, &state, &street); err != nil {
			return nil, err
		}
		shelter := Shelter{ID: id, City: city, Country: country, Name: name, PostalCode: postalCode, State: state, Street: street}
		response = append(response, &shelter)
	}
	return response, nil
}
