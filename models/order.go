package models

import "gorm.io/gorm"

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

// Order menyimpan transaksi pembelian
type Order struct {
	gorm.Model
	UserID          uint        `gorm:"not null;index"           json:"user_id"`
	Status          OrderStatus `gorm:"size:20;default:pending"  json:"status"`
	TotalAmount     float64     `gorm:"not null"                 json:"total_amount"`
	ShippingAddress string      `gorm:"type:text"                json:"shipping_address"`
	Notes           string      `gorm:"type:text"                json:"notes"`
	Items           []OrderItem `gorm:"foreignKey:OrderID"       json:"items,omitempty"`
}

// OrderItem menyimpan snapshot produk saat checkout
type OrderItem struct {
	gorm.Model
	OrderID     uint    `gorm:"not null;index"   json:"order_id"`
	ProductID   uint    `gorm:"not null"         json:"product_id"`
	ProductName string  `gorm:"size:200;not null" json:"product_name"`
	Price       float64 `gorm:"not null"         json:"price"`
	Quantity    int     `gorm:"not null"         json:"quantity"`
	Subtotal    float64 `gorm:"not null"         json:"subtotal"`
}

// DTO
type CheckoutRequest struct {
	ShippingAddress string `json:"shipping_address" binding:"required"`
	Notes           string `json:"notes"`
}
