// Copyright (c) 2022 Orlov Boris onlycodergod@gmail.com.
// This file postgres.go is subject to the terms and
// conditions defined in file 'LICENSE', which is part of this project source code.
package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go_trainee/models"

	"github.com/lib/pq"
	"go.uber.org/zap"
)

type Config struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DBName       string
	SSLMode      string
	MaxReconn    int
	MinReconn    int
	PingNoEvents int
}

const ordersCustomTable = "orderssave"

type Postgres struct {
	DB              *sql.DB
	listner         *pq.Listener
	config          Config
	listnerShutdown chan int
}

func NewPostgres(cfg Config) (*Postgres, error) {
	urlPostgres := fmt.Sprintf("host = %v port = %v user = %v dbname = %v password = %v sslmode = %v",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName,
		cfg.Password, cfg.SSLMode)

	db, err := sql.Open("postgres", urlPostgres)
	if err != nil {
		return nil, fmt.Errorf("failed to connect postgres %s", err)
	}

	err = db.Ping()

	if err != nil {
		return nil, fmt.Errorf("failed to ping postgres %s", err)
	}

	listner := createListner(cfg)

	return &Postgres{DB: db, listner: listner, config: cfg, listnerShutdown: make(chan int)}, nil
}

func createListner(cfg Config) *pq.Listener {
	urlPostgres := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Username, cfg.Password, cfg.Host,
		cfg.Port, cfg.DBName)

	minReconnInterval := time.Duration(cfg.MinReconn) * time.Second
	maxReconnInterval := time.Duration(cfg.MaxReconn) * time.Second

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	eventCallback := func(event pq.ListenerEventType, err error) {
		if err != nil {
			logger.Info("failed to initialize database Listner/Notify", zap.Error(err))
		}
	}

	return pq.NewListener(urlPostgres, minReconnInterval, maxReconnInterval, eventCallback)
}

func (db *Postgres) Close() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	if db.listner != nil {
		db.listnerShutdown <- 1
		err := db.listner.Close()
		if err != nil {
			logger.Error("failed to close listner", zap.Error(err))
		}
	}
	if db.DB != nil {
		err := db.DB.Close()
		if err != nil {
			logger.Error("failed to close database connection", zap.Error(err))
		}
	}
}

func (db *Postgres) RunService(channel string, f func([]byte)) error {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	if err := db.listner.Listen(channel); err != nil {
		return fmt.Errorf("failed to start listen a channel %s", channel)
	}

	go func() {
	loop:
		for {
			select {
			case <-db.listnerShutdown:
				break loop
			case notify := <-db.listner.Notify:
				f([]byte(notify.Extra)) // Payload
			case <-time.After(time.Duration(db.config.PingNoEvents) * time.Second):
				go func() {
					err := db.listner.Ping() // Ping the remote server
					if err != nil {
						logger.Fatal("failed to connect the remote server", zap.Error(err))
					}
				}()
			}
		}
	}()

	return nil
}

func (db *Postgres) DatabaseSaveOrder(order *models.Order) error {
	json, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order to json")
	}

	uid := order.OrderUID

	ok, err := db.CheckOrder(order)
	if err != nil {
		return err
	}

	if ok {
		return fmt.Errorf("failed to save order, because order already in database")
	}

	query := fmt.Sprintf("INSERT INTO %s (order_uid, data) VALUES ($1, $2)", ordersCustomTable)

	_, err = db.DB.Exec(query, uid, json)

	if err != nil {
		return fmt.Errorf("failed to save order in database %s", err)
	}

	return nil
}

func (db *Postgres) CheckOrder(order *models.Order) (bool, error) {
	orderToSave := models.OrderSave{}

	query := fmt.Sprintf("SELECT * FROM %s WHERE order_uid = $1", ordersCustomTable)

	row := db.DB.QueryRow(query, orderToSave)

	if err := row.Scan(&order.OrderUID); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		if err == nil {
			return true, nil
		}

		return false, fmt.Errorf("failed to check order in database %s", err)
	}

	return false, nil
}

func (db *Postgres) GetAllOrders() ([]models.OrderSave, error) {
	orderToSaveSlice := make([]models.OrderSave, 0)

	query := fmt.Sprintf("SELECT * FROM %s", ordersCustomTable)
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all orders from database %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&orderToSaveSlice); err != nil {
			return nil, fmt.Errorf("failed to get all orders from database %s", err)
		}
	}

	return orderToSaveSlice, nil
}
