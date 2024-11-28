package db_handler

import (
	"fmt"
	"url-shortener/db_interface"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDB collection handler type (implements IDBCollection)
type DBCollection struct {
	mongo_collection *mongo.Collection
}

// helpers
func genUpdateDoc(orig, upd bson.M) bson.M {
	// diff bson
	diff := bson.M{}
	// iterate over newer version
	for k, newVal := range upd {
		if origVal, exists := orig[k]; !exists || origVal != newVal {
			diff[k] = newVal
		}
	}
	// If diff is empty, return nil
	if len(diff) == 0 {
		return nil
	}
	// mongo needs an update doc with $ key
	return bson.M{"$set": diff}
}

// DBCollection methods

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
		if err == mongo.ErrNoDocuments {
			return db_interface.ErrNoDocuments
		}
		return err
	}

	return convertBsonToJson(doc, &result)
}

// update doc
func (collection *DBCollection) UpdateOne(old any, new any) error {
	old_doc, err := bsonFromAny(old)
	if err != nil {
		return err
	}
	new_doc, err := bsonFromAny(new)
	if err != nil {
		return err
	}
	// generate update doc (with $set)
	update := genUpdateDoc(old_doc, new_doc)
	if update == nil {
		return fmt.Errorf("no changes introduced")
	}
	ctx, cancel := getContext()
	defer cancel()
	res, err := collection.mongo_collection.UpdateOne(ctx, old_doc, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return db_interface.ErrNoDocuments
		}
		return err
	}

	if res.ModifiedCount == 0 {
		return fmt.Errorf("no documents were updated")
	}

	return nil
}

// delete doc
func (collection *DBCollection) DeleteOne(filter any) error {
	bson_filter, err := bsonFromAny(filter)
	if err != nil {
		return err
	}
	var res bson.M
	ctx, cancel := getContext()
	defer cancel()
	err = collection.mongo_collection.FindOneAndDelete(ctx, bson_filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return db_interface.ErrNoDocuments
		}
		return err
	}
	return nil
}
