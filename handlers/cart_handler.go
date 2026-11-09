package handlers

import (
    "errors"
    "net/http"
    "strconv"

    "ecommerce/dto"
    "ecommerce/repositories"
    "ecommerce/services"

    "github.com/gin-gonic/gin"
)

type CartHandler struct {
    Service *services.CartService
}

func NewCartHandler(service *services.CartService) *CartHandler {
    return &CartHandler{
        Service: service,
    }
}

// GET /api/v1/cart
func (h *CartHandler) GetCart(c *gin.Context) {

    userID := c.GetUint("user_id")

    cart, err := h.Service.GetCart(userID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Failed to retrieve cart",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Cart retrieved successfully",
        "data":    cart,
    })
}

// POST /api/v1/cart/items
func (h *CartHandler) AddItem(c *gin.Context) {

    userID := c.GetUint("user_id")

    var request dto.AddCartItemRequest

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid request data",
            "error":   err.Error(),
        })
        return
    }

    err := h.Service.AddItem(
        userID,
        request.ProductID,
        request.Quantity,
    )

    if errors.Is(err, repositories.ErrCartProductNotFound) {
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "message": "Product not found",
        })
        return
    }

    if errors.Is(err, repositories.ErrInsufficientStock) {
        c.JSON(http.StatusConflict, gin.H{
            "success": false,
            "message": "Insufficient product stock",
        })
        return
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Failed to add product to cart",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Product added to cart successfully",
    })
}

// PUT /api/v1/cart/items/:id
func (h *CartHandler) UpdateItem(c *gin.Context) {

    userID := c.GetUint("user_id")

    itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)

    if err != nil || itemID == 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid cart item ID",
        })
        return
    }

    var request dto.UpdateCartItemRequest

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid request data",
            "error":   err.Error(),
        })
        return
    }

    err = h.Service.UpdateItem(
        userID,
        uint(itemID),
        request.Quantity,
    )

    if errors.Is(err, repositories.ErrCartItemNotFound) {
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "message": "Cart item not found",
        })
        return
    }

    if errors.Is(err, repositories.ErrCartProductNotFound) {
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "message": "Product not found",
        })
        return
    }

    if errors.Is(err, repositories.ErrInsufficientStock) {
        c.JSON(http.StatusConflict, gin.H{
            "success": false,
            "message": "Insufficient product stock",
        })
        return
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Failed to update cart item",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Cart item updated successfully",
    })
}

// DELETE /api/v1/cart/items/:id
func (h *CartHandler) RemoveItem(c *gin.Context) {

    userID := c.GetUint("user_id")

    itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)

    if err != nil || itemID == 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid cart item ID",
        })
        return
    }

    err = h.Service.RemoveItem(
        userID,
        uint(itemID),
    )

    if errors.Is(err, repositories.ErrCartItemNotFound) {
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "message": "Cart item not found",
        })
        return
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Failed to remove cart item",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Cart item removed successfully",
    })
}