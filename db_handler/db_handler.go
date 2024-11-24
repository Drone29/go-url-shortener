package db_handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB  client type
type DBClient struct {
	handle      *mongo.Client      // db handle
	context     context.Context    // db context
	cancel      context.CancelFunc // context cancel function
	db          *mongo.Database    // db handler
	collections []*DBCollection    // collections
}

// MongoDB collection handler type
type DBCollection struct {
	mongo_collection *mongo.Collection
	context          *context.Context
}

var (
	db_timeout time.Duration = 30
)

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

// DBCollection methods

// get collection name
func (collection *DBCollection) GetName() string {
	if collection.mongo_collection != nil {
		return collection.mongo_collection.Name()
	}
	return ""
}

// insert one doc into collection
func (collection *DBCollection) InsertOne(doc interface{}) (id string, err error) {
	var bsonDoc interface{}
	// convert doc to bson
	switch d := doc.(type) {
	case string, []byte:
		// If `doc` is a JSON string/[]byte, parse it into a map
		var temp map[string]interface{}
		var jsonData []byte
		if str, ok := d.(string); ok {
			jsonData = []byte(str)
		} else {
			jsonData = d.([]byte)
		}
		err = json.Unmarshal(jsonData, &temp)
		if err != nil {
			return "", fmt.Errorf("failed to parse JSON: %v", err)
		}
		bsonDoc, err = bson.Marshal(temp)
		if err != nil {
			return "", fmt.Errorf("failed to convert JSON to BSON: %v", err)
		}
	default:
		// already a struct that can be converted to bson
		bsonDoc, err = bson.Marshal(d)
		if err != nil {
			return "", fmt.Errorf("failed to convert struct to BSON: %v", err)
		}
	}

	result, err := collection.mongo_collection.InsertOne(*collection.context, bsonDoc)
	if err == nil {
		if objectID, ok := result.InsertedID.(primitive.ObjectID); ok {
			id = objectID.Hex()
		} else {
			err = fmt.Errorf("invalid return type %T", objectID)
		}
	}
	return id, err
}

// functions

// sets timeout. should be called before Connect()
func SetTimeout(tmt int) {
	db_timeout = time.Duration(tmt)
}

// connect to db
func Connect(host string, port int) (*DBClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db_timeout*time.Second)

	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d/", host, port))
	handle, err := mongo.Connect(ctx, opts)
	if err != nil {
		cancel()
		return nil, err
	}
	if err := handle.Ping(ctx, nil); err != nil {
		cancel()
		return nil, err
	}

	return &DBClient{
		handle:  handle,
		context: ctx,
		cancel:  cancel,
	}, err
}

// find doc by id
func FindByID[T any](collection *DBCollection, id string) (T, error) {
	var result T
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, fmt.Errorf("invalid ID format: %v", err)
	}

	// filter to query by id
	filter := bson.M{"_id": objID}

	// find doc
	var doc bson.M
	err = collection.mongo_collection.FindOne(*collection.context, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, fmt.Errorf("doc %s not found", id)
		}
		return result, fmt.Errorf("failed to find doc %s: %v", id, err)
	}

	jsonData, err := json.Marshal(doc)
	if err != nil {
		return result, fmt.Errorf("failed to convert doc to JSON: %v", err)
	}
	err = json.Unmarshal(jsonData, &result)
	return result, err
}
