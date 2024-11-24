package db_handler

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDB  client type
type DBClient struct {
	handle      *mongo.Client      // db handle
	context     context.Context    // db context
	cancel      context.CancelFunc // context cancel function
	db          *mongo.Database    // db handler
	collections []*DBCollection    // collections
}

// DBClient methods

// disconnect from db
func (client *DBClient) Disconnect() error {
	return client.handle.Disconnect(client.context)
}

// get db names
func (client *DBClient) GetDBNames() (dbs []string, err error) {
	return client.handle.ListDatabaseNames(client.context, bson.D{})
}

// ping
func (client *DBClient) Ping() error {
	return client.handle.Ping(client.context, nil)
}

// cancel context (abandon all work immediately)
func (client *DBClient) CancelContext() {
	client.cancel()
}

// select a db, create one if it doesn't exist
func (client *DBClient) SelectDB(name string) error {
	if client.db = client.handle.Database(name); client.db == nil {
		return errors.New("Unable to select db " + name)
	}
	return nil
}

// get db name
func (client *DBClient) GetDBName() string {
	if client.db != nil {
		return client.db.Name()
	}
	return ""
}

// add collection
func (client *DBClient) AddCollection(name string) (*DBCollection, error) {
	if client.db == nil {
		return nil, errors.New("Unable to add collection " + name + ", DB is not selected")
	}
	collection := &DBCollection{
		mongo_collection: client.db.Collection(name),
		context:          &client.context,
	}
	client.collections = append(client.collections, collection)
	return collection, nil
}
