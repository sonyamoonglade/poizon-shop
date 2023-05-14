package usecase

import (
	"context"
	"fmt"

	"domain"
	"dto"
	"repositories"
	"services"
)

type HouseholdMakeOrder struct {
	promocodeRepo   repositories.Promocode
	orderService    services.Order[domain.HouseholdOrder]
	customerService services.HouseholdCustomer
}

func NewHouseholdMakeOrderUsecase(
	promocodeRepo repositories.Promocode,
	orderService services.Order[domain.HouseholdOrder],
	customerService services.HouseholdCustomer,
) *HouseholdMakeOrder {
	return &HouseholdMakeOrder{
		promocodeRepo:   promocodeRepo,
		orderService:    orderService,
		customerService: customerService,
	}
}

func (h *HouseholdMakeOrder) NewOrder(ctx context.Context, deliveryAddress string, customer domain.HouseholdCustomer,
	inStock bool) (domain.HouseholdOrder, error) {
	shortID, err := h.orderService.GetFreeShortID(ctx)
	if err != nil {
		return domain.HouseholdOrder{}, fmt.Errorf("orderService.GetFreeShortID: %w", err)
	}

	order := domain.NewHouseholdOrder(customer, deliveryAddress, shortID, inStock)
	if customer.HasPromocode() {
		promo := customer.MustGetPromocode()
		order.UseDiscount(promo.GetHouseholdDiscount())
	}

	if err := h.orderService.Save(ctx, order); err != nil {
		return domain.HouseholdOrder{}, fmt.Errorf("orderService.Save: %w", err)
	}

	customer.Cart.Clear()
	updateDTO := dto.UpdateHouseholdCustomerDTO{
		Cart:  &customer.Cart,
		State: &domain.StateDefault,
	}
	if err := h.customerService.Update(ctx, customer.CustomerID, updateDTO); err != nil {
		return domain.HouseholdOrder{}, fmt.Errorf("customerService.Update: %w", err)
	}

	return order, nil
}
