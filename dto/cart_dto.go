package dto

type AddCartItemRequest struct {
    ProductID uint `json:"product_id" binding:"required,gt=0"`
    Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

type UpdateCartItemRequest struct {
    Quantity int `json:"quantity" binding:"required,gt=0"`
}

type CartItemResponse struct {
    ID          uint   `json:"id"`
    ProductID   uint   `json:"product_id"`
    ProductName string `json:"product_name"`
    Image       string `json:"image"`
    Price       int64  `json:"price"`
    Quantity    int    `json:"quantity"`
    Subtotal    int64  `json:"subtotal"`
}

type CartResponse struct {
    ID          uint               `json:"id"`
    UserID      uint               `json:"user_id"`
    Items       []CartItemResponse `json:"items"`
    TotalItems  int                `json:"total_items"`
    TotalAmount int64              `json:"total_amount"`
}