package telegram

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrMessageAlreadyExists = errors.New("message already exists")

type CatalogMsg struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	MsgID  int                `bson:"msgId"`
	ChatID int64              `bson:"chatId"`
}
