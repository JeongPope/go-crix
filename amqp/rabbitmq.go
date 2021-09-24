package crixmq

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/streadway/amqp"

	"github.com/jeongpope/go-crix/logger"
	"github.com/jeongpope/go-crix/model"
)

var (
	EXCHANGE_KEY = []string{"UPBIT"}

	c             *amqp.Connection
	ch            *amqp.Channel
	reconnectLock *sync.Mutex

	chanReceive chan model.Ticker
)

func Initialize() (err error) {
	logger.Log.Println("[rabbitmq.go] Initialize")

	err = dial()
	if err != nil {
		logger.Log.Error("Check error string")

		return err
	}

	err = declare()
	if err != nil {
		Release()
		logger.Log.Error("Check error string")

		return err
	}

	reconnectLock = &sync.Mutex{}
	chanReceive = make(chan model.Ticker, 512)

	logger.Log.Println("[rabbitmq.go] Initialize Success")

	return err
}

func dial() (err error) {
	logger.Log.Println("[rabbitmq.go] dial")

	user := os.Getenv("RABBITMQ_USER_NAME")
	pwd := os.Getenv("RABBITMQ_USER_PASSWORD")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pwd, host, port)
	c, err = amqp.Dial(url)
	if err != nil {
		logger.Log.Error("Failed to connect to RabbitMQ")
		return err
	}
	logger.Log.Println("[rabbitmq.go] dial success")

	ch, err = c.Channel()
	if err != nil {
		Release()
		logger.Log.Error("Failed to open channel")
		logger.Log.Error(err.Error())

		return err
	}

	return err
}

func declare() (err error) {
	logger.Log.Println("[rabbitmq.go] declare")

	for _, v := range EXCHANGE_KEY {
		_, err := ch.QueueDeclare(
			v,     // Name
			false, // Durable
			false, // Delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)

		if err != nil {
			logger.Log.Error("Failed to declare a queue")
			return err
		}
	}

	logger.Log.Println("[rabbitmq.go] declare success")

	return err
}

func Reconnect() (err error) {
	logger.Log.Println("[rabbitmq.go] Reconnect")
	reconnectLock.Lock()
	ch.Close()
	c.Close()

	ticker := time.NewTicker(time.Second * 5)
	for {
		err = dial()
		if err != nil {
			logger.Log.Error("Failed to reconnect to RabbitMQ, retry after 10 second")
			logger.Log.Error(err.Error())
		}

		err = declare()
		if err != nil {
			logger.Log.Error("Failed to declare a queue, retry after 10 second")
			logger.Log.Error(err.Error())
		} else {
			logger.Log.Info("Successful reconnect")
			break
		}

		<-ticker.C
	}

	reconnectLock.Unlock()

	logger.Log.Println("[rabbitmq.go] Reconnect success")

	return err
}

func Release() {
	logger.Log.Println("[rabbitmq.go] Release")

	if ch != nil {
		ch.Close()
	}

	if c != nil {
		c.Close()
	}

	logger.Log.Println("[rabbitmq.go] Release success")
}

func Publish() (err error) {
	logger.Log.Println("[rabbitmq.go] Publish")
	for {
		msg := <-chanReceive
		jsonBytes, _ := json.Marshal(msg)

		err = ch.Publish(
			"",
			msg.Exchange,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        jsonBytes,
			})

		if err != nil {
			logger.Log.Error("Check error string")

			err := Reconnect()
			if err != nil {
				logger.Log.Error("Check error string")
			}
		}
	}
}

func GetChannel() chan model.Ticker {
	return chanReceive
}
