package rmqconsumer

import (
	"consumer/config"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

func RMQConsume(qName string) (<-chan amqp.Delivery, error) {

	productData, err := config.RMQChan.Consume(
		qName, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		zap.L().Error("Couldn't consume data from queue.",
			zap.Error(err))
		return nil, err
	}

	return productData, nil
}
