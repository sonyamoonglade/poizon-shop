package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HouseholdCart []HouseholdProduct

func NewHouseholdCart() HouseholdCart {
	return HouseholdCart(nil)
}

type GroupedHouesholdProducts struct {
	// Product
	P HouseholdProduct
	// Quantity
	Qty int
}

func (c *HouseholdCart) Group() []GroupedHouesholdProducts {
	groups := make(map[primitive.ObjectID]GroupedHouesholdProducts, c.Size())
	for _, product := range c.Slice() {
		if _, ok := groups[product.ProductID]; !ok {
			groups[product.ProductID] = GroupedHouesholdProducts{
				P:   product,
				Qty: 1,
			}
		} else {
			groups[product.ProductID] = GroupedHouesholdProducts{
				P:   groups[product.ProductID].P,
				Qty: groups[product.ProductID].Qty + 1,
			}
		}
	}
	out := make([]GroupedHouesholdProducts, 0, len(groups))
	for _, v := range groups {
		out = append(out, v)
	}
	return out
}

func (c *HouseholdCart) First() (HouseholdProduct, bool) {
	if len(*c) == 0 {
		return HouseholdProduct{}, false
	}
	return (*c)[0], true
}

func (c *HouseholdCart) Size() int {
	return len(*c)
}

func (c *HouseholdCart) Clear() {
	*c = nil
}

func (c *HouseholdCart) IsEmpty() bool {
	return len(*c) == 0
}

func (c *HouseholdCart) Add(p HouseholdProduct) {
	*c = append(*c, p)
}

func (c *HouseholdCart) RemoveAt(index int) {
	for i := range *c {
		if i == index {
			// swap to end and slice
			c.swap(i, len(*c)-1)
			*c = (*c)[:len(*c)-1]
			break
		}
	}
}

var RemoveByProductID = func(productIDStr string) RemovePredicate {
	id, _ := primitive.ObjectIDFromHex(productIDStr)
	return func(p HouseholdProduct, i int) bool {
		return id == p.ProductID
	}
}

type RemovePredicate func(p HouseholdProduct, i int) bool

func (c *HouseholdCart) Remove(predicate RemovePredicate) HouseholdProduct {
	for i := range *c {
		p := (*c)[i]
		if predicate(p, i) {
			c.swap(i, len(*c)-1)
			*c = (*c)[:len(*c)-1]
			return p
		}
	}
	return HouseholdProduct{}
}

func (c *HouseholdCart) swap(i, j int) {
	(*c)[i], (*c)[j] = (*c)[j], (*c)[i]
}

func (c *HouseholdCart) Slice() []HouseholdProduct {
	return *c
}
