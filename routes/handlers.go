// Copyright (c) 2022 Orlov Boris onlycodergod@gmail.com.
// This file handlers.go is subject to the terms and
// conditions defined in file 'LICENSE', which is part of this project source code.
package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (handler *Handler) getOrders(c *gin.Context) {
	orders := handler.service.GetAllUID()

	if orders == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND"})
	}

	c.JSON(http.StatusOK, orders)
}

func (handler *Handler) getOrderByUID(c *gin.Context) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	uid := c.Query("order_uid")

	saveOrder := handler.service.GetSaveOrderByUID(uid)

	if saveOrder == nil {
		logger.Error("failed found order", zap.Int("StatusNotFound", http.StatusNotFound))
		c.AbortWithStatusJSON(http.StatusNotFound, "failed found order")
		return
	}

	c.JSON(http.StatusOK, saveOrder)
}
