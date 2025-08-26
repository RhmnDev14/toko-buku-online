package dto

type Order struct {
	OrderItems []OrderItem `json:"orderItems"`
}

type OrderItem struct {
	BookId   int `json:"bookId"`
	Quantity int `json:"quantity"`
}
