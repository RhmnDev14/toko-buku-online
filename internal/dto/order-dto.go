package dto

type Order struct {
	OrderItems []OrderItem
}

type OrderItem struct {
	BookId   int
	Price    float64
	Quantity int
}
