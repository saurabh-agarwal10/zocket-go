package main

import (
	"consumer/config"
	"consumer/usecase/rmqconsumer"
	"log"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		zap.L().Fatal("unable to load .env file ")
	}

	initLogger()

	err = config.InitRabbitMQ()
	if err != nil {
		zap.L().Fatal("unable to initialize rabbitmq", zap.Error(err))
	}

	defer config.RMQDisconnect()

	// channel := make(chan bool)
	zap.L().Debug("Waiting for messages")

	//Consume data from product_data_queue
	productsData, err := rmqconsumer.RMQConsume("product_data_queue")
	if err != nil {
		zap.L().Error("Error in RMQConsume", zap.Error(err))
		return
	}

	zap.L().Info("", zap.Any("njnjn", productsData))

	for productData := range productsData {
		switch productData.Type {
		case "add":
			zap.L().Info("", zap.Any("njnjn", productData))
			err := rmqconsumer.DownloadImage(&productData)
			if err != nil {
				zap.L().Error("Error in DownloadImage",
					zap.Error(err))
				err := productData.Nack(true, true)
				if err != nil {
					zap.L().Error("Error in Nack",
						zap.Error(err))
				}
			}
		}
		productData.Ack(true)
	}
}

func initLogger() {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development: true,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	log.Println("log level development")

	logger, err := cfg.Build()
	if err != nil {
		log.Fatal("Unable to initiate the logger", err)
	}

	zap.L().Info("Log Enabled.")

	zap.ReplaceGlobals(logger)
	defer func() {
		_ = logger.Sync() // flushes buffer, if any
	}()
}
