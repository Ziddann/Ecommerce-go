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

type ProductHandler struct {
    Service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
    return &ProductHandler{
        Service: service,
    }
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
    products, err := h.Service.GetAll()

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Failed to retrieve products",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Products retrieved successfully",
        "data":    products,
    })
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)

    if err != nil || id == 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid product ID",
        })
        return
    }

    product, err := h.Service.GetByID(uint(id))

    if errors.Is(err, services.ErrProductNotFound) {
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "message": "Product not found",
        })
        return
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Failed to retrieve product",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Product retrieved successfully",
        "data":    product,
    })
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
    var request dto.CreateProductRequest

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid request data",
            "error":   err.Error(),
        })
        return
    }

    product := models.Product{
        Name:        request.Name,
        Description: request.Description,
        Price:       request.Price,
        Stock:       request.Stock,
        Image:       request.Image,
        CategoryID:  request.CategoryID,
    }

    err := h.Service.Create(&product)

    if errors.Is(err, services.ErrInvalidCategory) {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Category does not exist",
        })
        return
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Failed to create product",
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "message": "Product created successfully",
        "data":    product,
    })
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)

    if err != nil || id == 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid product ID",
        })
        return
    }

    var request dto.UpdateProductRequest

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid request data",
            "error":   err.Error(),
        })
        return
    }

    updatedProduct := models.Product{
        Name:        request.Name,
        Description: request.Description,
        Price:       request.Price,
        Stock:       request.Stock,
        Image:       request.Image,
        CategoryID:  request.CategoryID,
    }

    err = h.Service.Update(uint(id), &updatedProduct)

    if errors.Is(err, services.ErrProductNotFound) {
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "message": "Product not found",
        })
        return
    }

    if errors.Is(err, services.ErrInvalidCategory) {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Category does not exist",
        })
        return
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Failed to update product",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Product updated successfully",
    })
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)

    if err != nil || id == 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid product ID",
        })
        return
    }

    err = h.Service.Delete(uint(id))

    if errors.Is(err, services.ErrProductNotFound) {
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "message": "Product not found",
        })
        return
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Failed to delete product",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Product deleted successfully",
    })
}