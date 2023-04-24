package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSize(t *testing.T) {
	cart := NewHouseholdCart()
	cart.Add(HouseholdProduct{
		ProductID: primitive.NewObjectID(),
	})
	cart.Add(HouseholdProduct{
		ProductID: primitive.NewObjectID(),
	})
	cart.Add(HouseholdProduct{
		ProductID: primitive.NewObjectID(),
	})

	require.Equal(t, 3, cart.Size())
}
