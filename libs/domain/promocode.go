package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Promocode struct {
	PromocodeID primitive.ObjectID `json:"promocodeId" bson:"_id,omitempty"`
	Description string             `json:"description" bson:"description"`
	Discounts   map[Source]uint32  `json:"discounts" bson:"discounts"`
	Counters    PromoCounters      `json:"counters" bson:"counters"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
}

func NewPromocode(description string, discounts map[Source]uint32) Promocode {
	return Promocode{
		Description: description,
		Discounts:   discounts,
		CreatedAt:   time.Now().UTC(),
	}
}

func (p Promocode) GetDiscount(source Source) uint32 {
	return p.Discounts[source]
}

func (p *Promocode) IncrementAsFirst() Promocode {
	p.Counters.AsFirst++
	return *p
}

func (p *Promocode) IncrementAsSecondEtc() Promocode {
	p.Counters.AsSecondEtc++
	return *p
}

type PromoCounters struct {
	// Incremented when promocode is used to create first order for customer
	//
	// Incremented only once per customer
	AsFirst uint32 `json:"asFirst" bson:"asFirst"`

	// Incremented when promocode is used to create second,third etc... orders for customer
	AsSecondEtc uint32 `json:"asSecondEtc" bson:"asSecondEtc"`
}
