package repositories

import (
    "errors"
    "sort"

    "ecommerce/models"

    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

var (
    ErrCartEmpty = errors.New("cart is empty")
    ErrOrderProductUnavailable = errors.New("product is unavailable")
    ErrOrderInsufficientStock = errors.New("insufficient product stock")
    ErrOrderNotFound = errors.New("order not found")
	ErrOrderCannotBeCancelled = errors.New("order cannot be cancelled")
	ErrInvalidOrderStatus = errors.New("invalid order status transition")
)

type OrderRepository struct {
    DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
    return &OrderRepository{DB: db}
}

// Checkout mengubah cart menjadi order.
func (r *OrderRepository) Checkout(userID uint) (*models.Order, error) {
    var order models.Order

    err := r.DB.Transaction(func(tx *gorm.DB) error {
        // 1. Lock cart milik user.
        var cart models.Cart

        err := tx.Clauses(clause.Locking{
            Strength: "UPDATE",
        }).Where("user_id = ?", userID).First(&cart).Error

        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrCartEmpty
        }

        if err != nil {
            return err
        }

        // 2. Ambil semua item cart.
        var cartItems []models.CartItem

        err = tx.Where("cart_id = ?", cart.ID).
            Order("id ASC").
            Find(&cartItems).Error

        if err != nil {
            return err
        }

        if len(cartItems) == 0 {
            return ErrCartEmpty
        }

        // 3. Urutkan product ID untuk menjaga urutan lock.
        productIDs := make([]uint, 0, len(cartItems))

        for _, item := range cartItems {
            productIDs = append(productIDs, item.ProductID)
        }

        sort.Slice(productIDs, func(i, j int) bool {
            return productIDs[i] < productIDs[j]
        })

        // 4. Lock seluruh produk dan validasi stok.
        products := make(map[uint]models.Product)

        for _, productID := range productIDs {
            var product models.Product

            err := tx.Clauses(clause.Locking{
                Strength: "UPDATE",
            }).First(&product, productID).Error

            if errors.Is(err, gorm.ErrRecordNotFound) {
                return ErrOrderProductUnavailable
            }

            if err != nil {
                return err
            }

            products[productID] = product
        }

        // 5. Hitung total dan buat snapshot order item.
        var totalAmount int64
        orderItems := make([]models.OrderItem, 0, len(cartItems))

        for _, cartItem := range cartItems {
            product, exists := products[cartItem.ProductID]

            if !exists {
                return ErrOrderProductUnavailable
            }

            if cartItem.Quantity <= 0 {
                return ErrOrderInsufficientStock
            }

            if product.Stock < cartItem.Quantity {
                return ErrOrderInsufficientStock
            }

            subtotal := product.Price * int64(cartItem.Quantity)
            totalAmount += subtotal

            orderItems = append(orderItems, models.OrderItem{
                ProductID: product.ID,
                Quantity:  cartItem.Quantity,
                Price:     product.Price,
                Subtotal:  subtotal,
            })
        }

        // 6. Buat order dengan status menunggu pembayaran.
        order = models.Order{
            UserID:      userID,
            TotalAmount: totalAmount,
            Status:      "PENDING_PAYMENT",
        }

        if err := tx.Create(&order).Error; err != nil {
            return err
        }

        // 7. Simpan detail order.
        for i := range orderItems {
            orderItems[i].OrderID = order.ID
        }

        if err := tx.Create(&orderItems).Error; err != nil {
            return err
        }

        // 8. Kurangi stok produk.
        for _, cartItem := range cartItems {
            product := products[cartItem.ProductID]

            result := tx.Model(&models.Product{}).
                Where("id = ? AND stock >= ?",
                    product.ID,
                    cartItem.Quantity,
                ).
                Update("stock", gorm.Expr(
                    "stock - ?",
                    cartItem.Quantity,
                ))

            if result.Error != nil {
                return result.Error
            }

            if result.RowsAffected == 0 {
                return ErrOrderInsufficientStock
            }
        }

        // 9. Kosongkan cart setelah semua proses berhasil.
        if err := tx.Unscoped().
            Where("cart_id = ?", cart.ID).
            Delete(&models.CartItem{}).Error; err != nil {
            return err
        }

        order.Items = orderItems

        return nil
    })

    if err != nil {
        return nil, err
    }

    return &order, nil
}

