package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CalculatorMeta struct {
	NextOrderType *OrderType `json:"nextOrderType" bson:"nextOrderType"`
	Category      *Category  `bson:"category"`
}

type ClothingCustomer struct {
	CustomerID       primitive.ObjectID  `json:"customerId" json:"customerID,omitempty" bson:"_id,omitempty"`
	TelegramID       int64               `json:"telegramID" bson:"telegramId"`
	Username         *string             `json:"username,omitempty" bson:"username,omitempty"`
	FullName         *string             `json:"fullName,omitempty" bson:"fullName,omitempty"`
	PhoneNumber      *string             `json:"phoneNumber,omitempty" bson:"phoneNumber,omitempty"`
	CatalogOffset    uint                `json:"catalogOffset" bson:"catalogOffset"`
	TgState          State               `json:"state" bson:"state"`
	Cart             ClothingCart        `json:"cart" bson:"cart"`
	Meta             Meta                `json:"meta" bson:"meta"`
	CalculatorMeta   CalculatorMeta      `json:"calculatorMeta" bson:"calculatorMeta"`
	LastEditPosition *ClothingPosition   `json:"lastEditPosition,omitempty" bson:"lastEditPosition"`
	PromocodeID      *primitive.ObjectID `json:"promocodeId,omitempty" bson:"promocodeId,omitempty"`

	// Used to join promocode by PromocodeID
	Promocode *Promocode `json:"promocode,omitempty" bson:"-"`
}

func NewClothingCustomer(telegramID int64, username string) ClothingCustomer {
	return ClothingCustomer{
		TelegramID: telegramID,
		Username:   &username,
		TgState:    StateDefault,
	}
}

func (c *ClothingCustomer) UpdateState(newState State) {
	c.TgState = newState
}
func (c *ClothingCustomer) GetTgState() uint8 {
	return c.TgState.V
}
func (c *ClothingCustomer) SetLastEditPosition(p ClothingPosition) { c.LastEditPosition = &p }

func (c *ClothingCustomer) UpdateLastEditPositionSize(s string) {
	c.LastEditPosition.Size = s
}

func (c *ClothingCustomer) UpdateLastEditPositionCategory(cat Category) {
	if c.LastEditPosition == nil {
		c.LastEditPosition = &ClothingPosition{}
	}
	c.LastEditPosition.Category = cat
}

func (c *ClothingCustomer) UpdateLastEditPositionPrice(priceRub uint64, priceYuan uint64) {
	c.LastEditPosition.PriceRUB = priceRub
	c.LastEditPosition.PriceYUAN = priceYuan
}

func (c *ClothingCustomer) UpdateLastEditPositionLink(link string) {
	c.LastEditPosition.ShopLink = link
}

func (c *ClothingCustomer) UpdateLastEditPositionButtonColor(button Button) {
	c.LastEditPosition.Button = button
}

func (c *ClothingCustomer) UpdateMetaOrderType(typ OrderType) {
	c.Meta.NextOrderType = &typ
}

func (c *ClothingCustomer) UpdateCalculatorMetaCategory(cat Category) {
	c.CalculatorMeta.Category = &cat
}
func (c *ClothingCustomer) UpdateCalculatorMetaOrderType(typ OrderType) {
	c.CalculatorMeta.NextOrderType = &typ
}

func (c *ClothingCustomer) IncrementCatalogOffset() {
	c.CatalogOffset++
}

func (c *ClothingCustomer) NullifyCatalogOffset() {
	c.CatalogOffset = 0
}

func (c *ClothingCustomer) SetPromocode(p Promocode) ClothingCustomer {
	// Joined field
	c.Promocode = &p
	return *c
}
func (c *ClothingCustomer) UsePromocode(p Promocode) ClothingCustomer {
	c.Promocode = &p
	c.PromocodeID = &p.PromocodeID
	return *c
}

func (c *ClothingCustomer) HasPromocode() bool {
	return c.PromocodeID != nil
}

func (c *ClothingCustomer) GetPromocode() (*Promocode, bool) {
	if c.HasPromocode() {
		return c.Promocode, true
	}
	return nil, false
}
