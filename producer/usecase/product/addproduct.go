package product

import (
	"encoding/json"
	"errors"
	"net/http"
	"producer/dto"
	"producer/repository"
	"producer/utils"
	"producer/validation"

	"go.uber.org/zap"
)

func AddProductHandler(w http.ResponseWriter, r *http.Request) {
	zap.L().Debug("AddProduct")

	// Response parameter
	res := &dto.Response{
		Code:    http.StatusBadRequest,
		Message: "Unable to process request",
	}

	defer func() {
		if re := recover(); re != nil {
			err := re.(error)
			zap.L().Error("Something went wrong.",
				zap.Error(err))
			res.Error = err.Error()
		}
		utils.Response(w, res)
	}()

	// Takes the value from the request
	req := &dto.RequestData{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		res.Error = err.Error()
		return
	}

	// Validation for empty fields
	err = validation.EmptyFieldValidation(req)
	if err != nil {
		zap.L().Error("Error in EmptyFieldValidation",
			zap.Error(err))
		res.Error = err.Error()
		res.Message = "Invalid inputs"
		return
	}

	// Inserting data into MySQL and send id to RMQ
	err = insertProduct(req)
	if err != nil {
		zap.L().Error("Unable to insert data",
			zap.Error(err))
		res.Error = err.Error()
		return
	}

	// returning success response
	res.Code = http.StatusOK
	res.Message = "Request processed successfully."
}

func insertProduct(req *dto.RequestData) error {
	// check if user_id exists in user table
	userIdExists, err := repository.GetUserByID(req.UserID)
	if err != nil {
		zap.L().Error("Unable to get data from MySQL",
			zap.Error(err))
		return err
	}

	if !userIdExists {
		return errors.New("user_id not found in DB")
	}

	// if userID exists insert the product data
	lastInsertedId, err := repository.InsertProduct(req)
	if err != nil {
		zap.L().Error("Unable to insert data into DB",
			zap.Error(err))
		return err
	}

	data := &dto.RMQProductData{
		Id: lastInsertedId,
	}

	// marshal data
	rmqBody, err := json.Marshal(data)
	if err != nil {
		zap.L().Error("Couldn't marshal data to json", zap.Error(err))
		return err
	}

	zap.L().Info("", zap.Any("njnjn", data))
	zap.L().Info("", zap.Any("njnjn", rmqBody))

	// sending data to RMQ
	return utils.PublishToQueue("product_data_queue", rmqBody, "add")
}
