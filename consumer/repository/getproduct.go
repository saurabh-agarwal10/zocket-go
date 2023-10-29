package repository

import (
	"consumer/config"
	"database/sql"
	"errors"

	"go.uber.org/zap"
)

const (
	getProductFromID = "SELECT product_images FROM products WHERE product_id=?"
	updateProduct    = "UPDATE products SET compressed_product_images=? WHERE product_id=?"
)

func GetProductFromID(productID int) (string, error) {
	var images string
	err := config.MysqlConnection.QueryRow(getProductFromID, productID).Scan(&images)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("ID doesn't exist")
			return "", errors.New("no matching id found")
		}
		zap.L().Error("Unable to get data from MySQL",
			zap.Error(err))
		return "", err
	}

	return images, nil
}

func UpdateProductImage(productID int, path string) error {
	_, err := config.MysqlConnection.Exec(updateProduct, path, productID)
	if err != nil {
		zap.L().Error("Unable to update data in MySQL",
			zap.Error(err))
		return err
	}

	return nil
}