// Mendapatkan seluruh order milik user.
func (r *OrderRepository) GetOrders(
    userID uint,
) ([]models.Order, error) {
    var orders []models.Order

    err := r.DB.
        Preload("Items.Product").
        Where("user_id = ?", userID).
        Order("created_at DESC").
        Find(&orders).Error

    return orders, err
}

// Mendapatkan detail order berdasarkan ID dan pemiliknya.
func (r *OrderRepository) GetOrderByID(
    userID uint,
    orderID uint,
) (*models.Order, error) {
    var order models.Order

    err := r.DB.
        Preload("Items.Product").
        Where("id = ? AND user_id = ?", orderID, userID).
        First(&order).Error

    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, ErrOrderNotFound
    }

    if err != nil {
        return nil, err
    }

    return &order, nil
}

func (r *OrderRepository) CancelOrder(
    userID uint,
    orderID uint,
) error {

    return r.DB.Transaction(func(tx *gorm.DB) error {
        // 1. Lock order dan validasi kepemilikan.
        var order models.Order

        err := tx.Clauses(clause.Locking{
            Strength: "UPDATE",
        }).Where(
            "id = ? AND user_id = ?",
            orderID,
            userID,
        ).First(&order).Error

        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrOrderNotFound
        }

        if err != nil {
            return err
        }

        // 2. Order hanya bisa dibatalkan jika belum dibayar.
        if order.Status != "PENDING_PAYMENT" {
            return ErrOrderCannotBeCancelled
        }

        // 3. Pastikan tidak ada pembayaran yang berhasil.
        var successfulPayment models.Payment

        err = tx.Where(
            "order_id = ? AND status = ?",
            order.ID,
            "SUCCESS",
        ).First(&successfulPayment).Error

        if err == nil {
            return ErrOrderCannotBeCancelled
        }

        if !errors.Is(err, gorm.ErrRecordNotFound) {
            return err
        }

        // 4. Ambil seluruh order item.
        var items []models.OrderItem

        if err := tx.Where(
            "order_id = ?",
            order.ID,
        ).Order("product_id ASC").Find(&items).Error; err != nil {
            return err
        }

        // 5. Kembalikan stok produk.
        for _, item := range items {
            result := tx.Model(&models.Product{}).
                Where("id = ?", item.ProductID).
                Update(
                    "stock",
                    gorm.Expr("stock + ?", item.Quantity),
                )

            if result.Error != nil {
                return result.Error
            }

            if result.RowsAffected == 0 {
                return ErrOrderProductUnavailable
            }
        }

        // 6. Ubah payment yang masih pending menjadi expired.
        if err := tx.Model(&models.Payment{}).
            Where(
                "order_id = ? AND status = ?",
                order.ID,
                "PENDING",
            ).
            Update("status", "EXPIRED").Error; err != nil {
            return err
        }

        // 7. Ubah status order.
        order.Status = "CANCELLED"

        return tx.Save(&order).Error
    })
}

func (r *OrderRepository) GetAllOrders() ([]models.Order, error) {
    var orders []models.Order

    err := r.DB.
        Preload("Items.Product").
        Order("created_at DESC").
        Find(&orders).Error

    return orders, err
}

func (r *OrderRepository) UpdateOrderStatus(
    orderID uint,
    newStatus string,
) error {

    return r.DB.Transaction(func(tx *gorm.DB) error {
        var order models.Order

        err := tx.Clauses(clause.Locking{
            Strength: "UPDATE",
        }).First(&order, orderID).Error

        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrOrderNotFound
        }

        if err != nil {
            return err
        }

        // PAID -> PROCESSING
        // PROCESSING -> COMPLETED
        validTransition :=
            (order.Status == "PAID" && newStatus == "PROCESSING") ||
            (order.Status == "PROCESSING" && newStatus == "COMPLETED")

        if !validTransition {
            return ErrInvalidOrderStatus
        }

        order.Status = newStatus

        return tx.Save(&order).Error
    })
}