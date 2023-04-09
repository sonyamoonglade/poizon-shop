package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"logger"
)

type Mongo struct {
	c  *mongo.Client
	db *mongo.Database
}

const retries = 5

func Connect(ctx context.Context, uri string, DBName string) (*Mongo, error) {
	// Uses connection pool
	opts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		for i := 0; i < retries; i++ {
			client, err := mongo.Connect(ctx, opts)
			logger.Get().Info("reconnecting to db")
			if err != nil {
				return nil, err
			}
			if client != nil {
				err = client.Ping(ctx, readpref.Primary())
				if err != nil {
					return nil, err
				}
				return &Mongo{c: client, db: client.Database(DBName)}, nil
			}
		}
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return &Mongo{c: client, db: client.Database(DBName)}, nil
}

func (m *Mongo) Collection(collection string) *mongo.Collection {
	return m.db.Collection(collection)
}

func (m *Mongo) Close(ctx context.Context) error {
	return m.c.Disconnect(ctx)
}
