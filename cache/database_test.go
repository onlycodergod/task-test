// Copyright (c) 2022 Orlov Boris onlycodergod@gmail.com.
// This file database_test.go is subject to the terms and
// conditions defined in file 'LICENSE', which is part of this project source code.
package cache

import (
	"reflect"
	"testing"

	"go_trainee/models"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseCacheNewDatabaseCache(t *testing.T) {
	t.Run("NotNil", func(t *testing.T) {
		cache := NewDatabaseCache()
		assert.NotNil(t, cache.order)
	})
}

func TestDatabaseCacheChange(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		cache := NewDatabaseCache()
		order := map[string]*models.Order{
			"102": {OrderUID: "102"},
		}
		cache.Change(order)
		assert.Equal(t, order, cache.order)
	})
}

func TestDatabaseCacheGetOrderByUID(t *testing.T) {
	tests := []struct {
		name         string
		cache        map[string]*models.Order
		uid          string
		expectResult *models.Order
	}{
		{
			name: "OK",
			cache: map[string]*models.Order{
				"102": {OrderUID: "102"},
			},
			uid:          "102",
			expectResult: &models.Order{OrderUID: "102"},
		},
		{
			name: "Empty Fields",
			cache: map[string]*models.Order{
				"102": {OrderUID: "102"},
			},
			uid:          "100",
			expectResult: nil,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			cache := DatabaseCache{}
			cache.Change(v.cache)
			result := cache.GetOrderByUID(v.uid)

			assert.Equal(t, reflect.DeepEqual(result, v.expectResult), true)
		})
	}
}

func TestDatabaseCacheSaveOrder(t *testing.T) {
	tests := []struct {
		name      string
		cache     map[string]*models.Order
		order     *models.Order
		expectErr bool
	}{
		{
			name: "OK",
			cache: map[string]*models.Order{
				"102": {OrderUID: "102"},
			},
		},
		{
			name: "Already OK",
			cache: map[string]*models.Order{
				"102": {OrderUID: "102"},
			},
			order:     &models.Order{OrderUID: "102"},
			expectErr: true,
		},
	}

	for _, v := range tests {
		cache := DatabaseCache{}
		cache.Change(v.cache)
		err := cache.SaveOrder(v.order)

		if v.expectErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestDatabaseCacheGetAllOrders(t *testing.T) {
	tests := []struct {
		name         string
		cache        map[string]*models.Order
		expectResult []*models.Order
	}{
		{
			name: "OK",
			cache: map[string]*models.Order{
				"102": {OrderUID: "102"},
			},
			expectResult: []*models.Order{
				{OrderUID: "102"},
			},
		},
		{
			name:         "Empty Fields",
			cache:        make(map[string]*models.Order, 0),
			expectResult: make([]*models.Order, 0),
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			cache := DatabaseCache{}
			cache.Change(v.cache)
			res := cache.GetAllOrders()

			assert.Equal(t, res, v.expectResult)
		})
	}
}
