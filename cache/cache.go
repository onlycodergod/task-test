// Copyright (c) 2022 Orlov Boris onlycodergod@gmail.com.
// This file chache.go is subject to the terms and
// conditions defined in file 'LICENSE', which is part of this project source code.
package cache

import "go_trainee/models"

type Orders interface {
	GetOrderByUID(string) *models.Order
	SaveOrder(*models.Order) error
	Change(map[string]*models.Order)
	GetAllOrders() []*models.Order
}

type DatabaseCache struct {
	Orders
	order map[string]*models.Order
}
