package domain

type HouseholdCart []HouseholdProduct

func (c *HouseholdCart) Add(p HouseholdProduct) {
	*c = append(*c, p)
}

func (c *HouseholdCart) Remove(productID string) {
	for i, p := range *c {
		if p.ProductID.Hex() == productID {
			// swap to end and slice
			c.swap(i, len(*c)-1)
			*c = (*c)[:len(*c)-1]
			break
		}
	}
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

func (c *HouseholdCart) swap(i, j int) {
	(*c)[i], (*c)[j] = (*c)[j], (*c)[i]
}
