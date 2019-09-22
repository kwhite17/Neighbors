package managers

import (
	"context"
	"database/sql"

	"github.com/kwhite17/Neighbors/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

var createShelterQuery = "INSERT INTO shelters (City, Country, Name, Password, PostalCode, State, Street) VALUES ($1, $2, $3, $4, $5, $6, $7)"
var deleteShelterQuery = "DELETE FROM shelters WHERE id=$1"
var getSingleShelterQuery = "SELECT ID, City, Country, Name, PostalCode, State, Street FROM shelters where id=$1"
var getAllSheltersQuery = "SELECT ID, City, Country, Name, PostalCode, State, Street FROM shelters"
var updateShelterQuery = "UPDATE shelters SET City = $1, Country = $2, Name = $3, PostalCode = $4, State = $5, Street = $6 WHERE ID = $7"
var getPasswordForUsernameQuery = "SELECT ID, Password FROM shelters WHERE Name = $1"

type ShelterManager struct {
	Datasource database.Datasource
	database.DbManager
}

type ContactInformation struct {
	City       string
	Country    string
	Name       string
	PostalCode string
	State      string
	Street     string
}

type Shelter struct {
	ID       int64
	Password string
	*ContactInformation
}

func (sm *ShelterManager) GetShelter(ctx context.Context, id interface{}) (*Shelter, error) {
	result, err := sm.ReadEntity(ctx, id)

	if err != nil {
		return nil, err
	}

	shelter, err := sm.buildShelters(result)

	if err != nil {
		return nil, err
	}

	if len(shelter) < 1 {
		return nil, nil
	}

	return shelter[0], nil
}

func (sm *ShelterManager) GetPasswordForUsername(ctx context.Context, username string) (*Shelter, error) {
	row := sm.Datasource.ExecuteSingleReadQuery(ctx, getPasswordForUsernameQuery, []interface{}{username})

	var ID int64
	var password string
	if err := row.Scan(&ID, &password); err != nil {
		return nil, err
	}
	shelter := Shelter{ID: ID, Password: password}
	return &shelter, nil
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

func (sm *ShelterManager) WriteShelter(ctx context.Context, shelter *Shelter, unencryptedPassword string) (int64, error) {
	encryptedPassword, err := sm.encryptPassword(unencryptedPassword)
	if err != nil {
		return -1, err
	}

	values := []interface{}{shelter.City, shelter.Country, shelter.Name, encryptedPassword, shelter.PostalCode, shelter.State, shelter.Street}
	result, err := sm.WriteEntity(ctx, values, true)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func (sm *ShelterManager) UpdateShelter(ctx context.Context, shelter *Shelter) error {
	values := []interface{}{shelter.City, shelter.Country, shelter.Name, shelter.PostalCode, shelter.State, shelter.Street, shelter.ID}
	_, err := sm.Datasource.ExecuteWriteQuery(ctx, updateShelterQuery, values, true)
	return err
}

func (sm *ShelterManager) DeleteShelter(ctx context.Context, id interface{}) (int64, error) {
	result, err := sm.DeleteEntity(ctx, id, true)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (sm *ShelterManager) ReadEntity(ctx context.Context, id interface{}) (*sql.Rows, error) {
	return sm.Datasource.ExecuteBatchReadQuery(ctx, getSingleShelterQuery, []interface{}{id})
}

func (sm *ShelterManager) ReadEntities(ctx context.Context) (*sql.Rows, error) {
	return sm.Datasource.ExecuteBatchReadQuery(ctx, getAllSheltersQuery, nil)
}

func (sm *ShelterManager) WriteEntity(ctx context.Context, values []interface{}, returnResult bool) (sql.Result, error) {
	result, err := sm.Datasource.ExecuteWriteQuery(ctx, createShelterQuery, values, returnResult)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (sm *ShelterManager) DeleteEntity(ctx context.Context, id interface{}, returnResult bool) (sql.Result, error) {
	result, err := sm.Datasource.ExecuteWriteQuery(ctx, deleteShelterQuery, []interface{}{id}, returnResult)
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
		contactInfo := &ContactInformation{City: city, Country: country, Name: name, PostalCode: postalCode, State: state, Street: street}
		shelter := Shelter{ID: id, ContactInformation: contactInfo}
		response = append(response, &shelter)
	}
	return response, nil
}

func (sm *ShelterManager) encryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
