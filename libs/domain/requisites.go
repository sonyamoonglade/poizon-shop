package domain

type Requisites struct {
	SberID    string `bson:"sberId"`
	TinkoffID string `bson:"tinkoffId"`
}

var AdminRequisites = Requisites{
	SberID:    "2202 2062 2769 5751",
	TinkoffID: "2200 7007 7461 0942",
}
