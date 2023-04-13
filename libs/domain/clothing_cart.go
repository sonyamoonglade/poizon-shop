package domain

type ClothingCart []ClothingPosition

func (c *ClothingCart) Add(p ClothingPosition) {
	*c = append(*c, p)
}

func (c *ClothingCart) Remove(positionID string) {
	for i, p := range *c {
		if p.PositionID.Hex() == positionID {
			// swap to end and cut
			c.swap(i, len(*c)-1)
			*c = (*c)[:len(*c)-1]
			break
		}
	}
}

func (c *ClothingCart) RemoveAt(index int) {
	for i := range *c {
		if i == index {
			// swap to end and slice
			c.swap(i, len(*c)-1)
			*c = (*c)[:len(*c)-1]
			break
		}
	}
}

func (c *ClothingCart) swap(i, j int) {
	(*c)[i], (*c)[j] = (*c)[j], (*c)[i]
}
