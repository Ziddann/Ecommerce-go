package services

import (
    "errors"

    "ecommerce/models"
    "ecommerce/repositories"

    "gorm.io/gorm"
)

var (
    ErrProductNotFound = errors.New("product not found")
    ErrInvalidCategory = errors.New("category not found")
)

type ProductService struct {
    ProductRepository  *repositories.ProductRepository
    CategoryRepository *repositories.CategoryRepository
}

func NewProductService(
    productRepository *repositories.ProductRepository,
    categoryRepository *repositories.CategoryRepository,
) *ProductService {
    return &ProductService{
        ProductRepository:  productRepository,
        CategoryRepository: categoryRepository,
    }
}

func (s *ProductService) GetAll() ([]models.Product, error) {
    return s.ProductRepository.GetAll()
}

func (s *ProductService) GetByID(id uint) (*models.Product, error) {
    product, err := s.ProductRepository.GetByID(id)

    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, ErrProductNotFound
    }

    if err != nil {
        return nil, err
    }

    return product, nil
}

func (s *ProductService) validateCategory(categoryID uint) error {
    _, err := s.CategoryRepository.GetByID(categoryID)

    if errors.Is(err, gorm.ErrRecordNotFound) {
        return ErrInvalidCategory
    }

    return err
}

func (s *ProductService) Create(product *models.Product) error {
    if err := s.validateCategory(product.CategoryID); err != nil {
        return err
    }

    return s.ProductRepository.Create(product)
}

func (s *ProductService) Update(
    id uint,
    updatedData *models.Product,
) error {
    product, err := s.GetByID(id)

    if err != nil {
        return err
    }

    if err := s.validateCategory(updatedData.CategoryID); err != nil {
        return err
    }

    product.Name = updatedData.Name
    product.Description = updatedData.Description
    product.Price = updatedData.Price
    product.Stock = updatedData.Stock
    product.Image = updatedData.Image
    product.CategoryID = updatedData.CategoryID

    return s.ProductRepository.Update(product)
}

func (s *ProductService) Delete(id uint) error {
    _, err := s.GetByID(id)

    if err != nil {
        return err
    }

    return s.ProductRepository.Delete(id)
}