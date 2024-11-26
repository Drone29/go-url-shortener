package db_handler

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB collection handler type
type DBCollection struct {
	mongo_collection *mongo.Collection
}

var ErrNoDocuments = mongo.ErrNoDocuments

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
	ctx, cancel := getContext()
	defer cancel()
	result, err := collection.mongo_collection.InsertOne(ctx, bsonDoc)
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
		return err
	}
	ctx, cancel := getContext()
	defer cancel()
	err = collection.mongo_collection.FindOne(ctx, bson_filter).Decode(&doc)
	if err != nil {
		return err
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
	ctx, cancel := getContext()
	defer cancel()
	cursor, err := collection.mongo_collection.Find(ctx, bsonD)
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	// Decode the results into the result slice
	if err := cursor.All(ctx, result); err != nil {
		return err
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
	ctx, cancel := getContext()
	defer cancel()
	// find doc
	var doc bson.M
	err = collection.mongo_collection.FindOne(ctx, filter).Decode(&doc)
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
	ctx, cancel := getContext()
	defer cancel()
	var doc bson.M
	err = collection.mongo_collection.FindOneAndDelete(ctx, filter).Decode(&doc)
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
	ctx, cancel := getContext()
	defer cancel()
	cursor, err := collection.mongo_collection.Find(ctx,
		bson.D{{}},
		options.Find().SetProjection(bson.D{{"_id", 1}}))
	if err != nil {
		return nil, fmt.Errorf("failed to find any docs: %v", err)
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and print the document IDs
	for cursor.Next(ctx) {
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
	ctx, cancel := getContext()
	defer cancel()
	return collection.mongo_collection.Drop(ctx)
}
