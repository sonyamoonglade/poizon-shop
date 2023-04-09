package catalog

import (
	"sync"

	"domain"
	"repositories"
)

type CatalogProvider struct {
	mu    *sync.RWMutex
	items []domain.ClothingProduct
}

func NewCatalogProvider() *CatalogProvider {
	return &CatalogProvider{
		mu:    new(sync.RWMutex),
		items: nil,
	}
}

func (c *CatalogProvider) Load(items []domain.ClothingProduct) {
	c.mu.Lock()
	c.items = make([]domain.ClothingProduct, len(items), len(items))
	copy(c.items, items)
	c.mu.Unlock()
}

func (c *CatalogProvider) HasNext(offset uint) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return offset >= 0 && offset < uint(len(c.items)-1)
}

func (c *CatalogProvider) HasPrev(offset uint) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return offset > 0 && offset <= uint(len(c.items))
}

func (c *CatalogProvider) LoadNext(offset uint) domain.ClothingProduct {
	if !c.HasNext(offset) {
		return domain.ClothingProduct{}
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.items[offset+1]
}

func (c *CatalogProvider) LoadPrev(offset uint) domain.ClothingProduct {
	if !c.HasPrev(offset) {
		return domain.ClothingProduct{}
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.items[offset-1]
}

func (c *CatalogProvider) LoadFirst() domain.ClothingProduct {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.items) > 0 {
		return c.items[0]
	}
	return domain.ClothingProduct{}
}

func (c *CatalogProvider) LoadAt(offset uint) domain.ClothingProduct {
	if len(c.items) == 0 {
		return domain.ClothingProduct{}
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	_ = c.items[offset]
	return c.items[offset]
}

func MakeUpdateOnChangeFunc(catalogProvider *CatalogProvider) repositories.ClothingOnChangeFunc {
	return func(items []domain.ClothingProduct) {
		catalogProvider.Load(items)
	}
}
