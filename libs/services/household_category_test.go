package services

import (
	"testing"

	"domain"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDoAllExist(t *testing.T) {

	t.Run("one element missing, expect false", func(t *testing.T) {
		var cartProducts []domain.HouseholdProduct
		for i := 0; i < 5; i++ {
			p := domain.HouseholdProduct{
				ProductID: primitive.NewObjectID(),
				Price:     uint32(i),
			}
			cartProducts = append(cartProducts, p)
		}
		// missing cartProduct[0]
		products := cartProducts[1:]
		ok, missingProduct := doAllExist(cartProducts, products)

		require.False(t, ok)
		require.EqualValues(t, cartProducts[0], missingProduct)
	})

	t.Run("2 elements are missing, expect false", func(t *testing.T) {
		var cartProducts []domain.HouseholdProduct
		for i := 0; i < 5; i++ {
			p := domain.HouseholdProduct{
				ProductID: primitive.NewObjectID(),
				Price:     uint32(i),
			}
			cartProducts = append(cartProducts, p)
		}

		// missing cartProduct[3], cartProduct[4]
		products := append(cartProducts[:3])
		ok, missingProduct := doAllExist(cartProducts, products)
		require.False(t, ok)
		require.EqualValues(t, cartProducts[3], missingProduct)
	})

	t.Run("all products and id's on place, expect true", func(t *testing.T) {
		var cartProducts []domain.HouseholdProduct
		for i := 0; i < 5; i++ {
			p := domain.HouseholdProduct{
				ProductID: primitive.NewObjectID(),
				Price:     uint32(i),
			}
			cartProducts = append(cartProducts, p)
		}

		products := cartProducts
		ok, missingProduct := doAllExist(cartProducts, products)
		require.True(t, ok)
		require.Zero(t, missingProduct)
	})
}
