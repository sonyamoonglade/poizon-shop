package repositories

import (
	"context"
	"errors"

	"domain"
	"dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type clothingOrderRepo struct {
	orders *mongo.Collection
}

func NewClothingOrderRepo(orders *mongo.Collection) *clothingOrderRepo {
	return &clothingOrderRepo{
		orders: orders,
	}
}

func (o *clothingOrderRepo) AddComment(ctx context.Context, dto dto.AddCommentDTO) (domain.ClothingOrder, error) {
	filter := bson.M{"_id": dto.OrderID}
	update := bson.M{"$set": bson.M{"comment": dto.Comment}}
	return o.findOneAndUpdate(ctx, filter, update)
}

func (o *clothingOrderRepo) Approve(ctx context.Context, orderID primitive.ObjectID) (domain.ClothingOrder, error) {
	filter := bson.M{"_id": orderID}
	update := bson.M{"$set": bson.M{"isApproved": true}}
	return o.findOneAndUpdate(ctx, filter, update)
}

func (o *clothingOrderRepo) Delete(ctx context.Context, orderID primitive.ObjectID) error {
	filter := bson.M{"_id": orderID}
	_, err := o.orders.DeleteOne(ctx, filter)
	return err
}

func (o *clothingOrderRepo) ChangeStatus(ctx context.Context, dto dto.ChangeOrderStatusDTO) (domain.ClothingOrder, error) {
	filter := bson.M{"_id": dto.OrderID}
	update := bson.M{"$set": bson.M{"status": dto.NewStatus}}
	return o.findOneAndUpdate(ctx, filter, update)
}

func (o *clothingOrderRepo) GetAll(ctx context.Context) ([]domain.ClothingOrder, error) {
	findOpts := options.Find()
	findOpts.SetSort(bson.M{"isApproved": -1})
	res, err := o.orders.Find(ctx, bson.D{}, findOpts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNoOrders
		}
		return nil, err
	}
	var orders []domain.ClothingOrder
	if err := res.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *clothingOrderRepo) UpdateToPaid(ctx context.Context, customerID primitive.ObjectID, shortID string) error {
	filter := bson.M{"customer._id": customerID, "shortId": shortID}
	query := bson.M{
		"$set": bson.M{
			"isPaid": true,
		},
	}

	_, err := o.orders.UpdateOne(ctx, filter, query)
	if err != nil {
		return err
	}

	return nil
}
func (o *clothingOrderRepo) Save(ctx context.Context, order domain.ClothingOrder) error {
	_, err := o.orders.InsertOne(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

func (o *clothingOrderRepo) GetByShortID(ctx context.Context, shortID string) (domain.ClothingOrder, error) {
	res := o.orders.FindOne(ctx, bson.M{"shortId": shortID})
	err := res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.ClothingOrder{}, domain.ErrOrderNotFound
		}
		return domain.ClothingOrder{}, err
	}
	var ord domain.ClothingOrder
	if err := res.Decode(&ord); err != nil {
		return domain.ClothingOrder{}, err
	}
	return ord, nil
}

func (o *clothingOrderRepo) GetAllForCustomer(ctx context.Context, customerID primitive.ObjectID) ([]domain.ClothingOrder, error) {
	filter := bson.M{"customer._id": customerID}
	res, err := o.orders.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNoOrders
		}
		return nil, err
	}
	var orders []domain.ClothingOrder
	if err := res.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *clothingOrderRepo) findOneAndUpdate(ctx context.Context, filter, update any) (domain.ClothingOrder, error) {
	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(options.After)
	res := o.orders.FindOneAndUpdate(ctx, filter, update, opts)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return domain.ClothingOrder{}, domain.ErrOrderNotFound
		}
		return domain.ClothingOrder{}, res.Err()
	}
	var ord domain.ClothingOrder
	if err := res.Decode(&ord); err != nil {
		return domain.ClothingOrder{}, err
	}
	return ord, nil
}
