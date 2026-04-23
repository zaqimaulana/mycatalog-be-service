package repositories

import (
	"github.com/zaqimaulana/mycatalog-be/config"
	"github.com/zaqimaulana/mycatalog-be/models"
)

type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (r *OrderRepository) Create(order *models.Order) error {
	return config.DB.Create(order).Error
}

func (r *OrderRepository) GetByUserID(userID uint, page, limit int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	offset := (page - 1) * limit
	config.DB.Model(&models.Order{}).Where("user_id = ?", userID).Count(&total)
	err := config.DB.Preload("Items").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&orders).Error
	return orders, total, err
}

func (r *OrderRepository) GetByID(id, userID uint) (*models.Order, error) {
	var order models.Order
	err := config.DB.Preload("Items").
		Where("id = ? AND user_id = ?", id, userID).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetAll untuk admin — semua order
func (r *OrderRepository) GetAll(page, limit int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	offset := (page - 1) * limit
	config.DB.Model(&models.Order{}).Count(&total)
	err := config.DB.Preload("Items").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&orders).Error
	return orders, total, err
}

func (r *OrderRepository) UpdateStatus(id uint, status models.OrderStatus) error {
	return config.DB.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}
