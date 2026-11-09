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

type PaymentHandler struct {
	Service *services.PaymentService
}

func NewPaymentHandler(
	service *services.PaymentService,
) *PaymentHandler {
	return &PaymentHandler{
		Service: service,
	}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	userID := c.GetUint("user_id")

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil || orderID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid order ID",
		})
		return
	}

	var request dto.CreatePaymentRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid payment request",
		})
		return
	}

	payment, err := h.Service.CreatePayment(
		userID,
		uint(orderID),
		request.PaymentMethod,
	)

	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidPaymentMethod):
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid payment method",
			})

		case errors.Is(err, repositories.ErrPaymentOrderNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Order not found",
			})

		case errors.Is(err, repositories.ErrPaymentOrderNotPayable):
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Order is not payable",
			})

		case errors.Is(err, repositories.ErrPendingPaymentExists):
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "A pending payment already exists",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to create payment",
			})
		}

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Payment created successfully",
		"data":    payment,
	})
}

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	userID := c.GetUint("user_id")

	paymentID, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil || paymentID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid payment ID",
		})
		return
	}

	payment, err := h.Service.GetPayment(userID, uint(paymentID))

	if err != nil {
		if errors.Is(err, repositories.ErrPaymentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Payment not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve payment",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    payment,
	})
}

func (h *PaymentHandler) SimulatePayment(c *gin.Context) {
	userID := c.GetUint("user_id")

	paymentID, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil || paymentID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid payment ID",
		})
		return
	}

	var request dto.SimulatePaymentRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Status must be SUCCESS or FAILED",
		})
		return
	}

	payment, err := h.Service.SimulatePayment(
		userID,
		uint(paymentID),
		request.Status,
	)

	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidSimulationStatus):
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid simulation status",
			})

		case errors.Is(err, repositories.ErrPaymentNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Payment not found",
			})

		case errors.Is(err, repositories.ErrPaymentNotPending):
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Payment has already been processed",
			})

		case errors.Is(err, repositories.ErrPaymentOrderNotPayable):
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Order is not payable",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to simulate payment",
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment simulation completed",
		"data":    payment,
	})
}
