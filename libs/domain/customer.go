package domain

import (
	"errors"
	"regexp"
	"strings"
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
	StateWaitingToAddToCart            = State{14}
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrNoCustomers      = errors.New("no customers found")
)

type Meta struct {
	NextOrderType *OrderType `json:"nextOrderType" bson:"nextOrderType"`
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
