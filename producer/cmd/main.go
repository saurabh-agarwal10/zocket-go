package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"producer/config"
	"producer/usecase/product"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var err error

func main() {
	err = godotenv.Load()
	if err != nil {
		zap.L().Fatal("unable to load .env file ")
	}

	initLogger()

	zap.L().Info("ENV Load Complete.")

	err = config.ConnectMysql()
	if err != nil {
		zap.L().Fatal("Unable to connect mysql ", zap.Error(err))
	}

	defer func() {
		_ = config.MysqlConnection.Close()
	}()

	err = config.InitRabbitMQ()
	if err != nil {
		zap.L().Fatal("unable to initialize rabbitmq", zap.Error(err))
	}

	defer config.RMQDisconnect()

	router := mux.NewRouter()

	router.HandleFunc("/product/add", product.AddProductHandler).Methods(http.MethodPost)

	zap.L().Info(fmt.Sprintf("Listening & Serving on : %s", os.Getenv("APP_PORT")))

	err = http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("APP_PORT")), router)
	if err != nil {
		zap.L().Error("Listening & Serving, error :", zap.Error(err))
		return
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
		_ = logger.Sync()
	}()
}
