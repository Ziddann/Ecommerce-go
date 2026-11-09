package services

import (
    "errors"

    "ecommerce/models"
    "ecommerce/repositories"

    "gorm.io/gorm"
)

var (
    ErrCategoryNotFound    = errors.New("category not found")
    ErrCategoryHasProducts = errors.New("category still has products")
)

type CategoryService struct {
    Repository *repositories.CategoryRepository
}

func NewCategoryService(
    repository *repositories.CategoryRepository,
) *CategoryService {
    return &CategoryService{
        Repository: repository,
    }
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
    return s.Repository.GetAll()
}

func (s *CategoryService) GetByID(id uint) (*models.Category, error) {
    category, err := s.Repository.GetByID(id)

    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, ErrCategoryNotFound
    }

    if err != nil {
        return nil, err
    }

    return category, nil
}

func (s *CategoryService) Create(category *models.Category) error {
    return s.Repository.Create(category)
}

func (s *CategoryService) Update(
    id uint,
    updatedData *models.Category,
) error {
    category, err := s.GetByID(id)

    if err != nil {
        return err
    }

    category.Name = updatedData.Name
    category.Description = updatedData.Description

    return s.Repository.Update(category)
}

func (s *CategoryService) Delete(id uint) error {
    _, err := s.GetByID(id)

    if err != nil {
        return err
    }

    count, err := s.Repository.CountProducts(id)

    if err != nil {
        return err
    }

    if count > 0 {
        return ErrCategoryHasProducts
    }

    return s.Repository.Delete(id)
}