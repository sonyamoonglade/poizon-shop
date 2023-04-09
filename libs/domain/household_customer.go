package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type HouseholdCustomer struct {
	CustomerID  primitive.ObjectID `json:"customerId" json:"customerID,omitempty" bson:"_id,omitempty"`
	TelegramID  int64              `json:"telegramID" bson:"telegramId"`
	Username    *string            `json:"username,omitempty" bson:"username,omitempty"`
	FullName    *string            `json:"fullName,omitempty" bson:"fullName,omitempty"`
	PhoneNumber *string            `json:"phoneNumber,omitempty" bson:"phoneNumber,omitempty"`
	TgState     State              `json:"state" bson:"state"`
	Cart        HouseholdCart      `json:"cart" bson:"cart"`
	Meta        Meta               `json:"meta" bson:"meta"`
}
