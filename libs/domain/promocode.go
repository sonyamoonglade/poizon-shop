package domain

import (
	"encoding/json"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrNoPromocode  = errors.New("promocode not found")
	ErrNoPromocodes = errors.New("promocodes not found")
)

type Promocode struct {
	PromocodeID primitive.ObjectID `json:"promocodeId" bson:"_id,omitempty"`
	ShortID     string             `json:"shortId" bson:"shortId"`
	Description string             `json:"description" bson:"description"`
	Discounts   DiscountMap        `json:"discounts" bson:"discounts"`
	Counters    PromoCounters      `json:"counters" bson:"counters"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
}

func NewPromocode(description string, discounts DiscountMap, shortID string) Promocode {
	return Promocode{
		PromocodeID: primitive.NewObjectID(),
		Description: description,
		ShortID:     shortID,
		Discounts:   discounts,
		CreatedAt:   time.Now().UTC(),
	}
}

func (p Promocode) getDiscount(source Source) uint32 {
	if p.Discounts == nil {
		return 0
	}
	return p.Discounts[source]
}

func (p Promocode) GetHouseholdDiscount() uint32 {
	return p.getDiscount(SourceHousehold)
}

func (p Promocode) GetClothingDiscount() uint32 {
	return p.getDiscount(SourceClothing)
}

func (p *Promocode) IncrementAsFirst(x int) *Promocode {
	p.Counters.AsFirst += uint32(x)
	return p
}

func (p *Promocode) IncrementAsSecondEtc(x int) *Promocode {
	p.Counters.AsSecondEtc += uint32(x)
	return p
}

type PromoCounters struct {
	// Incremented when promocode is used to create first order for customer
	//
	// Incremented only once per customer
	AsFirst uint32 `json:"asFirst" bson:"asFirst"`

	// Incremented when promocode is used to create second,third etc... orders for customer
	AsSecondEtc uint32 `json:"asSecondEtc" bson:"asSecondEtc"`
}

type DiscountMap map[Source]uint32

func (d DiscountMap) MarshalBSON() ([]byte, error) {
	out := make(map[string]uint32)
	for k, v := range d {
		out[k.String()] = v
	}
	return bson.Marshal(out)
}
func (d *DiscountMap) UnmarshalBSON(raw []byte) error {
	*d = make(DiscountMap)
	recv := make(map[string]uint32)
	if err := bson.Unmarshal(raw, &recv); err != nil {
		return err
	}
	for k, v := range recv {
		(*d)[SourceFromString(k)] = v
	}
	return nil
}

func (d DiscountMap) MarshalJSON() ([]byte, error) {
	out := make(map[string]any)
	for k, v := range d {
		out[k.String()] = v
	}
	return json.Marshal(out)
}
func (d *DiscountMap) UnmarshalJSON(raw []byte) error {
	*d = make(DiscountMap)
	recv := make(map[string]any)
	if err := json.Unmarshal(raw, &recv); err != nil {
		return err
	}
	for k, v := range recv {
		(*d)[SourceFromString(k)] = v.(uint32)
	}
	return nil
}
