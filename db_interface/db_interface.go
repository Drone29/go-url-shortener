package db_interface

import "errors"

// db interface
type IDBCollection interface {
	InsertOne(doc any) (id string, err error)
	FindOne(filter any, result any) error
}

var ErrNoDocuments = errors.New("no documents")
