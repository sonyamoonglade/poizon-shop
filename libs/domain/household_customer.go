package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type HouseholdCustomer struct {
	CustomerID  primitive.ObjectID `json:"customerId" bson:"_id,omitempty"`
	TelegramID  int64              `json:"telegramID" bson:"telegramId"`
	Username    *string            `json:"username,omitempty" bson:"username,omitempty"`
	FullName    *string            `json:"fullName,omitempty" bson:"fullName,omitempty"`
	PhoneNumber *string            `json:"phoneNumber,omitempty" bson:"phoneNumber,omitempty"`
	State       State              `json:"state" bson:"state"`
	Cart        HouseholdCart      `json:"cart" bson:"cart"`
}

func NewHouseholdCustomer(telegramID int64, username string) HouseholdCustomer {
	return HouseholdCustomer{
		TelegramID: telegramID,
		Username:   &username,
		State:      StateDefault,
	}
}

func (c *HouseholdCustomer) UpdateState(newState State) {
	c.State = newState
}
func (c *HouseholdCustomer) GetTgState() uint8 {
	return c.State.V
}
