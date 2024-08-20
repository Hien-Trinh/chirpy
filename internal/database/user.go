package database

import (
	"errors"
	"sort"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	uniqueId := len(dbStructure.Users) + 1

	user := User{
		Id:       uniqueId,
		Email:    email,
		Password: password,
	}

	dbStructure.Users[user.Id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// GetUsers returns all users in the database
func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		users = append(users, user)
	}

	sort.Slice(users, func(i, j int) bool { return users[i].Id < users[j].Id })

	return users, nil
}

// GetUserById returns user with matching id in the database
func (db *DB) GetUserById(i int) (User, error) {
	user := User{}
	dbStructure, err := db.loadDB()
	if err != nil {
		return user, err
	}

	user, ok := dbStructure.Users[i]
	if !ok {
		return user, errors.New("User not found")
	}

	return user, nil
}

func (db *DB) UpdateUser(i int, new_email, new_password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	_, ok := dbStructure.Users[i]
	if !ok {
		return User{}, errors.New("User not found")
	}

	new_user := User{
		Id:       i,
		Email:    new_email,
		Password: new_password,
	}
	dbStructure.Users[i] = new_user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return new_user, nil
}
