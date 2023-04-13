package repositories

import (
	"domain"
	"onlineshop/database"
)

type Repositories struct {
	ClothingCustomer  *clothingCustomerRepo
	ClothingOrder     *orderRepo[domain.ClothingOrder]
	HouseholdOrder    *orderRepo[domain.HouseholdOrder]
	ClothingCatalog   *clothingCatalogRepo
	HouseholdCategory *householdCategoryRepo
	Rate              *rateRepo
}

const (
	customers    = "customers"
	customersHH  = "h-customers"
	orders       = "orders"
	ordersHH     = "h-orders"
	catalog      = "catalog"
	rateKeyValue = "rate_keyValue"
	//todo: rename to h-categories
	householdCategories = "household_categories"
)

func NewRepositories(db *database.Mongo, clothingOnChange ClothingOnChangeFunc, householdOnChange HouseholdOnChangeFunc) Repositories {
	return Repositories{
		ClothingCustomer:  NewClothingCustomerRepo(db.Collection(customers)),
		ClothingOrder:     NewClothingOrderRepo(db.Collection(orders)),
		HouseholdOrder:    NewHouseholdOrderRepo(db.Collection(ordersHH)),
		ClothingCatalog:   NewClothingCatalogRepo(db.Collection(catalog), clothingOnChange),
		Rate:              NewRateRepository(db.Collection(rateKeyValue)),
		HouseholdCategory: NewHouseholdCategoryRepo(db.Collection(householdCategories), householdOnChange),
	}
}
