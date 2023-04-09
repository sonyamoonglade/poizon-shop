package domain

import (
	"errors"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type State struct {
	V uint8 `json:"v" bson:"v"`
}

func (s State) Value() uint8 {
	return s.V
}

var (
	// Default state not waiting for any make order response
	StateDefault                       = State{0}
	StateWaitingForOrderType           = State{1}
	StateWaitingForCategory            = State{2}
	StateWaitingForCalculatorCategory  = State{3}
	StateWaitingForCalculatorOrderType = State{4}
	StateWaitingForSize                = State{5}
	StateWaitingForButton              = State{6}
	StateWaitingForPrice               = State{7}
	StateWaitingForLink                = State{8}
	StateWaitingForCartPositionToEdit  = State{9}
	StateWaitingForCalculatorInput     = State{10}
	StateWaitingForFIO                 = State{11}
	StateWaitingForPhoneNumber         = State{12}
	StateWaitingForDeliveryAddress     = State{13}
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrNoCustomers      = errors.New("no customers found")
)

type Meta struct {
	NextOrderType *OrderType `json:"nextOrderType" bson:"nextOrderType"`
}

type CalculatorMeta struct {
	NextOrderType *OrderType `json:"nextOrderType" bson:"nextOrderType"`
	Category      *Category  `bson:"category"`
}

type Customer struct {
	CustomerID       primitive.ObjectID `json:"customerId" json:"customerID,omitempty" bson:"_id,omitempty"`
	TelegramID       int64              `json:"telegramID" bson:"telegramId"`
	Username         *string            `json:"username,omitempty" bson:"username,omitempty"`
	FullName         *string            `json:"fullName,omitempty" bson:"fullName,omitempty"`
	PhoneNumber      *string            `json:"phoneNumber,omitempty" bson:"phoneNumber,omitempty"`
	TgState          State              `json:"state" bson:"state"`
	Cart             ClothingCart       `json:"cart" bson:"cart"`
	Meta             Meta               `json:"meta" bson:"meta"`
	CalculatorMeta   CalculatorMeta     `json:"calculatorMeta" bson:"calculatorMeta"`
	CatalogOffset    uint               `json:"catalogOffset" bson:"catalogOffset"`
	LastEditPosition *ClothingPosition  `json:"lastEditPosition,omitempty" bson:"lastEditPosition"`
}

func NewCustomer(telegramID int64, username string) Customer {
	return Customer{
		TelegramID: telegramID,
		Username:   &username,
		TgState:    StateDefault,
	}
}

func (c *Customer) UpdateState(newState State) {
	c.TgState = newState
}
func (c *Customer) GetTgState() uint8 {
	return c.TgState.V
}

func (c *Customer) SetLastEditPosition(p ClothingPosition) {
	c.LastEditPosition = &p
}

func (c *Customer) UpdateLastEditPositionSize(s string) {
	c.LastEditPosition.Size = s
}

func (c *Customer) UpdateLastEditPositionCategory(cat Category) {
	if c.LastEditPosition == nil {
		c.LastEditPosition = &ClothingPosition{}
	}
	c.LastEditPosition.Category = cat
}

func (c *Customer) UpdateLastEditPositionPrice(priceRub uint64, priceYuan uint64) {
	c.LastEditPosition.PriceRUB = priceRub
	c.LastEditPosition.PriceYUAN = priceYuan
}

func (c *Customer) UpdateLastEditPositionLink(link string) {
	c.LastEditPosition.ShopLink = link
}

func (c *Customer) UpdateLastEditPositionButtonColor(button Button) {
	c.LastEditPosition.Button = button
}

func (c *Customer) UpdateMetaOrderType(typ OrderType) {
	c.Meta.NextOrderType = &typ
}

func (c *Customer) UpdateCalculatorMetaCategory(cat Category) {
	c.CalculatorMeta.Category = &cat
}
func (c *Customer) UpdateCalculatorMetaOrderType(typ OrderType) {
	c.CalculatorMeta.NextOrderType = &typ
}

func (c *Customer) IncrementCatalogOffset() {
	c.CatalogOffset++
}

func (c *Customer) NullifyCatalogOffset() {
	c.CatalogOffset = 0
}

const defaultUsername = "User"

func MakeUsername(username string) string {
	if username == "" {
		return defaultUsername
	}
	return username
}

func IsValidFullName(fullName string) bool {
	spaceCount := strings.Count(fullName, " ")
	if spaceCount != 2 {
		return false
	}
	return true
}

var r = regexp.MustCompile(`^(8|7)((\d{10})|(\s\(\d{3}\)\s\d{3}\s\d{2}\s\d{2}))`)

func IsValidPhoneNumber(phoneNumber string) bool {
	if len(phoneNumber) != 11 {
		return false
	}
	return r.MatchString(phoneNumber)
}
