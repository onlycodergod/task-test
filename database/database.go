// Copyright (c) 2022 Orlov Boris onlycodergod@gmail.com.
// This file database.go is subject to the terms and
// conditions defined in file 'LICENSE', which is part of this project source code.
package database

import "go_trainee/models"

type Orders interface {
	DatabaseSaveOrder(*models.Order) error
	GetAllOrders() ([]models.OrderSave, error)
	CheckOrder(*models.Order) (bool, error)
	RunService(string, func([]byte)) error
	Close()
}

type Database struct {
	Orders
}
