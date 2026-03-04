package services

import (
	"errors"

	"github.com/veetmoradiya3628/go-shop/internal/dto"
	"github.com/veetmoradiya3628/go-shop/internal/models"
	"gorm.io/gorm"
)

type CartService struct {
	db *gorm.DB
}

func NewCartService(db *gorm.DB) *CartService {
	return &CartService{db: db}
}

func (s *CartService) GetCart(userID uint) (*dto.CartResponse, error) {
	var cart models.Cart
	err := s.db.Preload("CartItems.Product.Category").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return s.convertToCartResponse(&cart), nil
}

func (s *CartService) AddToCart(userID uint, req *dto.AddToCartRequest) (*dto.CartResponse, error) {
	// check if product exists
	var product models.Product
	err := s.db.First(&product, req.ProductID).Error
	if err != nil {
		return nil, errors.New("product not found")
	}

	// check if stock is sufficient
	if product.Stock < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	// check if cart exists, if not create one
	var cart models.Cart
	err = s.db.Where("user_id = ?", userID).First(&cart).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		cart = models.Cart{UserID: userID}
		if err := s.db.Create(&cart).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	var cartItem models.CartItem
	err = s.db.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&cartItem).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// if cart item does not exist, create new one
		cartItem = models.CartItem{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		if err := s.db.Create(&cartItem).Error; err != nil {
			return nil, err
		}
	} else if errors.Is(err, gorm.ErrDuplicatedKey) {
		// if cart item already exists, update quantity
		cartItem.Quantity += req.Quantity
		if cartItem.Quantity > product.Stock {
			return nil, errors.New("insufficient stock")
		}
		if err := s.db.Save(&cartItem).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		// if any other error occurs, return it
		return nil, err
	}

	// update stock
	product.Stock -= req.Quantity
	if err := s.db.Save(&product).Error; err != nil {
		return nil, err
	}
	return s.GetCart(userID)
}

func (s *CartService) UpdateCartItem(userID uint, itemID uint, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error) {
	var cartItem models.CartItem
	if err := s.db.Joins("JOIN carts on cart_items.cart_id = carts.id").
		Where("cart_items.id = ? AND carts.user_id = ?", itemID, userID).
		First(&cartItem).Error; err != nil {
		return nil, errors.New("cart item not found")
	}

	// check if stock is sufficient
	if cartItem.Product.Stock+cartItem.Quantity < req.Quantity {
		return nil, errors.New("insufficient stock")
	}
	// update stock
	cartItem.Product.Stock += cartItem.Quantity - req.Quantity
	if err := s.db.Save(&cartItem.Product).Error; err != nil {
		return nil, err
	}
	// update cart item quantity
	cartItem.Quantity = req.Quantity
	if err := s.db.Save(&cartItem).Error; err != nil {
		return nil, err
	}
	return s.GetCart(userID)
}

func (s *CartService) RemoveFromCart(userID, itemID uint) error {
	return s.db.Joins("JOIN carts on cart_items.cart_id = carts.id").
		Where("cart_items.id = ? AND carts.user_id = ?", itemID, userID).
		Delete(&models.CartItem{}).Error
}

func (s *CartService) convertToCartResponse(cart *models.Cart) *dto.CartResponse {
	cartItems := make([]dto.CartItemResponse, len(cart.CartItems)) // memory allocation for cart items
	var total float64
	for i := range cart.CartItems {
		subtotal := float64(cart.CartItems[i].Quantity) * cart.CartItems[i].Product.Price
		total += subtotal

		cartItems[i] = dto.CartItemResponse{
			ID: cart.CartItems[i].ID,
			Product: dto.ProductResponse{
				ID:          cart.CartItems[i].Product.ID,
				CategoryID:  cart.CartItems[i].Product.CategoryID,
				Name:        cart.CartItems[i].Product.Name,
				Description: cart.CartItems[i].Product.Description,
				Price:       cart.CartItems[i].Product.Price,
				Stock:       cart.CartItems[i].Product.Stock,
				SKU:         cart.CartItems[i].Product.SKU,
				IsActive:    cart.CartItems[i].Product.IsActive,
				Category: dto.CategoryResponse{
					ID:          cart.CartItems[i].Product.Category.ID,
					Name:        cart.CartItems[i].Product.Category.Name,
					Description: cart.CartItems[i].Product.Category.Description,
					IsActive:    cart.CartItems[i].Product.Category.IsActive,
				},
			},
			Quantity: cart.CartItems[i].Quantity,
			Subtotal: subtotal,
		}
	}
	return &dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CartItems: cartItems,
		Total:     total,
	}
}
