package repositories

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"household_bot/pkg/telegram"
)

type householdCatalogMsgRepo struct {
	msgs *mongo.Collection
}

func NewHouseholdCatalogMsgRepo(msgs *mongo.Collection) *householdCatalogMsgRepo {
	return &householdCatalogMsgRepo{
		msgs: msgs,
	}
}

func (h householdCatalogMsgRepo) Save(ctx context.Context, m telegram.CatalogMsg) error {
	_, err := h.msgs.InsertOne(ctx, m)
	if err != nil {
		// Fine because no need to store duplicates
		if mongo.IsDuplicateKeyError(err) {
			return nil
		}
		return err
	}
	return nil
}

func (h householdCatalogMsgRepo) GetAll(ctx context.Context) ([]telegram.CatalogMsg, error) {
	cur, err := h.msgs.Find(ctx, bson.D{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("find: %w", err)
	}
	msgs := make([]telegram.CatalogMsg, 0, cur.RemainingBatchLength())
	if err := cur.All(ctx, &msgs); err != nil {
		return nil, fmt.Errorf("cur all: %w", err)
	}
	return msgs, nil
}

func (h householdCatalogMsgRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := h.msgs.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (h householdCatalogMsgRepo) DeleteByMsgID(ctx context.Context, msgID int) error {
	_, err := h.msgs.DeleteOne(ctx, bson.M{"msgId": msgID})
	return err
}
