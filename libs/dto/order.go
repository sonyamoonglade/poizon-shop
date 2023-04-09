package dto

import (
	"domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddCommentDTO struct {
	OrderID primitive.ObjectID
	Comment string
}

type ChangeOrderStatusDTO struct {
	OrderID   primitive.ObjectID
	NewStatus domain.Status
}
