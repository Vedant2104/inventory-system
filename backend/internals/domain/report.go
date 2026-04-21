package domain

type ProductCountByCategory struct {
	CategoryId string `json:"category_id"`
	Category   string `json:"category"`
	Count      int    `json:"count"`
}

type LowStockProducts struct {
	Id       string `json:"id"`
	Product  string `json:"product"`
	Brand    string `json:"brand"`
	Quantity int    `json:"quantity"`
}

type PriceSegmentation struct {
	CategoryId   string `json:"category_id"`
	Category string `json:"category"`
	Budget       int    `json:"budget"`
	MidRange     int    `json:"midRange"`
	Premium      int    `json:"premium"`
}
