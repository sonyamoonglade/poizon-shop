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

func TestHasAt(t *testing.T) {

	p := NewProvider()
	items := []domain.HouseholdCategory{
		{
			Title: "category1",
			Subcategories: []domain.Subcategory{
				{Title: "subcategory1",
					Products: []domain.HouseholdProduct{
						{Name: "1"},
						{Name: "2"},
						{Name: "3"},
					}},
			},
			Rank: 0,
		},
	}

	p.Load(items)

	require.True(t, p.HasAt("category1", "subcategory1", 0))
	require.True(t, p.HasAt("category1", "subcategory1", 1))
	require.True(t, p.HasAt("category1", "subcategory1", 2))
	require.False(t, p.HasAt("category1", "subcategory1", 3))
}
func TestGetProductAt(t *testing.T) {
	p := NewProvider()
	items := []domain.HouseholdCategory{
		{
			Title: "category1",
			Subcategories: []domain.Subcategory{
				{Title: "subcategory1",
					Products: []domain.HouseholdProduct{
						{Name: "0"},
						{Name: "1"},
						{Name: "2"},
					}},
			},
			Rank: 0,
		},
	}

	p.Load(items)
	cTitle, sTitle := "category1", "subcategory1"
	p1, ok := p.GetProductAt(cTitle, sTitle, 0)
	require.True(t, ok)
	require.EqualValues(t, items[0].Subcategories[0].Products[0], p1)
	p2, ok := p.GetProductAt(cTitle, sTitle, 1)
	require.EqualValues(t, items[0].Subcategories[0].Products[1], p2)
	require.True(t, ok)
	p3, ok := p.GetProductAt(cTitle, sTitle, 2)
	require.EqualValues(t, items[0].Subcategories[0].Products[2], p3)
	require.True(t, ok)
}
