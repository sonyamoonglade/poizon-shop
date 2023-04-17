package repositories

import (
	"context"
	"errors"
	"logger"

	"domain"
	"dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type orderRepo[T domain.ClothingOrder | domain.HouseholdOrder] struct {
	orders *mongo.Collection
}

func NewClothingOrderRepo(orders *mongo.Collection) *orderRepo[domain.ClothingOrder] {
	repo := orderRepo[domain.ClothingOrder]{
		orders: orders,
	}
	return &repo
}
func NewHouseholdOrderRepo(orders *mongo.Collection) *orderRepo[domain.HouseholdOrder] {
	repo := orderRepo[domain.HouseholdOrder]{
		orders: orders,
	}
	return &repo
}

func (o *orderRepo[T]) AddComment(ctx context.Context, dto dto.AddCommentDTO) (T, error) {
	filter := bson.M{"_id": dto.OrderID}
	update := bson.M{"$set": bson.M{"comment": dto.Comment}}
	return o.findOneAndUpdate(ctx, filter, update)
}

func (o *orderRepo[T]) Approve(ctx context.Context, orderID primitive.ObjectID) (T, error) {
	filter := bson.M{"_id": orderID}
	update := bson.M{"$set": bson.M{"isApproved": true}}
	return o.findOneAndUpdate(ctx, filter, update)
}

func (o *orderRepo[T]) Delete(ctx context.Context, orderID primitive.ObjectID) error {
	filter := bson.M{"_id": orderID}
	_, err := o.orders.DeleteOne(ctx, filter)
	return err
}

func (o *orderRepo[T]) ChangeStatus(ctx context.Context, dto dto.ChangeOrderStatusDTO) (T, error) {
	filter := bson.M{"_id": dto.OrderID}
	update := bson.M{"$set": bson.M{"status": dto.NewStatus}}
	return o.findOneAndUpdate(ctx, filter, update)
}

func (o *orderRepo[T]) GetAll(ctx context.Context) ([]T, error) {
	findOpts := options.Find()
	findOpts.SetSort(bson.M{"isApproved": -1})
	res, err := o.orders.Find(ctx, bson.D{}, findOpts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNoOrders
		}
		return nil, err
	}
	var orders []T
	if err := res.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *orderRepo[T]) UpdateToPaid(ctx context.Context, customerID primitive.ObjectID, shortID string) error {
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
func (o *orderRepo[T]) Save(ctx context.Context, order T) error {
	_, err := o.orders.InsertOne(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

func (o *orderRepo[T]) GetByShortID(ctx context.Context, shortID string) (T, error) {
	res := o.orders.FindOne(ctx, bson.M{"shortId": shortID})
	err := res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return *new(T), domain.ErrOrderNotFound
		}
		return *new(T), err
	}
	var ord T
	if err := res.Decode(&ord); err != nil {
		return *new(T), err
	}
	return ord, nil
}

func (o *orderRepo[T]) GetAllForCustomer(ctx context.Context, customerID primitive.ObjectID) ([]T, error) {
	filter := bson.M{"customer._id": customerID}
	res, err := o.orders.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNoOrders
		}
		return nil, err
	}
	var orders []T
	if err := res.All(ctx, &orders); err != nil {
		return nil, err
	}
	logger.Get().Sugar().Debugf("cust: %s\ngot orders: %v\n", customerID.String(), orders)
	return orders, nil
}

func (o *orderRepo[T]) findOneAndUpdate(ctx context.Context, filter, update any) (T, error) {
	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(options.After)
	res := o.orders.FindOneAndUpdate(ctx, filter, update, opts)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return *new(T), domain.ErrOrderNotFound
		}
		return *new(T), res.Err()
	}
	var ord T
	if err := res.Decode(&ord); err != nil {
		return *new(T), err
	}
	return ord, nil
}
