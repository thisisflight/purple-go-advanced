package product

import (
	"purple/links/pkg/db"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *db.DB
}

func NewProductRepository(db *db.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (p *ProductRepository) Create(product *Product) (*Product, error) {
	result := p.db.Create(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil

}

func (p *ProductRepository) GetByID(id uint64) (*Product, error) {
	var product Product
	result := p.db.First(&product, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (p *ProductRepository) Update(product *Product) (*Product, error) {
	result := p.db.Updates(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (p *ProductRepository) Delete(id uint64) error {
	result := p.db.Delete(&Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
