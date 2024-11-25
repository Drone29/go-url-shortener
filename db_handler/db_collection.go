package db_handler

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB collection handler type
type DBCollection struct {
	mongo_collection *mongo.Collection
	context          *context.Context
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
func (collection *DBCollection) InsertOne(doc any) (id string, err error) {
	var bsonDoc any
	// convert doc to bson
	bsonDoc, err = bson.Marshal(doc)
	if err != nil {
		return "", fmt.Errorf("failed to convert struct to BSON: %v", err)
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

// find doc with filter
func (collection *DBCollection) FindOne(filter any, result any) error {
	var doc bson.M
	bson_filter, err := bsonFromAny(filter)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	err = collection.mongo_collection.FindOne(*collection.context, bson_filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("no docs found")
		}
		return fmt.Errorf("failed to find docs: %v", err)
	}

	return convertBsonToJson(doc, &result)
}

// find all with filter
func (collection *DBCollection) Find(filter any, result any) error {
	bson_filter, err := bsonFromAny(filter)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	var bsonD bson.D
	for key, value := range bson_filter {
		bsonD = append(bsonD, bson.E{Key: key, Value: value})
	}
	cursor, err := collection.mongo_collection.Find(*collection.context, bsonD)
	if err != nil {
		return fmt.Errorf("failed to find any docs: %v", err)
	}

	defer cursor.Close(*collection.context)

	// Decode the results into the result slice
	if err := cursor.All(*collection.context, result); err != nil {
		return fmt.Errorf("failed to decode documents: %v", err)
	}

	return nil
}

// find doc by id, and store into result
func (collection *DBCollection) FindByID(id string, result any) error {

	// filter to query by id
	filter, err := createFilterWithID(id)
	if err != nil {
		return fmt.Errorf("unable to create filter: %v", err)
	}

	// find doc
	var doc bson.M
	err = collection.mongo_collection.FindOne(*collection.context, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("doc %s not found", id)
		}
		return fmt.Errorf("failed to find doc %s: %v", id, err)
	}

	return convertBsonToJson(doc, &result)
}

// delete doc by id
func (collection *DBCollection) DeleteByID(id string) error {
	// filter to query by id
	filter, err := createFilterWithID(id)
	if err != nil {
		return fmt.Errorf("unable to create filter: %v", err)
	}
	var doc bson.M
	err = collection.mongo_collection.FindOneAndDelete(*collection.context, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("doc %s not found", id)
		}
		return fmt.Errorf("failed to find doc %s: %v", id, err)
	}
	return nil
}

// list all docs' ids
func (collection *DBCollection) ListDocsIDs() ([]string, error) {
	var ids []string
	cursor, err := collection.mongo_collection.Find(*collection.context,
		bson.D{{}},
		options.Find().SetProjection(bson.D{{"_id", 1}}))
	if err != nil {
		return nil, fmt.Errorf("failed to find any docs: %v", err)
	}
	defer cursor.Close(*collection.context)

	// Iterate through the cursor and print the document IDs
	for cursor.Next(*collection.context) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("decoding error: %v", err)
		}
		// Extract the _id from the result
		if id, ok := result["_id"].(primitive.ObjectID); ok {
			ids = append(ids, id.Hex())
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	return ids, nil
}

// drop collection
func (collection *DBCollection) Drop() error {
	return collection.mongo_collection.Drop(*collection.context)
}
