package services

import (
    "ecommerce/dto"
    "ecommerce/repositories"
)

type CartService struct {
    Repository *repositories.CartRepository
}

func NewCartService(
    repository *repositories.CartRepository,
) *CartService {
    return &CartService{
        Repository: repository,
    }
}

// Get cart and calculate totals
func (s *CartService) GetCart(userID uint) (*dto.CartResponse, error) {

    cart, err := s.Repository.GetCart(userID)

    if err != nil {
        return nil, err
    }

    response := &dto.CartResponse{
        ID:    cart.ID,
        UserID: cart.UserID,
        Items: []dto.CartItemResponse{},
    }

    for _, item := range cart.Items {

        subtotal := item.Product.Price * int64(item.Quantity)

        response.Items = append(
            response.Items,
            dto.CartItemResponse{
                ID:          item.ID,
                ProductID:   item.ProductID,
                ProductName: item.Product.Name,
                Image:       item.Product.Image,
                Price:       item.Product.Price,
                Quantity:    item.Quantity,
                Subtotal:    subtotal,
            },
        )

        response.TotalItems += item.Quantity
        response.TotalAmount += subtotal
    }

    return response, nil
}

// Add item
func (s *CartService) AddItem(
    userID uint,
    productID uint,
    quantity int,
) error {

    return s.Repository.AddItem(
        userID,
        productID,
        quantity,
    )
}

// Update item
func (s *CartService) UpdateItem(
    userID uint,
    itemID uint,
    quantity int,
) error {

    return s.Repository.UpdateItem(
        userID,
        itemID,
        quantity,
    )
}

// Remove item
func (s *CartService) RemoveItem(
    userID uint,
    itemID uint,
) error {

    return s.Repository.RemoveItem(
        userID,
        itemID,
    )
}