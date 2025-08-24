package entity

type Book struct {
	ID          uint    `gorm:"primaryKey;autoIncrement"`
	Title       string  `gorm:"size:255;not null"`
	Author      string  `gorm:"size:100;not null"`
	Price       float64 `gorm:"type:decimal(10,2);not null"`
	Stock       int     `gorm:"not null"`
	Year        int
	CategoryID  uint   `gorm:"index;not null"`
	ImageBase64 string `gorm:"type:text"`
}
