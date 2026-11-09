package repositories

import (
    "ecommerce/models"

    "gorm.io/gorm"
)

type CategoryRepository struct {
    DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
    return &CategoryRepository{
        DB: db,
    }
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
    var categories []models.Category

    err := r.DB.Find(&categories).Error

    return categories, err
}

func (r *CategoryRepository) GetByID(id uint) (*models.Category, error) {
    var category models.Category

    err := r.DB.First(&category, id).Error

    if err != nil {
        return nil, err
    }

    return &category, nil
}

func (r *CategoryRepository) Create(category *models.Category) error {
    return r.DB.Create(category).Error
}

func (r *CategoryRepository) Update(category *models.Category) error {
    return r.DB.Model(category).Updates(map[string]interface{}{
        "name":        category.Name,
        "description": category.Description,
    }).Error
}

func (r *CategoryRepository) Delete(id uint) error {
    return r.DB.Delete(&models.Category{}, id).Error
}

func (r *CategoryRepository) CountProducts(id uint) (int64, error) {
    var count int64

    err := r.DB.Model(&models.Product{}).
        Where("category_id = ?", id).
        Count(&count).Error

    return count, err
}