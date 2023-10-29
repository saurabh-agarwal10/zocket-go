package validation

import (
	"errors"
	"producer/dto"
)

func EmptyFieldValidation(req *dto.RequestData) error {
	if req.ProductName == "" {
		return errors.New("product_name field is required")
	}

	if req.ProductDescription == "" {
		return errors.New("product_description field is required")
	}

	if len(req.ProductImages) == 0 {
		return errors.New("product_images field is required")
	}

	if req.ProductPrice == 0 {
		return errors.New("product_price field is required")
	}

	return nil
}
