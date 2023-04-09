package repositories

import (
	"context"
	"errors"
	"io"

	"domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"logger"
)

type ClothingOnChangeFunc func(items []domain.ClothingProduct)

type clothingCatalogRepo struct {
	catalog  *mongo.Collection
	onChange ClothingOnChangeFunc
}

func NewClothingCatalogRepo(catalog *mongo.Collection, onChangeFunc ClothingOnChangeFunc) *clothingCatalogRepo {
	return &clothingCatalogRepo{
		catalog:  catalog,
		onChange: onChangeFunc,
	}
}

func (c *clothingCatalogRepo) GetIDByRank(ctx context.Context, rank uint) (primitive.ObjectID, error) {
	res := c.catalog.FindOne(ctx, bson.M{"rank": rank})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return primitive.ObjectID{}, domain.ErrItemNotFound
		}
		return primitive.ObjectID{}, res.Err()
	}
	var item domain.ClothingProduct
	if err := res.Decode(&item); err != nil {
		return primitive.ObjectID{}, err
	}
	return item.ItemID, nil

}

func (c *clothingCatalogRepo) GetRankByID(ctx context.Context, itemID primitive.ObjectID) (uint, error) {
	res := c.catalog.FindOne(ctx, bson.M{"_id": itemID})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return 0, domain.ErrItemNotFound
		}
		return 0, res.Err()
	}
	var item domain.ClothingProduct
	if err := res.Decode(&item); err != nil {
		return 0, err
	}
	return item.Rank, nil
}

func (c *clothingCatalogRepo) GetLastRank(ctx context.Context) (uint, error) {
	filter := bson.D{}
	opts := options.Find()
	opts.SetSort(bson.M{"rank": -1})
	opts.SetLimit(1)
	res, err := c.catalog.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, domain.ErrItemNotFound
		}
		if errors.Is(err, io.EOF) {
			return 0, nil
		}
		return 0, err
	}
	var item domain.ClothingProduct
	if err := res.Decode(&item); err != nil {
		if errors.Is(err, io.EOF) {
			return 0, nil
		}
		return 0, err
	}
	return item.Rank, nil
}

func (c *clothingCatalogRepo) GetCatalog(ctx context.Context) ([]domain.ClothingProduct, error) {
	opts := options.Find()
	opts.SetSort(bson.M{"rank": 1})
	res, err := c.catalog.Find(ctx, bson.D{}, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNoCatalog
		}
		return nil, err
	}

	var catalog []domain.ClothingProduct
	defer res.Close(ctx)
	if err := res.All(ctx, &catalog); err != nil {
		return nil, err
	}

	return catalog, nil
}

func (c *clothingCatalogRepo) AddItem(ctx context.Context, item domain.ClothingProduct) error {
	_, err := c.catalog.InsertOne(ctx, item)
	if err != nil {
		return err
	}
	defer func() {
		// Notify on change
		newCatalog, err := c.GetCatalog(ctx)
		if err != nil {
			logger.Get().Error("deferred catalog notify", zap.Error(err))
			return
		}
		if c.onChange != nil {
			c.onChange(newCatalog)
		}
	}()
	return nil
}

func (c *clothingCatalogRepo) RemoveItem(ctx context.Context, itemID primitive.ObjectID) error {
	_, err := c.catalog.DeleteOne(ctx, bson.M{"_id": itemID})
	if err != nil {
		return err
	}
	// TODO: move to service
	catalog, err := c.GetCatalog(ctx)
	if err != nil {
		return err
	}

	// Update ranks
	// TODO: move to service(bl)
	newCatalog := domain.UpdateRanks(catalog)
	for _, newItem := range newCatalog {
		if _, err := c.catalog.UpdateOne(ctx, bson.M{"_id": newItem.ItemID}, bson.M{"$set": bson.M{"rank": newItem.Rank}}); err != nil {
			return err
		}
	}

	defer func() {
		// Notify on change
		newCatalog, err := c.GetCatalog(ctx)
		if err != nil {
			logger.Get().Error("deferred catalog notify", zap.Error(err))
			return
		}
		if c.onChange != nil {
			c.onChange(newCatalog)
		}
	}()
	return nil
}
