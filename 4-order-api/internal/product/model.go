package product

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string         `json:"name" gorm:"size:255;not null"`
	Description string         `json:"description" gorm:"type:text"`
	Images      pq.StringArray `json:"images" gorm:"type:text[]"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
}

func (Product) TableName() string {
	return "product"
}
