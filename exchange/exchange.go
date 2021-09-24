package exchange

import (
	"errors"

	"github.com/jeongpope/go-crix/goredis"
	"github.com/jeongpope/go-crix/logger"
)

var instance *stCrix

var (
	ErrFailedGetRequest = errors.New("retry request after 10 minute")
)

type stCrix struct {
	upbit *Upbit

	supportAsset []string
}

func GetInstance() *stCrix {
	if instance != nil {
		return instance
	}

	instance = new(stCrix)
	err := instance.initialize()
	if err != nil {
		logger.Log.Error("Failed to redis instance intialize, check redis server status.")
		return nil
	}

	return instance
}

func (i *stCrix) initialize() error {
	logger.Log.Info("[exchange.go] Start initialize()")

	i.initExchange()

	logger.Log.Info("[exchange.go] End initialize()")
	return nil
}

func (i *stCrix) initExchange() {
	logger.Log.Info("[exchange.go] Start initExchange()")

	i.upbit = new(Upbit)
	i.supportAsset = i.upbit.Initialize(nil)
	i.upbit.AttatchChannel(goredis.GetInstance().GetTickerChannel())

	logger.Log.Info("[exchange.go] End initExchange()")
}

func (i *stCrix) Update() {
	stopC := make(chan struct{})

	logger.Log.Info("[exchange.go] Start Update()")

	go i.upbit.Execute()

	logger.Log.Info("[exchange.go] End Update()")

	<-stopC
}

func (i *stCrix) Release() {
	logger.Log.Info("[exchange.go] Start Release()")

	i.upbit.Release()

	logger.Log.Info("[exchange.go] End Release()")
}
