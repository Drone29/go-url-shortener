package db_handler

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// helpers

func createFilterWithID(id string) (bson.M, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	// filter to query by id
	return bson.M{"_id": objID}, nil
}

func convertBsonToJson(bson_doc bson.M, result any) error {
	jsonData, err := json.Marshal(bson_doc)
	if err != nil {
		return fmt.Errorf("failed to convert doc to JSON: %v", err)
	}
	return json.Unmarshal(jsonData, &result)
}

func bsonFromAny(s any) (bson.M, error) {
	data, err := bson.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct: %v", err)
	}
	// Unmarshal the BSON into a bson.M map
	var bsonMap bson.M
	err = bson.Unmarshal(data, &bsonMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal into bson.M: %v", err)
	}

	return bsonMap, nil
}
