package user

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	PhoneNumber string `json:"phone_number" gorm:"uniqueIndex;size:20;not null"`
	Name        string `json:"name" gorm:"type:varchar(100)"`
}

func (User) TableName() string {
	return "user"
}
