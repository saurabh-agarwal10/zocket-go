package config

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var RMQConn *amqp.Connection
var RMQChan *amqp.Channel

func rabbitMQConnection() (*amqp.Connection, error) {
	Host := os.Getenv("RABBITMQ_HOST")
	Port := os.Getenv("RABBITMQ_PORT")
	Username := os.Getenv("RABBITMQ_USER")
	Password := os.Getenv("RABBITMQ_PASS")

	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%s", Username, Password, Host, Port)

	conn, err := amqp.Dial(connectionString)

	if err != nil {
		zap.L().Fatal("RabbitMQ Connection failed")
		return nil, err
	}

	zap.L().Info("RabbitMQ Successfully Connected.")

	return conn, nil
}

func RMQDisconnect() {
	_ = RMQConn.Close()
	_ = RMQChan.Close()

	zap.L().Info("RMQ connections and Channels closed.")
}

func InitRabbitMQ() error {

	if RMQConn == nil || RMQConn.IsClosed() {
		RMQConn, err = rabbitMQConnection()
		if err != nil {
			return err
		}
	}

	// rabbitmq channel initialization
	RMQChan, err = RMQConn.Channel()
	if err != nil {
		return err
	}

	zap.L().Info("new RMQ connection and channel successfully created.")

	// now declaring all the required queues and exchanges
	return declareQueues()
}

// function to declare queues
func declareQueues() error {
	//if queue declare properties are same as used in this function.
	queueList := []string{
		"product_data_queue",
	}

	for _, queue := range queueList {
		_, err := RMQChan.QueueDeclare(
			queue, //name
			true,  //durable
			false, //auto-delete
			false, //exclusive
			false, //noWait
			nil,   //args
		)

		if err != nil {
			zap.L().Fatal("Error declaring queue", zap.Error(err))
			return err
		}
	}
	return nil
}

func IsRMQConnClosed() bool {
	return (RMQConn != nil && RMQConn.IsClosed())
}
