package repositories

import (
    "ecommerce/models"

    "gorm.io/gorm"
)

type ProductRepository struct {
    DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
    return &ProductRepository{
        DB: db,
    }
}

func (r *ProductRepository) GetAll() ([]models.Product, error) {
    var products []models.Product

    err := r.DB.Preload("Category").Find(&products).Error

    return products, err
}

func (r *ProductRepository) GetByID(id uint) (*models.Product, error) {
    var product models.Product

    err := r.DB.Preload("Category").First(&product, id).Error

    if err != nil {
        return nil, err
    }

    return &product, nil
}

func (r *ProductRepository) Create(product *models.Product) error {
    return r.DB.Create(product).Error
}

func (r *ProductRepository) Update(product *models.Product) error {
    return r.DB.Model(product).Updates(map[string]interface{}{
        "name":        product.Name,
        "description": product.Description,
        "price":       product.Price,
        "stock":       product.Stock,
        "image":       product.Image,
        "category_id": product.CategoryID,
    }).Error
}

func (r *ProductRepository) Delete(id uint) error {
    return r.DB.Delete(&models.Product{}, id).Error
}