package db_handler

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB  client type
type DBClient struct {
	handle  *mongo.Client
	context context.Context
}

// DBClient methods

func (client *DBClient) Connect(host string, port int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", host, port))
	handle, err := mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}
	client.context = ctx
	client.handle = handle
	return nil
}

func (client *DBClient) Disconnect() error {
	return client.handle.Disconnect(client.context)
}

func (client *DBClient) GetDBNames() (dbs []string, err error) {
	return client.handle.ListDatabaseNames(client.context, bson.D{})
}
