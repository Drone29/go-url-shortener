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

func bsonFromString(filter string) (bson.M, error) {
	var bson_filter bson.M
	err := json.Unmarshal([]byte(filter), &bson_filter)
	if err != nil {
		return nil, err
	}
	return bson_filter, nil
}

func bsonFromStruct(s any) (bson.M, error) {
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

func bsonFromAny(val any) (bson.M, error) {
	var bson_filter bson.M
	switch f := val.(type) {
	case string:
		var err error
		bson_filter, err = bsonFromString(f)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON filter: %v", err)
		}
		// Handle _id conversion (if it exists in the filter)
		if id, ok := bson_filter["_id"].(string); ok {
			// Convert the string _id to ObjectID
			objectID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				return nil, fmt.Errorf("invalid _id format: %v", err)
			}
			bson_filter["_id"] = objectID // Replace string _id with ObjectID
		}
	default:
		var err error
		bson_filter, err = bsonFromStruct(f)
		if err != nil {
			return nil, fmt.Errorf("invalid BSON filter: %v", err)
		}
	}
	return bson_filter, nil
}
