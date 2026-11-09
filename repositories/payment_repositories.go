package repositories

import (
    "errors"

    "ecommerce/models"

    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

var (
    ErrPaymentOrderNotFound = errors.New("order not found")
    ErrPaymentOrderNotPayable = errors.New("order is not payable")
    ErrPaymentNotFound = errors.New("payment not found")
    ErrPaymentNotPending = errors.New("payment is not pending")
    ErrPendingPaymentExists = errors.New("pending payment already exists")
)

type PaymentRepository struct {
    DB *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
    return &PaymentRepository{DB: db}
}

// CreatePayment membuat transaksi pembayaran baru.
func (r *PaymentRepository) CreatePayment(
    userID uint,
    orderID uint,
    referenceCode string,
    paymentMethod string,
) (*models.Payment, error) {

    var payment models.Payment

    err := r.DB.Transaction(func(tx *gorm.DB) error {
        // 1. Lock order dan pastikan order dimiliki user.
        var order models.Order

        err := tx.Clauses(clause.Locking{
            Strength: "UPDATE",
        }).Where(
            "id = ? AND user_id = ?",
            orderID,
            userID,
        ).First(&order).Error

        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrPaymentOrderNotFound
        }

        if err != nil {
            return err
        }

        // 2. Pastikan order belum dibayar.
        if order.Status != "PENDING_PAYMENT" {
            return ErrPaymentOrderNotPayable
        }

        // 3. Pastikan belum ada payment yang masih pending.
        var existingPayment models.Payment

        err = tx.Where(
            "order_id = ? AND status = ?",
            order.ID,
            "PENDING",
        ).First(&existingPayment).Error

        if err == nil {
            return ErrPendingPaymentExists
        }

        if !errors.Is(err, gorm.ErrRecordNotFound) {
            return err
        }

        // 4. Buat payment baru.
        payment = models.Payment{
            OrderID:       order.ID,
            ReferenceCode: referenceCode,
            PaymentMethod: paymentMethod,
            Amount:        order.TotalAmount,
            Status:        "PENDING",
        }

        return tx.Create(&payment).Error
    })

    if err != nil {
        return nil, err
    }

    return &payment, nil
}

// GetPayment mengambil payment berdasarkan ID dan pemilik order.
func (r *PaymentRepository) GetPayment(
    userID uint,
    paymentID uint,
) (*models.Payment, error) {

    var payment models.Payment

    err := r.DB.
        Joins("JOIN orders ON orders.id = payments.order_id").
        Where(
            "payments.id = ? AND orders.user_id = ?",
            paymentID,
            userID,
        ).
        First(&payment).Error

    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, ErrPaymentNotFound
    }

    if err != nil {
        return nil, err
    }

    return &payment, nil
}

// SimulatePayment mengubah status pembayaran secara simulasi.
func (r *PaymentRepository) SimulatePayment(
    userID uint,
    paymentID uint,
    newStatus string,
) (*models.Payment, error) {

    var payment models.Payment

    err := r.DB.Transaction(func(tx *gorm.DB) error {
        // 1. Ambil informasi payment untuk mengetahui order ID.
        var initialPayment models.Payment

        err := tx.First(
            &initialPayment,
            paymentID,
        ).Error

        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrPaymentNotFound
        }

        if err != nil {
            return err
        }

        // 2. Lock order terlebih dahulu.
        var order models.Order

        err = tx.Clauses(clause.Locking{
            Strength: "UPDATE",
        }).Where(
            "id = ? AND user_id = ?",
            initialPayment.OrderID,
            userID,
        ).First(&order).Error

        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrPaymentNotFound
        }

        if err != nil {
            return err
        }

        // 3. Lock payment.
        err = tx.Clauses(clause.Locking{
            Strength: "UPDATE",
        }).First(&payment, paymentID).Error

        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrPaymentNotFound
        }

        if err != nil {
            return err
        }

        // 4. Pastikan payment masih pending.
        if payment.Status != "PENDING" {
            return ErrPaymentNotPending
        }

        // 5. Update status payment.
        payment.Status = newStatus

        if err := tx.Save(&payment).Error; err != nil {
            return err
        }

        // 6. Jika pembayaran berhasil, update order.
        if newStatus == "SUCCESS" {
            if order.Status != "PENDING_PAYMENT" {
                return ErrPaymentOrderNotPayable
            }

            order.Status = "PAID"

            if err := tx.Save(&order).Error; err != nil {
                return err
            }
        }

        return nil
    })

    if err != nil {
        return nil, err
    }

    return &payment, nil
}