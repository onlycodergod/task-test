// Copyright (c) 2022 Orlov Boris onlycodergod@gmail.com.
// This file database.go is subject to the terms and
// conditions defined in file 'LICENSE', which is part of this project source code.
package cache

import (
	"fmt"

	"go_trainee/models"
)

func NewDatabaseCache() *DatabaseCache {
	return &DatabaseCache{order: make(map[string]*models.Order)}
}

func (cache *DatabaseCache) Change(order map[string]*models.Order) {
	cache.order = order
}

func (cache *DatabaseCache) GetOrderByUID(uid string) *models.Order {
	v, ok := cache.order[uid]

	if !ok {
		return nil
	}

	return v
}

func (cache *DatabaseCache) SaveOrder(order *models.Order) error {
	uid := order.OrderUID

	if _, ok := cache.order[uid]; ok {
		return fmt.Errorf("uid %s already in cache", uid)
	}

	cache.order[uid] = order

	return nil
}

func (cache *DatabaseCache) GetAllOrders() []*models.Order {
	orderSlice := make([]*models.Order, 0)

	for _, v := range cache.order {
		orderSlice = append(orderSlice, v)
	}

	return orderSlice
}
