package entity

import "time"

type Order struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	UserID     uint      `gorm:"index;not null"`
	TotalPrice float64   `gorm:"type:decimal(10,2);not null"`
	Status     string    `gorm:"type:order_status;not null;default:'PENDING'"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

type OrderItem struct {
	ID       uint    `gorm:"primaryKey;autoIncrement"`
	OrderID  uint    `gorm:"index;not null"`
	BookID   uint    `gorm:"index;not null"`
	Quantity int     `gorm:"not null"`
	Price    float64 `gorm:"type:decimal(10,2);not null"`
}
