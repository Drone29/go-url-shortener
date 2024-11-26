package db_handler

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDB  client type
type DBClient struct {
	handle      *mongo.Client   // db handle
	db          *mongo.Database // db handler
	collections []*DBCollection // collections
}

// DBClient methods

// disconnect from db
func (client *DBClient) Disconnect() error {
	return client.handle.Disconnect(context.Background())
}

// get db names
func (client *DBClient) GetDBNames() (dbs []string, err error) {
	ctx, cancel := getContext()
	defer cancel()
	return client.handle.ListDatabaseNames(ctx, bson.D{})
}

// ping
func (client *DBClient) Ping() error {
	ctx, cancel := getContext()
	defer cancel()
	return client.handle.Ping(ctx, nil)
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

// get collection (create if doesn't exist)
func (client *DBClient) GetCollection(name string) (*DBCollection, error) {
	if client.db == nil {
		return nil, errors.New("Unable to add collection " + name + ", DB is not selected")
	}
	collection := &DBCollection{
		mongo_collection: client.db.Collection(name),
	}
	client.collections = append(client.collections, collection)
	return collection, nil
}
