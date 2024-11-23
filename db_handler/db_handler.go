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
	cancel  context.CancelFunc
}

var (
	db_timeout time.Duration = 10
)

// DBClient methods

func (client *DBClient) Connect(host string, port int) error {

	ctx, cancel := context.WithTimeout(context.Background(), db_timeout*time.Second)

	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", host, port))
	handle, err := mongo.Connect(ctx, opts)
	if err != nil {
		cancel()
		return err
	}
	if err := handle.Ping(ctx, nil); err != nil {
		cancel()
		return err
	}
	client.context = ctx
	client.handle = handle
	client.cancel = cancel
	return nil
}

func (client *DBClient) Disconnect() error {
	return client.handle.Disconnect(client.context)
}

func (client *DBClient) GetDBNames() (dbs []string, err error) {
	return client.handle.ListDatabaseNames(client.context, bson.D{})
}

// functions

// sets timeout. should be called before Connect()
func SetTimeout(tmt int) {
	db_timeout = time.Duration(tmt)
}
