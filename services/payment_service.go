package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"

	"ecommerce/dto"
	"ecommerce/models"
	"ecommerce/repositories"
)

var (
	ErrInvalidPaymentMethod    = errors.New("invalid payment method")
	ErrInvalidSimulationStatus = errors.New("invalid simulation status")
)

type PaymentService struct {
	Repository *repositories.PaymentRepository
}

func NewPaymentService(
	repository *repositories.PaymentRepository,
) *PaymentService {
	return &PaymentService{
		Repository: repository,
	}
}

func generateReferenceCode() (string, error) {
	randomBytes := make([]byte, 8)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	return "PAY-" + strings.ToUpper(
		hex.EncodeToString(randomBytes),
	), nil
}

func mapPaymentToResponse(
	payment *models.Payment,
) dto.PaymentResponse {
	return dto.PaymentResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		ReferenceCode: payment.ReferenceCode,
		PaymentMethod: payment.PaymentMethod,
		Amount:        payment.Amount,
		Status:        payment.Status,
	}
}

func (s *PaymentService) CreatePayment(
	userID uint,
	orderID uint,
	paymentMethod string,
) (dto.PaymentResponse, error) {

	paymentMethod = strings.ToUpper(
		strings.TrimSpace(paymentMethod),
	)

	switch paymentMethod {
	case "VIRTUAL_ACCOUNT", "E_WALLET", "QRIS":
		// Metode pembayaran valid.
	default:
		return dto.PaymentResponse{}, ErrInvalidPaymentMethod
	}

	referenceCode, err := generateReferenceCode()
	if err != nil {
		return dto.PaymentResponse{}, err
	}

	payment, err := s.Repository.CreatePayment(
		userID,
		orderID,
		referenceCode,
		paymentMethod,
	)

	if err != nil {
		return dto.PaymentResponse{}, err
	}

	return mapPaymentToResponse(payment), nil
}

func (s *PaymentService) GetPayment(
	userID uint,
	paymentID uint,
) (dto.PaymentResponse, error) {

	payment, err := s.Repository.GetPayment(userID, paymentID)
	if err != nil {
		return dto.PaymentResponse{}, err
	}

	return mapPaymentToResponse(payment), nil
}

func (s *PaymentService) SimulatePayment(
	userID uint,
	paymentID uint,
	status string,
) (dto.PaymentResponse, error) {

	status = strings.ToUpper(strings.TrimSpace(status))

	if status != "SUCCESS" && status != "FAILED" {
		return dto.PaymentResponse{}, ErrInvalidSimulationStatus
	}

	payment, err := s.Repository.SimulatePayment(
		userID,
		paymentID,
		status,
	)

	if err != nil {
		return dto.PaymentResponse{}, err
	}

	return mapPaymentToResponse(payment), nil
}
