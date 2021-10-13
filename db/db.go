package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const dbURI = "mongodb+srv://gvozdev:1993gvozdev1993@gvozdev-cluster0.0cxjj.mongodb.net/test?authSource=admin&replicaSet=atlas-2koq8f-shard-0&readPreference=primary&appname=MongoDB%20Compass&ssl=true"

// ClientDB содержит пул подключений к БД.
var ClientDB *mongo.Client

// ConnectAndPingToClientDB создает нового клиента для подключения его к БД,
// устанавливает это соединение и проверяет его живость.
func ConnectAndPingToClientDB() (*mongo.Client, error) {

	var err error

	client, err := mongo.NewClient(options.Client().ApplyURI(dbURI))
	if err != nil {
		return nil, err
	}

	if err = client.Connect(context.Background()); err != nil {
		return nil, err
	}

	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}