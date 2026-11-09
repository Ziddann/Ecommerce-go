package dto

type CreateProductRequest struct {
    Name        string `json:"name" binding:"required,min=3,max=150"`
    Description string `json:"description"`
    Price       int64  `json:"price" binding:"required,gt=0"`
    Stock       int    `json:"stock" binding:"gte=0"`
    Image       string `json:"image"`
    CategoryID  uint   `json:"category_id" binding:"required,gt=0"`
}

type UpdateProductRequest struct {
    Name        string `json:"name" binding:"required,min=3,max=150"`
    Description string `json:"description"`
    Price       int64  `json:"price" binding:"required,gt=0"`
    Stock       int    `json:"stock" binding:"gte=0"`
    Image       string `json:"image"`
    CategoryID  uint   `json:"category_id" binding:"required,gt=0"`
}