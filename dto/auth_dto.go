package dto

type RegisterRequest struct {
    Name     string `json:"name" binding:"required,min=3,max=100"`
    Email    string `json:"email" binding:"required,email,max=150"`
    Password string `json:"password" binding:"required,min=8,max=72"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type UserResponse struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  string `json:"role"`
}