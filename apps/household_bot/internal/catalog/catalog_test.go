package catalog

import (
	"testing"

	"domain"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	p := NewProvider()
	items := []domain.HouseholdCategory{
		{
			Title: "category1",
			Subcategories: []domain.Subcategory{
				{Title: "subcategory1"},
				{Title: "subcategory2"},
			},
			Rank: 0,
		},
		{
			Title: "category2",
			Subcategories: []domain.Subcategory{
				{Title: "subcategory3"},
			},
			Rank: 1,
		},
	}

	p.Load(items)

	require.Equal(t, 2, len(p.items))
	c, ok := p.GetCategory(items[0].Title)
	require.True(t, ok)
	require.Equal(t, items[0], c)
	c, ok = p.GetCategory(items[1].Title)
	require.True(t, ok)
	require.Equal(t, items[1], c)
}
func TestGetCategory(t *testing.T) {
	p := NewProvider()
	items := []domain.HouseholdCategory{
		{
			Title:  "category1",
			Active: true,
			Subcategories: []domain.Subcategory{
				{Title: "subcategory1", Active: true},
				{Title: "subcategory2", Active: true},
			},
		},
	}

	p.Load(items)

	// test existing category
	c, ok := p.GetCategory("category1")
	require.True(t, ok)
	require.Equal(t, "category1", c.Title)

	// test non-existing category
	_, ok = p.GetCategory("non-existing")
	require.False(t, ok)

	// test active filter
	titles := p.GetCategoryTitles(true)
	require.Equal(t, []string{"category1"}, titles)
}

func TestGetSubcategory(t *testing.T) {
	p := NewProvider()
	items := []domain.HouseholdCategory{
		{
			Title: "category1",
			Subcategories: []domain.Subcategory{
				{Title: "subcategory1", Active: true},
				{Title: "subcategory2", Active: true},
			},
		},
	}

	p.Load(items)

	// test existing subcategory
	s, ok := p.GetSubcategory("category1", "subcategory1")
	require.True(t, ok)
	require.Equal(t, items[0].Subcategories[0], s)

	// test non-existing subcategory
	_, ok = p.GetSubcategory("category1", "non-existing")
	require.False(t, ok)

	// test active filter
	titles := p.GetSubcategoryTitles("category1", true)
	require.EqualValues(t, []string{"subcategory1", "subcategory2"}, titles)
}
