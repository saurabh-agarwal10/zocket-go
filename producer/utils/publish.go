package utils

import (
	"producer/config"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

func PublishToQueue(qName string, rmqBody []byte, qType string) error {
	zap.L().Info("Publishing to queue")

	// Publish to rmq
	err := config.RMQChan.Publish(
		"",    // exchange
		qName, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        rmqBody,
			Type:        qType,
		})
	if err != nil {
		zap.L().Fatal("Error while publishing", zap.Error(err))
		return err
	}

	return nil
}
