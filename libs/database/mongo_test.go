package database

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestTransactor(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	mongoURI, dbName := os.Getenv("MONGO_URI"), os.Getenv("DB_NAME")
	if mongoURI == "" || dbName == "" {
		t.Fatalf("empty env's\n")
	}
	db, err := Connect(context.Background(), mongoURI, dbName)
	require.NoError(t, err)

	col := db.Collection("test-col1")
	ctx := context.Background()
	t.Cleanup(func() {
		col.Drop(ctx)
	})
	t.Run("test transaction", func(t *testing.T) {
		// Insert one document
		_, err := col.InsertOne(ctx, bson.M{"key": "value"})
		require.NoError(t, err)

		doChan := make(chan struct{})
		afterInsertChan := make(chan struct{})
		go func() {
			<-doChan
			// when triggered, insert one document without transaction
			_, err := col.InsertOne(ctx, bson.M{"key3": "value3"})
			require.NoError(t, err)
			afterInsertChan <- struct{}{}
		}()
		// Start countDocument within transaction and expect one document
		db.WithTransaction(ctx, func(tx mongo.SessionContext) error {
			n, err := col.CountDocuments(tx, bson.D{})
			if err != nil {
				return fmt.Errorf("count docs: %w", err)
			}
			require.Equal(t, int64(1), n)
			if _, err := col.InsertOne(tx, bson.M{"key1": "value2"}); err != nil {
				return fmt.Errorf("insert one: %w", err)
			}
			doChan <- struct{}{}
			<-afterInsertChan
			n, err = col.CountDocuments(tx, bson.D{})
			if err != nil {
				return fmt.Errorf("count docs: %w", err)
			}
			require.Equal(t, int64(2), n)
			// Here it commits
			return nil
		})
		time.Sleep(time.Millisecond * 500)
		n, err := col.CountDocuments(ctx, bson.D{})
		require.NoError(t, err)

		// Tx snapshot with inserted doc within gorutine results in 3 documents
		require.Equal(t, int64(3), n)
	})

	db.Close(ctx)
}
