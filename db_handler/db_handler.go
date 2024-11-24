package db_handler

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db_timeout time.Duration = 30
)

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
