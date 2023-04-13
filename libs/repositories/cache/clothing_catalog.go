package cache

import (
	"context"
	"domain"
	"redis"
	"repositories"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type clothingCatalogCached struct {
	client *redis.Client

	delegate repositories.ClothingCatalog
}

const gtpKey = "gtp-1"

func (c *clothingCatalogCached) GetTopRank(ctx context.Context) (uint, error) {
	v, ok, err := c.client.Get(ctx, gtpKey)
	if err != nil {
		// handle err
		return c.delegate.GetTopRank(ctx)
	}
	if !ok {
		return c.delegate.GetTopRank(ctx)
	}
	rank, err := strconv.ParseUint(string(v), 10, 64)
	if err != nil {
		// handler err
		return c.delegate.GetTopRank(ctx)
	}
	return uint(rank), nil
}

const gckey = "gc-1"

func (c *clothingCatalogCached) GetCatalog(ctx context.Context) ([]domain.ClothingProduct, error) {
	panic("not implemented") // TODO: Implement
}

func (c *clothingCatalogCached) GetIDByRank(ctx context.Context, rank uint) (primitive.ObjectID, error) {
	panic("not implemented") // TODO: Implement
}

func (c *clothingCatalogCached) GetRankByID(ctx context.Context, itemID primitive.ObjectID) (uint, error) {
	panic("not implemented") // TODO: Implement
}

func (c *clothingCatalogCached) GetLastRank(ctx context.Context) (uint, error) {
	panic("not implemented") // TODO: Implement
}

func (c *clothingCatalogCached) AddItem(ctx context.Context, item domain.ClothingProduct) error {
	panic("not implemented") // TODO: Implement
}

func (c *clothingCatalogCached) RemoveItem(ctx context.Context, itemID primitive.ObjectID) error {
	panic("not implemented") // TODO: Implement
}
