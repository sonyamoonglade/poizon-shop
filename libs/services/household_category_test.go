package services

import (
	"domain"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDoAllExist(t *testing.T) {

	t.Run("one element missing, expect false", func(t *testing.T) {
		var ids []primitive.ObjectID
		for i := 0; i < 5; i++ {
			ids = append(ids, primitive.NewObjectID())
		}
		var products []domain.HouseholdProduct
		for i := 0; i < 5; i++ {
			// Skip 3rd
			if i == 3 {
				continue
			}
			p := domain.HouseholdProduct{
				ProductID: ids[i],
			}
			products = append(products, p)
		}
		ok := doAllExist(ids, products)
		require.False(t, ok)
	})

	t.Run("first element missing, expect false", func(t *testing.T) {
		var ids []primitive.ObjectID
		for i := 0; i < 5; i++ {
			ids = append(ids, primitive.NewObjectID())
		}
		var products []domain.HouseholdProduct
		for i := 0; i < 5; i++ {
			// Skip 3rd
			if i == 0 {
				continue
			}
			p := domain.HouseholdProduct{
				ProductID: ids[i],
			}
			products = append(products, p)
		}

		ok := doAllExist(ids, products)
		require.False(t, ok)
	})

	t.Run("all products and id's on place, expect true", func(t *testing.T) {
		var ids []primitive.ObjectID
		for i := 0; i < 5; i++ {
			ids = append(ids, primitive.NewObjectID())
		}
		var products []domain.HouseholdProduct
		for i := 0; i < 5; i++ {
			p := domain.HouseholdProduct{
				ProductID: ids[i],
			}
			products = append(products, p)
		}

		ok := doAllExist(ids, products)
		require.True(t, ok)
	})
}
