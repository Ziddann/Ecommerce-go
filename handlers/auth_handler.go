package handlers

import (
    "errors"
    "net/http"

    "ecommerce/dto"
    "ecommerce/services"

    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    Service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
    return &AuthHandler{
        Service: service,
    }
}

func (h *AuthHandler) Register(c *gin.Context) {
    var request dto.RegisterRequest

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid request data",
            "error":   err.Error(),
        })
        return
    }

    user, err := h.Service.Register(request)

    if errors.Is(err, services.ErrEmailAlreadyRegistered) {
        c.JSON(http.StatusConflict, gin.H{
            "success": false,
            "message": "Email is already registered",
        })
        return
    }

    if errors.Is(err, services.ErrPasswordTooLong) {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Password must not exceed 72 bytes",
        })
        return
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Failed to register user",
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "message": "Registration successful",
        "data": dto.UserResponse{
            ID:    user.ID,
            Name:  user.Name,
            Email: user.Email,
            Role:  user.Role,
        },
    })
}

func (h *AuthHandler) Login(c *gin.Context) {
    var request dto.LoginRequest

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid request data",
            "error":   err.Error(),
        })
        return
    }

    token, user, err := h.Service.Login(
        request.Email,
        request.Password,
    )

    if errors.Is(err, services.ErrInvalidCredentials) {
        c.JSON(http.StatusUnauthorized, gin.H{
            "success": false,
            "message": "Invalid email or password",
        })
        return
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Failed to login",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Login successful",
        "data": gin.H{
            "token": token,
            "user": dto.UserResponse{
                ID:    user.ID,
                Name:  user.Name,
                Email: user.Email,
                Role:  user.Role,
            },
        },
    })
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
    userID := c.GetUint("user_id")

    user, err := h.Service.GetProfile(userID)

    if errors.Is(err, services.ErrUserNotFound) {
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "message": "User not found",
        })
        return
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Failed to retrieve profile",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Profile retrieved successfully",
        "data": dto.UserResponse{
            ID:    user.ID,
            Name:  user.Name,
            Email: user.Email,
            Role:  user.Role,
        },
    })
}