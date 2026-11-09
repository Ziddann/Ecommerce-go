package handlers

import (
    "errors"
    "net/http"
    "strconv"

    "ecommerce/dto"
    "ecommerce/models"
    "ecommerce/services"

    "github.com/gin-gonic/gin"
)

type CategoryHandler struct {
    Service *services.CategoryService
}

func NewCategoryHandler(
    service *services.CategoryService,
) *CategoryHandler {
    return &CategoryHandler{
        Service: service,
    }
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
    categories, err := h.Service.GetAll()

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Failed to retrieve categories",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Categories retrieved successfully",
        "data":    categories,
    })
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)

    if err != nil || id == 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid category ID",
        })
        return
    }

    category, err := h.Service.GetByID(uint(id))

    if errors.Is(err, services.ErrCategoryNotFound) {
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "message": "Category not found",
        })
        return
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Failed to retrieve category",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Category retrieved successfully",
        "data":    category,
    })
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
    var request dto.CreateCategoryRequest

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid request data",
            "error":   err.Error(),
        })
        return
    }

    category := models.Category{
        Name:        request.Name,
        Description: request.Description,
    }

    if err := h.Service.Create(&category); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Failed to create category",
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "message": "Category created successfully",
        "data":    category,
    })
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)

    if err != nil || id == 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid category ID",
        })
        return
    }

    var request dto.UpdateCategoryRequest

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid request data",
            "error":   err.Error(),
        })
        return
    }

    updatedCategory := models.Category{
        Name:        request.Name,
        Description: request.Description,
    }

    err = h.Service.Update(uint(id), &updatedCategory)

    if errors.Is(err, services.ErrCategoryNotFound) {
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "message": "Category not found",
        })
        return
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Failed to update category",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Category updated successfully",
    })
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)

    if err != nil || id == 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid category ID",
        })
        return
    }

    err = h.Service.Delete(uint(id))

    if errors.Is(err, services.ErrCategoryNotFound) {
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "message": "Category not found",
        })
        return
    }

    if errors.Is(err, services.ErrCategoryHasProducts) {
        c.JSON(http.StatusConflict, gin.H{
            "success": false,
            "message": "Cannot delete category because it still has products",
        })
        return
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Failed to delete category",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Category deleted successfully",
    })
}