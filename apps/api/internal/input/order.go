package input

import (
	"domain"
	"dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddCommentToOrderInput struct {
	Comment string `json:"comment"`
}

func (a AddCommentToOrderInput) ToDTO(orderID primitive.ObjectID) dto.AddCommentDTO {
	return dto.AddCommentDTO{
		OrderID: orderID,
		Comment: a.Comment,
	}
}

type ChangeOrderStatusInput struct {
	NewStatus int `json:"newStatus"`
}

func (c ChangeOrderStatusInput) ToDTO(orderID primitive.ObjectID) dto.ChangeOrderStatusDTO {
	return dto.ChangeOrderStatusDTO{
		OrderID:   orderID,
		NewStatus: domain.Status(c.NewStatus),
	}
}
