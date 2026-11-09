package handlers

import (
    "errors"
    "net/http"
    "strconv"

    "ecommerce/repositories"
    "ecommerce/services"
	"ecommerce/dto"

    "github.com/gin-gonic/gin"
)

type OrderHandler struct {
    Service *services.OrderService
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
    return &OrderHandler{
        Service: service,
    }
}

func (h *OrderHandler) Checkout(c *gin.Context) {
    userID := c.GetUint("user_id")

    order, err := h.Service.Checkout(userID)

    if err != nil {
        switch {
        case errors.Is(err, repositories.ErrCartEmpty):
            c.JSON(http.StatusBadRequest, gin.H{
                "success": false,
                "message": "Your cart is empty",
            })

        case errors.Is(err, repositories.ErrOrderProductUnavailable):
            c.JSON(http.StatusConflict, gin.H{
                "success": false,
                "message": "One or more products are unavailable",
            })

        case errors.Is(err, repositories.ErrOrderInsufficientStock):
            c.JSON(http.StatusConflict, gin.H{
                "success": false,
                "message": "Insufficient product stock",
            })

        default:
            c.JSON(http.StatusInternalServerError, gin.H{
                "success": false,
                "message": "Failed to checkout",
            })
        }

        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "message": "Order created successfully",
        "data":    order,
    })
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
    userID := c.GetUint("user_id")

    orders, err := h.Service.GetOrders(userID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Failed to retrieve orders",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    orders,
    })
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
    userID := c.GetUint("user_id")

    id, err := strconv.ParseUint(c.Param("id"), 10, 32)

    if err != nil || id == 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid order ID",
        })
        return
    }

    order, err := h.Service.GetOrderByID(userID, uint(id))

    if err != nil {
        if errors.Is(err, repositories.ErrOrderNotFound) {
            c.JSON(http.StatusNotFound, gin.H{
                "success": false,
                "message": "Order not found",
            })
            return
        }

        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Failed to retrieve order",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    order,
    })
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
    userID := c.GetUint("user_id")

    id, err := strconv.ParseUint(c.Param("id"), 10, 32)

    if err != nil || id == 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid order ID",
        })
        return
    }

    err = h.Service.CancelOrder(userID, uint(id))

    if err != nil {
        switch {
        case errors.Is(err, repositories.ErrOrderNotFound):
            c.JSON(http.StatusNotFound, gin.H{
                "success": false,
                "message": "Order not found",
            })

        case errors.Is(err, repositories.ErrOrderCannotBeCancelled):
            c.JSON(http.StatusConflict, gin.H{
                "success": false,
                "message": "Order cannot be cancelled",
            })

        default:
            c.JSON(http.StatusInternalServerError, gin.H{
                "success": false,
                "message": "Failed to cancel order",
            })
        }

        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Order cancelled successfully",
    })
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
    orders, err := h.Service.GetAllOrders()

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Failed to retrieve orders",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    orders,
    })
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)

    if err != nil || id == 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid order ID",
        })
        return
    }

    var request dto.UpdateOrderStatusRequest

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid order status",
        })
        return
    }

    err = h.Service.UpdateOrderStatus(uint(id), request.Status)

    if err != nil {
        switch {
        case errors.Is(err, repositories.ErrOrderNotFound):
            c.JSON(http.StatusNotFound, gin.H{
                "success": false,
                "message": "Order not found",
            })

        case errors.Is(err, repositories.ErrInvalidOrderStatus):
            c.JSON(http.StatusConflict, gin.H{
                "success": false,
                "message": "Invalid order status transition",
            })

        default:
            c.JSON(http.StatusInternalServerError, gin.H{
                "success": false,
                "message": "Failed to update order status",
            })
        }

        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Order status updated successfully",
    })
}
