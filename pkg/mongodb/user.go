package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// ErrUserNotFound user not found error
	ErrUserNotFound = errors.New("user not found")
)

// User user settings
type User struct {
	ID              string       `json:"_id" bson:"_id"`
	Username        string       `json:"username" bson:"username"`
	DefaultDatabase string       `json:"default_database" bson:"default_database"`
	Collections     []Collection `json:"collections" bson:"collections"`
}

// GetUser gets user settings from the database
func (m *Instance) GetUser(name string) (User, error) {
	filter := bson.D{bson.E{Key: "_id", Value: name}}

	var rets []User

	err := m.UnmarshallCollection(SettingsDB, SettingsUsersCol, filter, &rets)
	if err != nil {
		return User{}, err
	}

	if len(rets) == 0 {
		return User{}, ErrUserNotFound
	}

	return rets[0], nil
}

// CollectionByName find collection by name
func (u *User) CollectionByName(name string) (Collection, error) {
	i := u.CollectionIndex(name)
	if i >= 0 {
		return u.Collections[i], nil
	}
	return Collection{}, ErrCollectionNotFound
}

// CollectionIndex find collection index by name
func (u *User) CollectionIndex(name string) int {
	for i, c := range u.Collections {
		if c.Name == name {
			return i
		}
	}
	return -1
}

// UserUpdate update user settings
func (m *Instance) UserUpdate(userID string, update interface{}) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxTimeout)
	defer cancel()

	filter := bson.M{"_id": userID}

	col := m.Cli.Database(SettingsDB).Collection(SettingsUsersCol)
	res, err := col.UpdateOne(ctx, filter, update)

	return res, err
}

// userUpdateCollections finds the user by id and makes changes in the database in "collections" according to the current values of "user" variable
func (m *Instance) userUpdateCollections(user *User) error {
	_, err := m.UserUpdate(user.ID, bson.D{
		{Key: "$set", Value: bson.D{{Key: "collections", Value: user.Collections}}},
	})
	if err != nil {
		return err
	}

	return nil
}

// UserUpdateCollection updates collections of the user
func (m *Instance) UserUpdateCollection(user *User, col Collection) error {
	colIndex := user.CollectionIndex(col.Name)
	if colIndex < 0 {
		return ErrCollectionNotFound
	}

	user.Collections[colIndex] = col

	err := m.userUpdateCollections(user)
	if err != nil {
		return err
	}

	return nil
}

// UserAddCollection adds collection to the user
func (m *Instance) UserAddCollection(user *User, col Collection) error {
	user.Collections = append(user.Collections, col)

	err := m.userUpdateCollections(user)
	if err != nil {
		return err
	}

	return nil
}

// UserDeleteCollection deletes collection from user settings
func (m *Instance) UserDeleteCollection(user *User, col Collection) error {
	colIndex := user.CollectionIndex(col.Name)
	if colIndex < 0 {
		return ErrCollectionNotFound
	}

	// delete slice item by the colIndex
	user.Collections = append(user.Collections[:colIndex], user.Collections[colIndex+1:]...)

	err := m.userUpdateCollections(user)
	if err != nil {
		return err
	}

	return nil
}
