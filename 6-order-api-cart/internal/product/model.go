package product

import (
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string          `json:"name" gorm:"size:255;not null"`
	Description string          `json:"description" gorm:"type:text"`
	Images      pq.StringArray  `json:"images" gorm:"type:text[]"`
	IsActive    bool            `json:"is_active" gorm:"default:true"`
	Price       decimal.Decimal `json:"price" gorm:"type:decimal(17,2);default:0"`
}

func (Product) TableName() string {
	return "product"
}

func NewProduct(data *ProductCreateRequest) *Product {
	return &Product{
		Name:        data.Name,
		Description: data.Description,
		Images:      data.Images,
		IsActive:    data.IsActive,
		Price:       data.Price,
	}
}
