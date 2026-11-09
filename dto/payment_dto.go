package dto

type CreatePaymentRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required"`
}

type SimulatePaymentRequest struct {
	Status string `json:"status" binding:"required,oneof=SUCCESS FAILED"`
}

type PaymentResponse struct {
	ID            uint   `json:"id"`
	OrderID       uint   `json:"order_id"`
	ReferenceCode string `json:"reference_code"`
	PaymentMethod string `json:"payment_method"`
	Amount        int64  `json:"amount"`
	Status        string `json:"status"`
}
