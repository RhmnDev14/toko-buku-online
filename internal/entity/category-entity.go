package entity

type Category struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"size:100;not null"`
}

func (Category) TableName() string {
	return "categories"
}
