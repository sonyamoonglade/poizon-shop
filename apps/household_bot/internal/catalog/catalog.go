package catalog

import (
	"sync"

	"domain"
	"repositories"
)

type Provider struct {
	items []domain.HouseholdCategory
	mu    *sync.RWMutex
}

func NewProvider() *Provider {
	return &Provider{
		mu:    new(sync.RWMutex),
		items: nil,
	}
}

// Load categories and it's subcategories are expected to be already sorted
func (p *Provider) Load(items []domain.HouseholdCategory) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.items = make([]domain.HouseholdCategory, len(items), len(items))
	for i, it := range items {
		p.items[i] = it
	}
}

func (p *Provider) GetCategory(title string, inStock bool) (domain.HouseholdCategory, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, c := range p.items {
		if c.Title == title && c.InStock == inStock {
			return c, true
		}
	}
	return domain.HouseholdCategory{}, false
}

func (p *Provider) GetCategoryTitles(active bool, inStock bool) []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	titles := make([]string, 0, len(p.items))
	for _, c := range p.items {
		if c.Active == active && c.InStock == inStock {
			titles = append(titles, c.Title)
		}
	}
	return titles
}

func (p *Provider) GetSubcategory(cTitle, sTitle string, inStock bool) (domain.Subcategory, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, c := range p.items {
		if c.Title == cTitle && c.InStock == inStock {
			for _, subC := range c.Subcategories {
				if subC.Title == sTitle {
					return subC, true
				}
			}
		}
	}
	return domain.Subcategory{}, false
}

func (p *Provider) GetSubcategoryTitles(cTitle string, active bool, inStock bool) []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	var titles []string
	for _, c := range p.items {
		if c.Title == cTitle && c.InStock == inStock {
			for _, subC := range c.Subcategories {
				if subC.Active == active {
					titles = append(titles, subC.Title)
				}
			}
		}
	}
	return titles
}

func (p *Provider) GetProducts(cTitle, sTitle string, inStock bool) []domain.HouseholdProduct {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, c := range p.items {
		if c.InStock == inStock && c.Title == cTitle {
			for _, s := range c.Subcategories {
				if s.Title == sTitle {
					return s.Products
				}
			}
		}
	}
	return nil
}

func (p *Provider) MakeOnChangeHook() repositories.HouseholdOnChangeFunc {
	return func(items []domain.HouseholdCategory) {
		p.Load(items)
	}
}
