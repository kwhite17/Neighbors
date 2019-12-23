package managers

import (
	"context"
	"database/sql"

	"github.com/kwhite17/Neighbors/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

var createUserQuery = "INSERT INTO users (City, Country, Name, Password, PostalCode, State, Street, UserType) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
var deleteUserQuery = "DELETE FROM users WHERE ID=$1"
var getSingleUserQuery = "SELECT ID, City, Country, Name, PostalCode, State, Street, UserType FROM users where ID=$1"
var getAllSheltersQuery = "SELECT ID, City, Country, Name, PostalCode, State, Street, UserType FROM users WHERE UserType=1"
var updateUserQuery = "UPDATE users SET City = $1, Country = $2, Name = $3, PostalCode = $4, State = $5, Street = $6 WHERE ID = $7"
var getPasswordForUsernameQuery = "SELECT ID, Password FROM users WHERE Name = $1"

type UserManager struct {
	Datasource database.Datasource
}

type UserType int

const (
	SHELTER   UserType = 1
	SAMARITAN UserType = 2
)

type ContactInformation struct {
	City       string
	Country    string
	Name       string
	PostalCode string
	State      string
	Street     string
}

type User struct {
	ID       int64
	Password string
	UserType UserType
	*ContactInformation
}

func (um *UserManager) GetUser(ctx context.Context, id interface{}) (*User, error) {
	result, err := um.Datasource.ExecuteBatchReadQuery(ctx, getSingleUserQuery, []interface{}{id})

	if err != nil {
		return nil, err
	}

	user, err := um.buildUsers(result)

	if err != nil {
		return nil, err
	}

	if len(user) < 1 {
		return nil, nil
	}

	return user[0], nil
}

func (um *UserManager) GetPasswordForUsername(ctx context.Context, username string) (*User, error) {
	row := um.Datasource.ExecuteSingleReadQuery(ctx, getPasswordForUsernameQuery, []interface{}{username})

	var ID int64
	var password string
	if err := row.Scan(&ID, &password); err != nil {
		return nil, err
	}
	user := User{ID: ID, Password: password}
	return &user, nil
}

func (um *UserManager) GetUsers(ctx context.Context) ([]*User, error) {
	result, err := um.Datasource.ExecuteBatchReadQuery(ctx, getAllSheltersQuery, nil)
	if err != nil {
		return nil, err
	}
	users, err := um.buildUsers(result)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (um *UserManager) WriteUser(ctx context.Context, user *User, unencryptedPassword string) (int64, error) {
	encryptedPassword, err := um.encryptPassword(unencryptedPassword)
	if err != nil {
		return -1, err
	}

	values := []interface{}{user.City, user.Country, user.Name, encryptedPassword, user.PostalCode, user.State, user.Street, SHELTER}
	result, err := um.Datasource.ExecuteWriteQuery(ctx, createUserQuery, values, true)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func (um *UserManager) UpdateUser(ctx context.Context, user *User) error {
	values := []interface{}{user.City, user.Country, user.Name, user.PostalCode, user.State, user.Street, user.ID}
	_, err := um.Datasource.ExecuteWriteQuery(ctx, updateUserQuery, values, true)
	return err
}

func (um *UserManager) DeleteUser(ctx context.Context, id interface{}) (int64, error) {
	result, err := um.Datasource.ExecuteWriteQuery(ctx, deleteUserQuery, []interface{}{id}, true)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (um *UserManager) buildUsers(result *sql.Rows) ([]*User, error) {
	response := make([]*User, 0)
	for result.Next() {
		var id int64
		var city string
		var country string
		var name string
		var postalCode string
		var state string
		var street string
		var userType int
		if err := result.Scan(&id, &city, &country, &name, &postalCode, &state, &street, &userType); err != nil {
			return nil, err
		}
		contactInfo := &ContactInformation{City: city, Country: country, Name: name, PostalCode: postalCode, State: state, Street: street}
		user := User{ID: id, ContactInformation: contactInfo, UserType: UserType(userType)}
		response = append(response, &user)
	}
	return response, nil
}

func (um *UserManager) encryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
