package cache

import (
	"sync"

	"github.com/ds124wfegd/tech_wildberries_Go/internal/entity"
)

type Cache struct {
	orders map[string]*entity.Order
	mutex  sync.RWMutex
}

// create cache
func NewCache() *Cache {
	return &Cache{
		orders: make(map[string]*entity.Order),
	}
}

// save oder in cach
func (c *Cache) SetOrder(order *entity.Order) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.orders[order.OrderUID] = order
}

// retrieves an order from the cache
func (c *Cache) GetOrder(orderUID string) (*entity.Order, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	order, exists := c.orders[orderUID]
	return order, exists
}

// load orders in cache
func (c *Cache) LoadOrders(orders []*entity.Order) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, order := range orders {
		c.orders[order.OrderUID] = order
	}
}

// retrieves all order from the cache
func (c *Cache) GetAllOrders() []*entity.Order {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	orders := make([]*entity.Order, 0, len(c.orders))
	for _, order := range c.orders {
		orders = append(orders, order)
	}
	return orders
}

// Clear chache
func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.orders = make(map[string]*entity.Order)
}

// compatibility wit interface OrderCache type
func (c *Cache) Set(order *entity.Order)                   { c.SetOrder(order) }
func (c *Cache) Get(orderUID string) (*entity.Order, bool) { return c.GetOrder(orderUID) }
func (c *Cache) Load(orders []*entity.Order)               { c.LoadOrders(orders) }
