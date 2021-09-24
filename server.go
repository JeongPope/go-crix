package main

import (
	"errors"
	"os"
	"time"

	"github.com/jeongpope/go-crix/exchange"
	"github.com/jeongpope/go-crix/goredis"
	"github.com/jeongpope/go-crix/logger"
	"github.com/jeongpope/go-crix/routes"
)

var (
	ErrFailedInitRedis    = errors.New("failed to initialize redis")
	ErrFailedInitExchange = errors.New("failed to initialize exchange")
)

func main() {
	logger.Log.Info("[server.go] Start main()")

	err := initialize()
	if err != nil {
		logger.Log.Error(err.Error())

		return
	}

	update()
	release()

	logger.Log.Info("[server.go] End main()")
}

func initialize() error {
	logger.Log.Info("[server.go] Start initialize()")

	// Logger
	t := time.Now()
	fileName := "logger/logs/" + t.Format(time.RFC3339) + ".log"
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	logger.Log.SetOut(f)
	logger.Log.SetLevel(logger.ERROR)

	// Routes
	routes.Initialize()

	// Redis
	if goredis.GetInstance() == nil {
		return ErrFailedInitRedis
	}

	// Exchange
	if exchange.GetInstance() == nil {
		return ErrFailedInitExchange
	}

	logger.Log.Info("[server.go] End initialize()")
	return nil
}

func update() {
	logger.Log.Info("[server.go] Start update()")

	goredis.GetInstance().Update()
	exchange.GetInstance().Update()

	logger.Log.Info("[server.go] End update()")
}

func release() {
	logger.Log.Info("[server.go] Start release()")

	logger.Log.Info("[server.go] End release()")
}
