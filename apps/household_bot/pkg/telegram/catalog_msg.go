package telegram

import "go.mongodb.org/mongo-driver/bson/primitive"

type CatalogMsg struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	MsgID  int                `bson:"msgId"`
	ChatID int64              `bson:"chatId"`
}
