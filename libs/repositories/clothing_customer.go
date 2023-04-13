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

type clothingCustomerRepo struct {
	customers *mongo.Collection
}

func NewClothingCustomerRepo(db *mongo.Collection) *clothingCustomerRepo {
	return &clothingCustomerRepo{
		customers: db,
	}
}

func (c *clothingCustomerRepo) Save(ctx context.Context, customer domain.ClothingCustomer) error {
	if _, err := c.customers.InsertOne(ctx, customer); err != nil {
		return err
	}
	return nil
}

func (c *clothingCustomerRepo) Delete(ctx context.Context, customerID primitive.ObjectID) error {
	if _, err := c.customers.DeleteOne(ctx, bson.M{"_id": customerID}); err != nil {
		return err
	}
	return nil
}

func (c *clothingCustomerRepo) NullifyCatalogOffsets(ctx context.Context) error {
	filter := bson.D{}
	update := bson.M{"$set": bson.M{"catalogOffset": 0}}
	_, err := c.customers.UpdateMany(ctx, filter, update)
	return err
}

func (c *clothingCustomerRepo) Update(ctx context.Context, customerID primitive.ObjectID, dto dto.UpdateClothingCustomerDTO) error {

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
	if dto.LastPosition != nil {
		if dto.LastPosition.PositionID.IsZero() {
			dto.LastPosition.PositionID = primitive.NewObjectID()
		}
		update["lastEditPosition"] = *dto.LastPosition
	}

	if dto.Username != nil {
		update["username"] = *dto.Username
	}

	if dto.FullName != nil {
		update["fullName"] = *dto.FullName
	}

	if dto.Meta != nil {

		if dto.Meta.NextOrderType != nil {
			update["meta.nextOrderType"] = dto.Meta.NextOrderType
		}
	}
	if dto.CalculatorMeta != nil {
		if dto.CalculatorMeta.Category != nil {
			update["calculatorMeta.category"] = *dto.CalculatorMeta.Category
		}

		if dto.CalculatorMeta.NextOrderType != nil {
			update["calculatorMeta.nextOrderType"] = *dto.CalculatorMeta.NextOrderType
		}
	}

	if dto.CatalogOffset != nil {
		update["catalogOffset"] = *dto.CatalogOffset
	}

	_, err := c.customers.UpdateByID(ctx, customerID, bson.M{"$set": update})
	if err != nil {
		return err
	}
	return nil
}

func (c *clothingCustomerRepo) UpdateState(ctx context.Context, telegramID int64, newState domain.State) error {
	filter := bson.M{"telegramId": telegramID}
	updateQuery := bson.M{"$set": bson.M{"state": newState}}
	_, err := c.customers.UpdateOne(ctx, filter, updateQuery)
	return err
}

func (c *clothingCustomerRepo) GetByTelegramID(ctx context.Context, telegramID int64) (domain.ClothingCustomer, error) {
	query := bson.M{"telegramId": telegramID}
	res := c.customers.FindOne(ctx, query)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.ClothingCustomer{}, domain.ErrCustomerNotFound
		}
		return domain.ClothingCustomer{}, err
	}
	var customer domain.ClothingCustomer
	if err := res.Decode(&customer); err != nil {
		return domain.ClothingCustomer{}, fmt.Errorf("cant decode customer: %w", err)
	}
	return customer, nil
}
func (c *clothingCustomerRepo) GetState(ctx context.Context, telegramID int64) (domain.State, error) {
	customer, err := c.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return domain.StateDefault, err
	}
	return customer.TgState, nil
}

func (c *clothingCustomerRepo) All(ctx context.Context) ([]domain.ClothingCustomer, error) {
	res, err := c.customers.Find(ctx, bson.D{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNoCustomers
		}
		return nil, err
	}
	var customers []domain.ClothingCustomer
	if err := res.All(ctx, &customers); err != nil {
		return nil, err
	}
	return customers, nil
}
