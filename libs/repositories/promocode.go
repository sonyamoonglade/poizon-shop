package repositories

import (
	"context"
	"errors"

	"domain"
	"dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type promocodeRepo struct {
	promocodes *mongo.Collection
}

func NewPromocodeRepo(promocodes *mongo.Collection) *promocodeRepo {
	return &promocodeRepo{
		promocodes: promocodes,
	}
}

func (p promocodeRepo) GetAll(ctx context.Context) ([]domain.Promocode, error) {
	cur, err := p.promocodes.Find(ctx, bson.D{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNoPromocodes
		}
		return nil, err
	}
	var promos []domain.Promocode
	if err := cur.All(ctx, &promos); err != nil {
		return nil, err
	}
	return promos, nil
}

func (p promocodeRepo) GetByID(ctx context.Context, promocodeID primitive.ObjectID) (domain.Promocode, error) {
	res := p.promocodes.FindOne(ctx, bson.M{"_id": promocodeID})
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Promocode{}, domain.ErrNoPromocode
		}
		return domain.Promocode{}, nil
	}
	var promo domain.Promocode
	return promo, res.Decode(&promo)
}

func (p promocodeRepo) Save(ctx context.Context, promo domain.Promocode) error {
	_, err := p.promocodes.InsertOne(ctx, promo)
	return err
}

func (p promocodeRepo) Delete(ctx context.Context, promocodeID primitive.ObjectID) error {
	_, err := p.promocodes.DeleteOne(ctx, bson.M{"_id": promocodeID})
	return err
}

func (p promocodeRepo) Update(ctx context.Context, promocodeID primitive.ObjectID, dto dto.UpdatePromocodeDTO) error {
	update := bson.M{}
	if dto.Description != nil {
		update["description"] = *dto.Description
	}
	if dto.Discounts != nil {
		update["discounts"] = *dto.Discounts
	}
	if dto.Counters != nil {
		update["counters"] = *dto.Counters
	}
	_, err := p.promocodes.UpdateOne(ctx, bson.M{"_id": promocodeID}, bson.M{"$set": update})
	return err
}

func (p promocodeRepo) GetByShortID(ctx context.Context, shortID string) (domain.Promocode, error) {
	res := p.promocodes.FindOne(ctx, bson.M{"shortId": shortID})
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Promocode{}, domain.ErrNoPromocode
		}
		return domain.Promocode{}, err
	}
	var promo domain.Promocode
	return promo, res.Decode(&promo)
}
