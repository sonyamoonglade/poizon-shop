package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type rateRepo struct {
	rate *mongo.Collection
}

func NewRateRepository(rate *mongo.Collection) *rateRepo {
	return &rateRepo{
		rate: rate,
	}
}

func (r *rateRepo) GetYuanRate(ctx context.Context) (float64, error) {
	cur, err := r.rate.Find(ctx, bson.D{})
	if err != nil {
		return 0, err
	}
	type rateResponse struct {
		Rate float64 `bson:"rate"`
	}
	var resp rateResponse
	if cur.Next(ctx) {
		if err := cur.Decode(&resp); err != nil {
			return 0, err
		}
	}
	return resp.Rate, nil
}

func (r *rateRepo) UpdateRate(ctx context.Context, rate float64) error {
	opts := options.Update()
	opts.SetUpsert(true)
	_, err := r.rate.UpdateOne(ctx, bson.D{}, bson.M{"$set": bson.M{
		"rate": rate,
	}}, opts)
	if err != nil {
		return err
	}
	return nil
}
