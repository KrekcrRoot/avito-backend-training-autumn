package dto

type Product struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	IsAvailable bool    `json:"is_available"`
}

type ProductsResponse struct {
	Products []Product `json:"products"`
}
