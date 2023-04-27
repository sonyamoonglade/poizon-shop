package repositories

import (
	"domain"
	"onlineshop/database"
)

const (
	customersHH  = "h-customers"
	categoriesHH = "h-categories"
	ordersHH     = "h-orders"
	catalogMsgHH = "h-catalog_msg"

	orders    = "orders"
	customers = "customers"
	catalog   = "catalog"

	rateKeyValue = "rate_keyValue"
	promocodes   = "promocodes"
)

type Repositories struct {
	HouseholdCustomer   *householdCustomerRepo
	HouseholdOrder      *orderRepo[domain.HouseholdOrder]
	HouseholdCategory   *householdCategoryRepo
	HouseholdCatalogMsg *householdCatalogMsgRepo

	ClothingCustomer *clothingCustomerRepo
	ClothingOrder    *orderRepo[domain.ClothingOrder]
	ClothingCatalog  *clothingCatalogRepo

	Rate      *rateRepo
	Promocode *promocodeRepo
}

func NewRepositories(db *database.Mongo, clothingOnChange ClothingOnChangeFunc, householdOnChange HouseholdOnChangeFunc) Repositories {
	return Repositories{
		ClothingCustomer:    NewClothingCustomerRepo(db.Collection(customers)),
		HouseholdCustomer:   NewHouseholdCustomerRepo(db.Collection(customersHH)),
		ClothingOrder:       NewClothingOrderRepo(db.Collection(orders)),
		HouseholdOrder:      NewHouseholdOrderRepo(db.Collection(ordersHH)),
		ClothingCatalog:     NewClothingCatalogRepo(db.Collection(catalog), clothingOnChange),
		Rate:                NewRateRepository(db.Collection(rateKeyValue)),
		HouseholdCategory:   NewHouseholdCategoryRepo(db.Collection(categoriesHH), householdOnChange),
		HouseholdCatalogMsg: NewHouseholdCatalogMsgRepo(db.Collection(catalogMsgHH)),
		Promocode:           NewPromocodeRepo(db.Collection(promocodes)),
	}
}
