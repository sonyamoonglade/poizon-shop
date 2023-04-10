package database

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	ctx := context.Background()
	t.Run("count documents", func(t *testing.T) {
		col := db.Collection("test-col1")
		t.Cleanup(func() {
			col.Drop(ctx)
		})
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
		err = db.WithTransaction(ctx, func(tx context.Context) error {
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
		require.NoError(t, err)
		time.Sleep(time.Millisecond * 500)
		n, err := col.CountDocuments(ctx, bson.D{})
		require.NoError(t, err)

		// Tx snapshot with inserted doc within goroutine results in 3 documents
		require.Equal(t, int64(3), n)
	})

	t.Run("test with snapshot counter", func(t *testing.T) {
		N := 10
		col := db.Collection("test-col1")
		t.Cleanup(func() {
			col.Drop(ctx)
		})
		key := "counter"
		res, err := col.InsertOne(ctx, bson.M{key: 0})
		require.NoError(t, err)

		id := res.InsertedID
		doChan := make(chan struct{})
		assertChan := make(chan struct{})
		go func() {
			<-doChan
			res := col.FindOne(ctx, bson.M{"_id": id})
			type Resp struct {
				Counter int `bson:"counter"`
			}
			var resp Resp
			require.NoError(t, res.Decode(&resp))
			// transaction has not committed yet
			require.Equal(t, 0, resp.Counter)
			assertChan <- struct{}{}
		}()
		err = db.WithTransaction(ctx, func(tx context.Context) error {
			for i := 0; i < N; i++ {
				_, err := col.UpdateOne(tx, bson.M{"_id": id.(primitive.ObjectID)}, bson.M{"$set": bson.M{key: i + 1}})
				if err != nil {
					return err
				}
			}
			// Run check on counter
			doChan <- struct{}{}
			<-assertChan
			return nil
		})
		require.NoError(t, err)

		// Check after transaction has committed
		afterTxRes := col.FindOne(ctx, bson.M{"_id": id})
		type Resp struct {
			Counter int `bson:"counter"`
		}
		var resp Resp
		require.NoError(t, afterTxRes.Decode(&resp))
		require.Equal(t, N, resp.Counter)
	})

	t.Run("test correct error return", func(t *testing.T) {
		col := db.Collection("test-col1")
		t.Cleanup(func() {
			col.Drop(ctx)
		})

		err := db.WithTransaction(ctx, func(tx context.Context) error {
			return fmt.Errorf("my error")
		})
		require.Equal(t, "my error", err.Error())
	})
	db.Close(ctx)
}
