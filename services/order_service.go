package services

import (
    "ecommerce/dto"
    "ecommerce/models"
    "ecommerce/repositories"
)

type OrderService struct {
    Repository *repositories.OrderRepository
}

func NewOrderService(
    repository *repositories.OrderRepository,
) *OrderService {
    return &OrderService{
        Repository: repository,
    }
}

func mapOrderToResponse(order *models.Order) dto.OrderResponse {
    response := dto.OrderResponse{
        ID:          order.ID,
        UserID:      order.UserID,
        TotalAmount: order.TotalAmount,
        Status:      order.Status,
        Items:       make([]dto.OrderItemResponse, 0),
    }

    for _, item := range order.Items {
        response.Items = append(response.Items, dto.OrderItemResponse{
            ID:          item.ID,
            ProductID:   item.ProductID,
            ProductName: item.Product.Name,
            Image:       item.Product.Image,
            Price:       item.Price,
            Quantity:    item.Quantity,
            Subtotal:    item.Subtotal,
        })
    }

    return response
}

func (s *OrderService) Checkout(
    userID uint,
) (dto.OrderResponse, error) {
    order, err := s.Repository.Checkout(userID)
    if err != nil {
        return dto.OrderResponse{}, err
    }

    return mapOrderToResponse(order), nil
}

func (s *OrderService) GetOrders(
    userID uint,
) ([]dto.OrderResponse, error) {
    orders, err := s.Repository.GetOrders(userID)
    if err != nil {
        return nil, err
    }

    responses := make([]dto.OrderResponse, 0, len(orders))

    for i := range orders {
        responses = append(
            responses,
            mapOrderToResponse(&orders[i]),
        )
    }

    return responses, nil
}

func (s *OrderService) GetOrderByID(
    userID uint,
    orderID uint,
) (dto.OrderResponse, error) {
    order, err := s.Repository.GetOrderByID(userID, orderID)
    if err != nil {
        return dto.OrderResponse{}, err
    }

    return mapOrderToResponse(order), nil
}

func (s *OrderService) CancelOrder(
    userID uint,
    orderID uint,
) error {
    return s.Repository.CancelOrder(userID, orderID)
}

func (s *OrderService) GetAllOrders() ([]dto.OrderResponse, error) {
    orders, err := s.Repository.GetAllOrders()
    if err != nil {
        return nil, err
    }

    responses := make([]dto.OrderResponse, 0, len(orders))

    for i := range orders {
        responses = append(
            responses,
            mapOrderToResponse(&orders[i]),
        )
    }

    return responses, nil
}

func (s *OrderService) UpdateOrderStatus(
    orderID uint,
    status string,
) error {
    return s.Repository.UpdateOrderStatus(orderID, status)
}