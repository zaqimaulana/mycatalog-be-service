package services

import (
	"errors"

	"github.com/zaqimaulana/mycatalog-be/models"
	"github.com/zaqimaulana/mycatalog-be/repositories"
)

type OrderService struct {
	orderRepo   *repositories.OrderRepository
	cartRepo    *repositories.CartRepository
	productRepo *repositories.ProductRepository
}

func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo:   repositories.NewOrderRepository(),
		cartRepo:    repositories.NewCartRepository(),
		productRepo: repositories.NewProductRepository(),
	}
}

func (s *OrderService) Checkout(userID uint, req *models.CheckoutRequest) (*models.Order, error) {
	// Ambil cart user
	cartItems, err := s.cartRepo.GetByUserID(userID)
	if err != nil || len(cartItems) == 0 {
		return nil, errors.New("keranjang belanja kosong")
	}

	// Bangun order items + hitung total
	var orderItems []models.OrderItem
	var totalAmount float64

	for _, item := range cartItems {
		product, err := s.productRepo.FindByID(item.ProductID)
		if err != nil {
			return nil, errors.New("produk " + item.Product.Name + " tidak ditemukan")
		}
		if product.Stock < item.Quantity {
			return nil, errors.New("stok produk " + product.Name + " tidak mencukupi")
		}

		subtotal := product.Price * float64(item.Quantity)
		totalAmount += subtotal

		orderItems = append(orderItems, models.OrderItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			Price:       product.Price,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})

		// Kurangi stok
		newStock := product.Stock - item.Quantity
		if err := s.productRepo.UpdateStock(product.ID, newStock); err != nil {
			return nil, err
		}
	}

	// Buat order
	order := &models.Order{
		UserID:          userID,
		Status:          models.OrderStatusPending,
		TotalAmount:     totalAmount,
		ShippingAddress: req.ShippingAddress,
		Notes:           req.Notes,
		Items:           orderItems,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	// Kosongkan cart
	_ = s.cartRepo.ClearCart(userID)

	return order, nil
}

// CreateOrderDirect — buat order dari items yang dikirim Flutter (local cart)
func (s *OrderService) CreateOrderDirect(userID uint, req *models.DirectOrderRequest) (*models.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("tidak ada item dalam pesanan")
	}

	var orderItems []models.OrderItem
	var totalAmount float64

	for _, item := range req.Items {
		product, err := s.productRepo.FindByID(uint(item.ProductID))
		if err != nil {
			return nil, errors.New("produk tidak ditemukan")
		}
		if product.Stock < item.Quantity {
			return nil, errors.New("stok produk " + product.Name + " tidak mencukupi")
		}

		subtotal := product.Price * float64(item.Quantity)
		totalAmount += subtotal

		orderItems = append(orderItems, models.OrderItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			Price:       product.Price,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})

		// Kurangi stok
		if err := s.productRepo.UpdateStock(product.ID, product.Stock-item.Quantity); err != nil {
			return nil, err
		}
	}

	order := &models.Order{
		UserID:           userID,
		Status:           models.OrderStatusPaid,
		TotalAmount:      totalAmount,
		PaymentReference: req.PaymentReference,
		PaymentMethod:    req.PaymentMethod,
		Items:            orderItems,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetMyOrders(userID uint, page, limit int) ([]models.Order, int64, error) {
	return s.orderRepo.GetByUserID(userID, page, limit)
}

func (s *OrderService) GetOrderByID(orderID, userID uint) (*models.Order, error) {
	return s.orderRepo.GetByID(orderID, userID)
}

func (s *OrderService) GetAllOrders(page, limit int) ([]models.Order, int64, error) {
	return s.orderRepo.GetAll(page, limit)
}

func (s *OrderService) UpdateOrderStatus(orderID uint, status models.OrderStatus) error {
	validStatuses := map[models.OrderStatus]bool{
		models.OrderStatusPending:    true,
		models.OrderStatusPaid:       true,
		models.OrderStatusProcessing: true,
		models.OrderStatusShipped:    true,
		models.OrderStatusDelivered:  true,
		models.OrderStatusCancelled:  true,
	}
	if !validStatuses[status] {
		return errors.New("status tidak valid")
	}
	return s.orderRepo.UpdateStatus(orderID, status)
}
