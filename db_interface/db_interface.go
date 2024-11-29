package db_interface

import "errors"

// db interface
type IDBCollection interface {
	InsertOne(doc any) (id string, err error)
	FindOne(filter any, result any) error
	UpdateOne(filter any, update_with any) error
	DeleteOne(filter any) error
}

var ErrNoDocuments = errors.New("no records found")
