// Copyright (c) 2022 Orlov Boris onlycodergod@gmail.com.
// This file service.go is subject to the terms and
// conditions defined in file 'LICENSE', which is part of this project source code.
package service

import "go_trainee/models"

type Orders interface {
	GetOrderByUID(string) *models.Order
	SaveOrderCache(*models.Order) error
	SaveOrderData(*models.Order) error
	LoadAllOrders() error
	UpdateCache([]byte)
	GetSaveOrderByUID(string) *models.OrderSave
	GetAllUID() []string
}

type Service interface {
	Orders
}
