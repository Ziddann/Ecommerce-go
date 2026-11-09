package repositories

import (
    "errors"

    "ecommerce/models"

    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

var (
    ErrCartProductNotFound = errors.New("product not found")
    ErrCartItemNotFound    = errors.New("cart item not found")
    ErrInsufficientStock   = errors.New("insufficient product stock")
)

type CartRepository struct {
    DB *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
    return &CartRepository{
        DB: db,
    }
}

// Get or create cart
func (r *CartRepository) getOrCreateCart(
    tx *gorm.DB,
    userID uint,
) (*models.Cart, error) {

    var cart models.Cart

    err := tx.Where("user_id = ?", userID).First(&cart).Error

    if err == nil {
        return &cart, nil
    }

    if !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, err
    }

    cart = models.Cart{
        UserID: userID,
    }

    // Avoid duplicate cart creation for the same user.
    err = tx.Clauses(clause.OnConflict{
        Columns: []clause.Column{
            {Name: "user_id"},
        },
        DoNothing: true,
    }).Create(&cart).Error

    if err != nil {
        return nil, err
    }

    err = tx.Where("user_id = ?", userID).First(&cart).Error

    if err != nil {
        return nil, err
    }

    return &cart, nil
}

// Load cart with product information
func (r *CartRepository) loadCart(
    db *gorm.DB,
    userID uint,
) (*models.Cart, error) {

    var cart models.Cart

    err := db.
        Preload("Items.Product").
        Where("user_id = ?", userID).
        First(&cart).Error

    if err != nil {
        return nil, err
    }

    return &cart, nil
}

// Get cart
func (r *CartRepository) GetCart(
    userID uint,
) (*models.Cart, error) {

    err := r.DB.Transaction(func(tx *gorm.DB) error {
        _, err := r.getOrCreateCart(tx, userID)
        return err
    })

    if err != nil {
        return nil, err
    }

    return r.loadCart(r.DB, userID)
}

// Add product to cart
func (r *CartRepository) AddItem(
    userID uint,
    productID uint,
    quantity int,
) error {

    err := r.DB.Transaction(func(tx *gorm.DB) error {

        // 1. Dapatkan atau buat cart terlebih dahulu.
        cart, err := r.getOrCreateCart(tx, userID)
        if err != nil {
            return err
        }

        // 2. Lock product.
        var product models.Product

        err = tx.Clauses(clause.Locking{
            Strength: "UPDATE",
        }).First(&product, productID).Error

        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrCartProductNotFound
        }

        if err != nil {
            return err
        }

        // 3. Cari existing cart item.
        var item models.CartItem

        err = tx.Clauses(clause.Locking{
            Strength: "UPDATE",
        }).Where(
            "cart_id = ? AND product_id = ?",
            cart.ID,
            productID,
        ).First(&item).Error

        if err == nil {
            newQuantity := item.Quantity + quantity

            if newQuantity > product.Stock {
                return ErrInsufficientStock
            }

            return tx.Model(&item).Update(
                "quantity",
                newQuantity,
            ).Error
        }

        if !errors.Is(err, gorm.ErrRecordNotFound) {
            return err
        }

        // 4. Validasi stock.
        if quantity > product.Stock {
            return ErrInsufficientStock
        }

        // 5. Buat cart item baru.
        item = models.CartItem{
            CartID:    cart.ID,
            ProductID: productID,
            Quantity:  quantity,
        }

        return tx.Create(&item).Error
    })

    return err
}

// Update cart item quantity
func (r *CartRepository) UpdateItem(
    userID uint,
    itemID uint,
    quantity int,
) error {

    return r.DB.Transaction(func(tx *gorm.DB) error {

        cart, err := r.getOrCreateCart(tx, userID)

        if err != nil {
            return err
        }

        var item models.CartItem

        err = tx.Clauses(clause.Locking{
            Strength: "UPDATE",
        }).Where(
            "id = ? AND cart_id = ?",
            itemID,
            cart.ID,
        ).First(&item).Error

        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrCartItemNotFound
        }

        if err != nil {
            return err
        }

        var product models.Product

        err = tx.Clauses(clause.Locking{
            Strength: "UPDATE",
        }).First(&product, item.ProductID).Error

        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrCartProductNotFound
        }

        if err != nil {
            return err
        }

        if quantity > product.Stock {
            return ErrInsufficientStock
        }

        return tx.Model(&item).Update(
            "quantity",
            quantity,
        ).Error
    })
}

// Remove item from cart
func (r *CartRepository) RemoveItem(
    userID uint,
    itemID uint,
) error {

    return r.DB.Transaction(func(tx *gorm.DB) error {

        cart, err := r.getOrCreateCart(tx, userID)

        if err != nil {
            return err
        }

        var item models.CartItem

        err = tx.Where(
            "id = ? AND cart_id = ?",
            itemID,
            cart.ID,
        ).First(&item).Error

        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrCartItemNotFound
        }

        if err != nil {
            return err
        }

        // Permanently delete cart item.
        // This allows the same product to be added again.
        return tx.Unscoped().Delete(&item).Error
    })
}