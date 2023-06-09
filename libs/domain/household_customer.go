package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type HouseholdCustomer struct {
	CustomerID  primitive.ObjectID  `json:"customerId" bson:"_id,omitempty"`
	TelegramID  int64               `json:"telegramID" bson:"telegramId"`
	Username    *string             `json:"username,omitempty" bson:"username,omitempty"`
	FullName    *string             `json:"fullName,omitempty" bson:"fullName,omitempty"`
	PhoneNumber *string             `json:"phoneNumber,omitempty" bson:"phoneNumber,omitempty"`
	State       State               `json:"state" bson:"state"`
	Cart        HouseholdCart       `json:"cart" bson:"cart"`
	PromocodeID *primitive.ObjectID `json:"promocodeId,omitempty" bson:"promocodeId,omitempty"`

	// Used to join promocode by PromocodeID
	Promocode *Promocode `json:"promocode,omitempty" bson:"-"`
}

func NewHouseholdCustomer(telegramID int64, username string) HouseholdCustomer {
	return HouseholdCustomer{
		CustomerID: primitive.NewObjectID(),
		TelegramID: telegramID,
		Username:   &username,
		State:      StateDefault,
		Cart:       NewHouseholdCart(),
	}
}

func (c *HouseholdCustomer) UpdateState(newState State) *HouseholdCustomer {
	c.State = newState
	return c
}

func (c *HouseholdCustomer) SetFullName(fullName string) *HouseholdCustomer {
	c.FullName = &fullName
	return c
}

func (c *HouseholdCustomer) SetPhoneNumber(phoneNumber string) *HouseholdCustomer {
	c.PhoneNumber = &phoneNumber
	return c
}

func (c *HouseholdCustomer) GetTgState() uint8 {
	return c.State.V
}

// SetPromocode Should only be used to join Promocode field
func (c *HouseholdCustomer) SetPromocode(p Promocode) *HouseholdCustomer {
	c.Promocode = &p
	return c
}

func (c *HouseholdCustomer) UsePromocode(p Promocode) *HouseholdCustomer {
	c.Promocode = &p
	c.PromocodeID = &p.PromocodeID
	return c
}

func (c *HouseholdCustomer) HasPromocode() bool {
	return c.PromocodeID != nil
}

func (c *HouseholdCustomer) GetPromocode() (Promocode, bool) {
	if c.HasPromocode() {
		return *c.Promocode, true
	}
	return Promocode{}, false
}

func (c *HouseholdCustomer) MustGetPromocode() Promocode {
	if c.HasPromocode() {
		return *c.Promocode
	}
	return Promocode{}
}
