package domain

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

var ErrUnknownOrderType = errors.New("unknown order type")

type OrderType int

func NewOrderTypeFromString(s string) (OrderType, error) {
	parsed, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("strconv atoi: %w", err)
	}
	if parsed == int(OrderTypeExpress) || parsed == int(OrderTypeNormal) {
		return OrderType(parsed), nil
	}
	return 0, ErrUnknownOrderType
}

const (
	OrderTypeExpress OrderType = iota + 1
	OrderTypeNormal
)

func (o OrderType) String() string {
	return strconv.Itoa(int(o))
}

type Status int

const (
	StatusNotApproved Status = iota + 1
	StatusApproved
	StatusBuyout
	StatusTransferToPoison
	StatusSentFromPoison
	StatusGotToRussia
	StatusCheckTrack
	StatusGotToOrdererCity
)

var StatusTexts = map[Status]string{
	StatusNotApproved:      "Не подтвержден",
	StatusApproved:         "Подтвержден",
	StatusBuyout:           "Выкуплен",
	StatusTransferToPoison: "Передан на склад POIZON",
	StatusSentFromPoison:   "Отправлен со склада POIZON в Россию",
	StatusGotToRussia:      "Пришл на склад распределения",
	StatusCheckTrack:       "Трэк номер",
	StatusGotToOrdererCity: "Пришел в город назначения",
}

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrNoOrders      = errors.New("no orders")
)

type Source struct {
	V string
}

func (s Source) String() string {
	return s.V
}

var (
	SourceClothing  = Source{"Clothing"}
	SourceHousehold = Source{"Household"}
)

func IsValidOrderStatus(s Status) bool {
	_, ok := StatusTexts[s]
	return ok
}

type formula func(x uint64, rate float64) (rub uint64)

type FormulaMap = map[OrderType]map[Category]formula

const (
	othMul   = 0.5
	lightMul = 1.6
	heavyMul = 2.6
)

var formulas = FormulaMap{
	OrderTypeExpress: {
		CategoryOther: expressfn(othMul, 764),
		CategoryLight: expressfn(lightMul, 764),
		CategoryHeavy: expressfn(heavyMul, 764),
	},
	OrderTypeNormal: {
		CategoryOther: normalfn(othMul, 715),
		CategoryLight: normalfn(lightMul, 715),
		CategoryHeavy: normalfn(heavyMul, 715),
	},
}

func expressfn(kg_mul float64, fee float64) formula {
	return func(x uint64, rate float64) (rub uint64) {
		v := (float64(x)*rate)*1.09 + (170.0 * kg_mul * rate) + fee
		return uint64(math.Ceil(v))
	}
}

func normalfn(kg_mul float64, fee float64) formula {
	return func(x uint64, rate float64) (rub uint64) {
		v := (float64(x)*rate)*1.09 + (50.0 * kg_mul * rate) + fee
		return uint64(math.Ceil(v))
	}
}

type ConvertYuanArgs struct {
	X         uint64
	Rate      float64
	OrderType OrderType
	Category  Category
}

func ConvertYuan(args ConvertYuanArgs) (rub uint64) {
	return formulas[args.OrderType][args.Category](args.X, args.Rate)
}
