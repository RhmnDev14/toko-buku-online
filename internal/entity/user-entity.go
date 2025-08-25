package entity

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:100;uniqueIndex;not null"`
	Password  string    `gorm:"size:255;not null"`
	Role      string    `gorm:"type:user_role;not null;default:'user'"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (User) TableName() string {
	return "users"
}
