package goredis

import (
	"errors"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jeongpope/go-crix/logger"
	"github.com/jeongpope/go-crix/model"
)

var instance *stRedis

type stRedis struct {
	pool       *redis.Pool
	chanTicker chan model.Ticker

	// Environment
	host      string
	port      string
	dbNumber  string
	maxIdle   int
	maxActive int
}

func GetInstance() *stRedis {
	if instance != nil {
		return instance
	}

	err := initialize()
	if err != nil {
		logger.Log.Error("Failed to redis instance intialize, check redis server status.")
		return nil
	}

	return instance
}

func (i *stRedis) GetTickerChannel() chan model.Ticker {
	return i.chanTicker
}

func initialize() error {
	logger.Log.Info("[redis.go] Start initialze()")

	instance = new(stRedis)

	instance.host = os.Getenv("REDIS_HOST")
	instance.port = os.Getenv("REDIS_PORT")
	instance.dbNumber = os.Getenv("REDIS_DB_NUMBER")
	instance.maxIdle = 80
	instance.maxActive = 12000

	instance.pool = &redis.Pool{
		MaxIdle:   instance.maxIdle,
		MaxActive: instance.maxActive,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", instance.host+":"+instance.port)
			if err != nil {
				logger.Log.Errorf("Failed to dial redis, %s", err.Error())
			}

			return conn, err
		},
	}

	instance.chanTicker = make(chan model.Ticker)

	logger.Log.Info("[redis.go] End Initialze()")
	return nil
}

func ping(conn redis.Conn) error {
	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		return errors.New("Failed receive PING/PONG response, " + err.Error())
	}

	return nil
}

func (i *stRedis) Update() {
	logger.Log.Info("[redis.go] Start Update()")

	go func() {
		for {
			var conn redis.Conn

			for {
				conn = instance.pool.Get()
				err := ping(conn)
				if err != nil {
					logger.Log.Errorf(err.Error())
					time.Sleep(time.Second * 5)
				} else {
					break
				}
			}

			for {
				if msg, openChannel := <-instance.chanTicker; openChannel {
					_, err := redis.String(conn.Do("SELECT", instance.dbNumber))
					if err != nil {
						logger.Log.Errorf("Failed select db, %s", err.Error())
						continue
					}

					_, err = redis.Int64(conn.Do("RPUSH", "CRIX", msg))
					if err != nil {
						logger.Log.Errorf("Failed ticker push message, %s", err.Error())
						continue
					}
				} else {
					logger.Log.Info("Redis receive channel is closed.")
					break
				}
			}

			logger.Log.Info("[redis.go] End ticker Update()")
			break
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Second * 10)

		for {
			<-ticker.C

			var conn redis.Conn

			for {
				conn = instance.pool.Get()
				err := ping(conn)
				if err != nil {
					logger.Log.Errorf(err.Error())
					time.Sleep(time.Second * 5)
				} else {
					break
				}
			}

			//conn.Do("flushall")
		}
	}()
}

func (i *stRedis) Release() {
	instance.pool.Close()
	close(instance.chanTicker)
}
