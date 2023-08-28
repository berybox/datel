package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"
)

const (
	// EnvURI key in environment variables to retrieve Mongo URI
	//EnvURI = "DATEL_MONGO_URI"

	// SettingsDB database with settings
	SettingsDB = "datel"

	// SettingsUsersCol collection with user settings
	SettingsUsersCol = "users"
)

var (
	currentInstance = &Instance{}

	// CtxTimeout default timeout for database connection
	CtxTimeout = 10 * time.Second
)

// Instance database connection instance
type Instance struct {
	Cli *mongo.Client
	URI string
}

// CurrentURI get URI of currently used Instance
func CurrentURI() string {
	return currentInstance.URI
}

// GetInstance gets the current instance so that it doesn't have to reconnect.
func GetInstance(URI string) (*Instance, error) {
	if currentInstance.Cli == nil || URI == "" || currentInstance.URI != URI {
		var err error
		currentInstance, err = createAndConnect(URI)
		if err != nil {
			return &Instance{}, err
		}
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), CtxTimeout)
		defer cancel()
		err := currentInstance.Cli.Ping(ctx, nil)
		if err != nil {
			currentInstance, err = createAndConnect(URI)
			if err != nil {
				return &Instance{}, err
			}
		}
	}

	return currentInstance, nil
}

// createAndConnect creates a mongoDB instance and connects to URI
func createAndConnect(URI string) (*Instance, error) {
	c := &Instance{}
	err := c.connectURI(URI)
	return c, err
}

// connectURI connect to the database with URI
func (m *Instance) connectURI(URI string) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxTimeout)
	defer cancel()

	var err error
	m.Cli, err = mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		return err
	}

	err = m.Cli.Ping(ctx, nil)
	if err != nil {
		return err
	}

	m.URI = URI

	return nil
}

// InsertOne inserts one document to the database
func (m *Instance) InsertOne(database, collection string, v interface{}) error {
	//*
	ctx, cancel := context.WithTimeout(context.Background(), CtxTimeout)
	defer cancel()
	//*/

	db := m.Cli.Database(database)

	_, err := db.Collection(collection).InsertOne(ctx, v)
	return err
}

// UnmarshallCollection retrieves data from the collection using a filter. Similar to JSON Unmarshall
func (m *Instance) UnmarshallCollection(database, collection string, filter interface{}, v interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxTimeout)
	defer cancel()

	db := m.Cli.Database(database)

	cur, err := db.Collection(collection).Find(ctx, filter)
	if err != nil {
		return err
	}

	err = cur.All(ctx, v)
	if err != nil {
		return err
	}

	return nil
}

// UpdateOne updates one document
func (m *Instance) UpdateOne(database, collection string, filter, update interface{}) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxTimeout)
	defer cancel()

	res, err := m.Cli.Database(database).Collection(collection).UpdateOne(ctx, filter, update)

	return res, err
}

// ReplaceOne replaces one document
func (m *Instance) ReplaceOne(database, collection string, filter, v interface{}) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxTimeout)
	defer cancel()

	res, err := m.Cli.Database(database).Collection(collection).ReplaceOne(ctx, filter, v)

	return res, err
}

// DeleteOne deletes one document
func (m *Instance) DeleteOne(database, collection string, filter interface{}) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxTimeout)
	defer cancel()

	res, err := m.Cli.Database(database).Collection(collection).DeleteOne(ctx, filter)

	return res, err
}

// HasCollection finds out if the database contains the given collection
func (m *Instance) HasCollection(database, collection string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxTimeout)
	defer cancel()

	cols, err := m.Cli.Database(database).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return false, err
	}

	return slices.Contains(cols, collection), nil
}

// DropCollection drops / deletes specified collection
func (m *Instance) DropCollection(database, collection string) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxTimeout)
	defer cancel()
	err := m.Cli.Database(database).Collection(collection).Drop(ctx)
	if err != nil {
		return err
	}
	return nil
}

// IDFilter creates a bson filter for filtering based on _id
func IDFilter(id string) bson.M {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return bson.M{"_id": id}
	}
	return bson.M{"_id": objID}
}
