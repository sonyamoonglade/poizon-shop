package services

import (
	"context"
	"domain"
	"dto"
	"repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type customerService[T customerConstraint, D dtoConstraint] struct {
	promocodeRepo         repositories.Promocode
	householdCustomerRepo repositories.HouseholdCustomer
	clothingCustomerRepo  repositories.ClothingCustomer
}

func NewHouseholdCustomerService(promocodeRepo repositories.Promocode, hhCustomerRepo repositories.HouseholdCustomer) *customerService[domain.HouseholdCustomer, dto.UpdateHouseholdCustomerDTO] {
	return &customerService[domain.HouseholdCustomer, dto.UpdateHouseholdCustomerDTO]{
		promocodeRepo:         promocodeRepo,
		householdCustomerRepo: hhCustomerRepo,
	}
}

func (c *customerService[T, D]) GetByTelegramID(ctx context.Context, telegramID int64) (T, error) {
	
}

func (c *customerService[T, D]) All(ctx context.Context) ([]T, error) {
	panic("not implemented") // TODO: Implement
}

func (c *customerService[T, D]) Save(ctx context.Context, c T) error {
	panic("not implemented") // TODO: Implement
}

func (c *customerService[T, D]) UpdateState(ctx context.Context, telegramID int64, newState domain.State) error {
	panic("not implemented") // TODO: Implement
}

func (c *customerService[T, D]) Update(ctx context.Context, customerID primitive.ObjectID, dto D) error {
	panic("not implemented") // TODO: Implement
}

func (c *customerService[T, D]) Delete(ctx context.Context, customerID primitive.ObjectID) error {
	panic("not implemented") // TODO: Implement
}

func (c *customerService[T, D]) isHousehold[T customerConstraint](customer T) (domain.HouseholdCustomer, bool) {
	v, ok := any(customer).(domain.HouseholdCustomer)
	return v, ok
}

func (c *customerService[T, D]) isClothing[T customerConstraint](customer T) bool {
	zero := new(T)
	_, ok := any(zero).(domain.ClothingCustomer)
	return v, ok
}
