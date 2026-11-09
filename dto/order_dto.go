package dto

type OrderItemResponse struct {
    ID          uint   `json:"id"`
    ProductID   uint   `json:"product_id"`
    ProductName string `json:"product_name"`
    Image       string `json:"image"`
    Price       int64  `json:"price"`
    Quantity    int    `json:"quantity"`
    Subtotal    int64  `json:"subtotal"`
}

type OrderResponse struct {
    ID          uint                `json:"id"`
    UserID      uint                `json:"user_id"`
    TotalAmount int64               `json:"total_amount"`
    Status      string              `json:"status"`
    Items       []OrderItemResponse `json:"items"`
}

type UpdateOrderStatusRequest struct {
    Status string `json:"status" binding:"required,oneof=PROCESSING COMPLETED"`
}