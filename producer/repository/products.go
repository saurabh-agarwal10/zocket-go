package repository

import (
	"producer/config"
	"producer/dto"
	"time"

	"go.uber.org/zap"
)

const (
	insertProduct = "INSERT INTO products(product_name, product_description, product_images, product_price, created_at) VALUES(?,?,?,?,?)"
)

func InsertProduct(product *dto.RequestData) (int, error) {
	result, err := config.MysqlConnection.Exec(insertProduct, product.ProductName, product.ProductDescription, product.ProductImages, product.ProductPrice, time.Now().Local())
	if err != nil {
		zap.L().Error("Unable to insert data into MySQL",
			zap.Error(err))
		return 0, err
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		zap.L().Error("Unable to retrieve last Id from MySQL",
			zap.Error(err))
		return 0, err
	}

	return int(lastInsertedId), nil
}
