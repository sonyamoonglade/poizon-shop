package services

import (
	"domain"
	"dto"
	"onlineshop/database"
	"repositories"
)

type Services struct {
	HouseholdCustomer   *customerService[domain.HouseholdCustomer, dto.UpdateHouseholdCustomerDTO]
	HouseholdOrder      *orderService[domain.HouseholdOrder]
	HouseholdCategory   *householdCategoryService
	HouseholdCatalogMsg *householdCatalogMsgService

	ClothingCustomer *customerService[domain.ClothingCustomer, dto.UpdateClothingCustomerDTO]
	ClothingOrder    *orderService[domain.ClothingOrder]
}

func NewServices(repos repositories.Repositories, transactor database.Transactor) Services {
	return Services{
		HouseholdCustomer:   NewHouseholdCustomerService(repos.Promocode, repos.HouseholdCustomer),
		HouseholdOrder:      NewHouseholdOrderService(repos.HouseholdOrder),
		HouseholdCategory:   NewHouseholdCategoryService(repos.HouseholdCategory, transactor),
		HouseholdCatalogMsg: NewHouseholdCatalogMsgService(repos.HouseholdCatalogMsg, transactor),
		ClothingCustomer:    NewClothingCustomerService(repos.Promocode, repos.ClothingCustomer),
		ClothingOrder:       NewClothingOrderService(repos.ClothingOrder),
	}
}
