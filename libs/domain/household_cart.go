package domain

type HouseholdCart []HouseholdProduct

func (c *HouseholdCart) First() (HouseholdProduct, bool) {
	if len(*c) == 0 {
		return HouseholdProduct{}, false
	}
	return (*c)[0], true
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

func (c *HouseholdCart) swap(i, j int) {
	(*c)[i], (*c)[j] = (*c)[j], (*c)[i]
}
