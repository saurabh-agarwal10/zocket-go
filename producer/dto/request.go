package dto

type RequestData struct {
	UserID             int     `json:"user_id"`
	ProductName        string  `json:"product_name"`
	ProductDescription string  `json:"product_description"`
	ProductImages      string  `json:"product_images"`
	ProductPrice       float64 `json:"product_price"`
}
