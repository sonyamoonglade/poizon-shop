package repositories

import (
	"context"
	"errors"
	"fmt"

	"domain"
	"dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type householdCustomerRepo struct {
	customers *mongo.Collection
}

func NewHouseholdCustomerRepo(customers *mongo.Collection) *householdCustomerRepo {
	return &householdCustomerRepo{
		customers: customers,
	}
}

func (h *householdCustomerRepo) GetByTelegramID(ctx context.Context, telegramID int64) (domain.HouseholdCustomer, error) {
	query := bson.M{"telegramId": telegramID}
	res := h.customers.FindOne(ctx, query)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.HouseholdCustomer{}, domain.ErrCustomerNotFound
		}
		return domain.HouseholdCustomer{}, err
	}
	var customer domain.HouseholdCustomer
	if err := res.Decode(&customer); err != nil {
		return domain.HouseholdCustomer{}, fmt.Errorf("cant decode customer: %w", err)
	}
	return customer, nil
}

func (h *householdCustomerRepo) All(ctx context.Context) ([]domain.HouseholdCustomer, error) {
	res, err := h.customers.Find(ctx, bson.D{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNoCustomers
		}
		return nil, err
	}
	var customers []domain.HouseholdCustomer
	if err := res.All(ctx, &customers); err != nil {
		return nil, err
	}
	return customers, nil

}

func (h *householdCustomerRepo) Save(ctx context.Context, customer domain.HouseholdCustomer) error {
	if _, err := h.customers.InsertOne(ctx, customer); err != nil {
		return err
	}
	return nil
}

func (h *householdCustomerRepo) UpdateState(ctx context.Context, telegramID int64, newState domain.State) error {
	filter := bson.M{"telegramId": telegramID}
	updateQuery := bson.M{"$set": bson.M{"state": newState}}
	_, err := h.customers.UpdateOne(ctx, filter, updateQuery)
	return err

}

func (h *householdCustomerRepo) Update(ctx context.Context, customerID primitive.ObjectID, dto dto.UpdateHouseholdCustomerDTO) error {
	update := bson.M{}
	if dto.Cart != nil {
		update["cart"] = *dto.Cart
	}

	if dto.PhoneNumber != nil {
		update["phoneNumber"] = *dto.PhoneNumber
	}

	if dto.State != nil {
		update["state"] = *dto.State
	}

	if dto.Username != nil {
		update["username"] = *dto.Username
	}

	if dto.FullName != nil {
		update["fullName"] = *dto.FullName
	}

	if dto.PromocodeID != nil {
		update["promocodeId"] = *dto.PromocodeID
	}

	_, err := h.customers.UpdateByID(ctx, customerID, bson.M{"$set": update})
	return err
}

func (h *householdCustomerRepo) Delete(ctx context.Context, customerID primitive.ObjectID) error {
	if _, err := h.customers.DeleteOne(ctx, bson.M{"_id": customerID}); err != nil {
		return err
	}
	return nil

}

func (h *householdCustomerRepo) GetAllByPromocodeID(ctx context.Context, promocodeID primitive.ObjectID) ([]domain.ClothingCustomer, error) {
	cur, err := h.customers.Find(ctx, bson.M{"promocodeId": promocodeID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNoCustomers
		}
		return nil, err
	}

	var customers []domain.ClothingCustomer
	return customers, cur.All(ctx, &customers)
}

func (h *householdCustomerRepo) GetState(ctx context.Context, telegramID int64) (domain.State, error) {
	customer, err := h.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return domain.StateDefault, err
	}
	return customer.State, nil
}
