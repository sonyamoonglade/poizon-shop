package database

import (
	"context"
	"fmt"

	"logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

type txFunc func(tx mongo.SessionContext) error
type Transactor interface {
	WithTransaction(ctx context.Context, f txFunc) error
}

func (m *Mongo) WithTransaction(ctx context.Context, txFn txFunc) error {
	sess, err := m.db.Client().StartSession()
	if err != nil {
		return fmt.Errorf("start session: %w", err)
	}
	if err := sess.StartTransaction(); err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}
	defer sess.EndSession(context.Background())
	_, err = sess.WithTransaction(ctx, func(tx mongo.SessionContext) (interface{}, error) {
		if err := txFn(tx); err != nil {
			return nil, tx.AbortTransaction(context.Background())
		}
		return nil, tx.CommitTransaction(context.Background())
	})
	return err
}
