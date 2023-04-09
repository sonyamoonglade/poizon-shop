package repositories

import "onlineshop/database"

type Repositories struct {
	ClothingCustomer  *clothingCustomerRepo
	ClothingOrder     *clothingOrderRepo
	ClothingCatalog   *clothingCatalogRepo
	HouseholdCategory *householdCategoryRepo
	Rate              *rateRepo
}

const (
	customers           = "customers"
	orders              = "orders"
	catalog             = "catalog"
	rateKeyValue        = "rate_keyValue"
	householdCategories = "household_categories"
)

func NewRepositories(db *database.Mongo, clothingOnChange ClothingOnChangeFunc, householdOnChange HouseholdOnChangeFunc) Repositories {
	return Repositories{
		ClothingCustomer:  NewClothingCustomerRepo(db.Collection(customers)),
		ClothingOrder:     NewClothingOrderRepo(db.Collection(orders)),
		ClothingCatalog:   NewClothingCatalogRepo(db.Collection(catalog), clothingOnChange),
		Rate:              NewRateRepository(db.Collection(rateKeyValue)),
		HouseholdCategory: NewHouseholdCategoryRepo(db.Collection(householdCategories), householdOnChange),
	}
}
