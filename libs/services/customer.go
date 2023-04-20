package services

import (
	"context"

	"domain"
	"dto"
	"repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: generic repos
type customerService[T customerConstraint, D dtoConstraint] struct {
	promocodeRepo         repositories.Promocode
	householdCustomerRepo repositories.HouseholdCustomer
	clothingCustomerRepo  repositories.ClothingCustomer
}

func NewHouseholdCustomerService(
	promocodeRepo repositories.Promocode,
	hhCustomerRepo repositories.HouseholdCustomer,
) *customerService[domain.HouseholdCustomer, dto.UpdateHouseholdCustomerDTO] {
	return &customerService[domain.HouseholdCustomer, dto.UpdateHouseholdCustomerDTO]{
		promocodeRepo:         promocodeRepo,
		householdCustomerRepo: hhCustomerRepo,
	}
}

func NewClothingCustomerService(
	promocodeRepo repositories.Promocode,
	clothingCustomerRepo repositories.ClothingCustomer,
) *customerService[domain.ClothingCustomer, dto.UpdateClothingCustomerDTO] {
	return &customerService[domain.ClothingCustomer, dto.UpdateClothingCustomerDTO]{
		promocodeRepo:        promocodeRepo,
		clothingCustomerRepo: clothingCustomerRepo,
	}
}

func (c *customerService[T, D]) GetByTelegramID(ctx context.Context, telegramID int64) (T, error) {
	if c.isHousehold() {
		customer, err := c.householdCustomerRepo.GetByTelegramID(ctx, telegramID)
		if err != nil {
			return *new(T), err
		}
		if customer.HasPromocode() {
			promo, err := c.promocodeRepo.GetByID(ctx, *customer.PromocodeID)
			if err != nil {
				return *new(T), err
			}
			customer.SetPromocode(promo)
		}

		return any(customer).(T), nil
	} else {
		customer, err := c.clothingCustomerRepo.GetByTelegramID(ctx, telegramID)
		if err != nil {
			return *new(T), err
		}
		if customer.HasPromocode() {
			promo, err := c.promocodeRepo.GetByID(ctx, *customer.PromocodeID)
			if err != nil {
				return *new(T), err
			}
			customer.SetPromocode(promo)
		}
		return any(customer).(T), nil
	}
}

func (c *customerService[T, D]) All(ctx context.Context) ([]T, error) {
	if c.isHousehold() {
		customers, err := c.householdCustomerRepo.All(ctx)
		if err != nil {
			return nil, err
		}
		for i, customer := range customers {
			if customer.HasPromocode() {
				promo, err := c.promocodeRepo.GetByID(ctx, *customer.PromocodeID)
				if err != nil {
					return nil, err
				}
				customers[i].SetPromocode(promo)
			}
		}
		return any(customers).([]T), nil
	} else {
		customers, err := c.clothingCustomerRepo.All(ctx)
		if err != nil {
			return nil, err
		}
		for i, customer := range customers {
			if customer.HasPromocode() {
				promo, err := c.promocodeRepo.GetByID(ctx, *customer.PromocodeID)
				if err != nil {
					return nil, err
				}
				customers[i].SetPromocode(promo)
			}
		}
		return any(customers).([]T), nil
	}
}

func (c *customerService[T, D]) Save(ctx context.Context, customer T) error {
	if c.isHousehold() {
		return c.householdCustomerRepo.Save(ctx, any(customer).(domain.HouseholdCustomer))
	}
	return c.clothingCustomerRepo.Save(ctx, any(customer).(domain.ClothingCustomer))
}

func (c *customerService[T, D]) UpdateState(ctx context.Context, telegramID int64, newState domain.State) error {
	if c.isHousehold() {
		return c.householdCustomerRepo.UpdateState(ctx, telegramID, newState)
	}
	return c.clothingCustomerRepo.UpdateState(ctx, telegramID, newState)
}

func (c *customerService[T, D]) Update(ctx context.Context, customerID primitive.ObjectID, d D) error {
	if c.isHousehold() {
		return c.householdCustomerRepo.Update(ctx, customerID, any(d).(dto.UpdateHouseholdCustomerDTO))
	}
	return c.clothingCustomerRepo.Update(ctx, customerID, any(d).(dto.UpdateClothingCustomerDTO))
}

func (c *customerService[T, D]) Delete(ctx context.Context, customerID primitive.ObjectID) error {
	if c.isHousehold() {
		return c.householdCustomerRepo.Delete(ctx, customerID)
	}
	return c.clothingCustomerRepo.Delete(ctx, customerID)
}

func (c *customerService[T, D]) isHousehold() bool {
	zero := *new(T)
	_, ok := any(zero).(domain.HouseholdCustomer)
	return ok
}

func (c *customerService[T, D]) isClothing() bool {
	zero := *new(T)
	_, ok := any(zero).(domain.ClothingCustomer)
	return ok
}
