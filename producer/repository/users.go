package repository

import (
	"producer/config"

	"go.uber.org/zap"
)

const (
	getUserByID = "SELECT EXISTS(SELECT * FROM users WHERE user_id=?)"
)

func GetUserByID(userID int) (bool, error) {
	var exists int
	err := config.MysqlConnection.QueryRow(getUserByID, userID).Scan(&exists)
	if err != nil {
		zap.L().Error("Unable to get data from MySQL",
			zap.Error(err))
		return false, err
	}

	return exists == 1, nil
}
