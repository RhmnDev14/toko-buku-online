package dto

type BookReq struct {
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Year        int     `json:"year"`
	CategoryId  int     `json:"category_id"`
	ImageBase64 string  `json:"image_base64"`
}
