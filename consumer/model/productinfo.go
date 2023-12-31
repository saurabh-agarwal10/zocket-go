package model

import "time"

type ProductData struct {
	ProductID               int       `json:"product_id"`
	ProductName             string    `json:"product_name"`
	ProductDescription      string    `json:"product_description"`
	ProductImages           string    `json:"product_images"`
	ProductPrice            float64   `json:"product_price"`
	CompressedProductImages string    `json:"compressed_product_images"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}
