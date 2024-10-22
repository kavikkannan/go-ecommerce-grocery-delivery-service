package models

import "time"

type Login struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password []byte `json:"-"`
	IsAdmin  bool   `json:"is_admin"` // New field to identify if the user is an admin
}


type Products struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Price       float64 `json:"price"`
	Stock       int    `json:"stock"`
	Description string `json:"description"`
}

type CartItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	UserID    uint    `json:"user_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Product   Products `gorm:"foreignKey:ProductID"` // Relation to the Products table
}

type Order struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	UserID         uint           `json:"user_id"`
	TotalAmount    float64        `json:"total_amount"`
	Status         string         `json:"status"` // e.g., "Pending", "Paid", "Shipped"
	DeliveryStatus string         `json:"delivery_status"` // e.g., "Out for Delivery", "Delivered"
	DeliveryPartnerID uint        `json:"delivery_partner_id"` // Foreign key to a DeliveryPartner model
	Items          []OrderItem    `json:"items" gorm:"foreignKey:OrderID"` // List of ordered items
	CreatedAt      time.Time      `json:"created_at"`
}

type OrderItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"` // Unit price at the time of order
	Product   Products `gorm:"foreignKey:ProductID"` // Relation to the Products table
}

type PaymentRequest struct {
	OrderID  uint    `json:"order_id"`
	Amount   float64 `json:"amount"`
	Method   string  `json:"method"` // e.g., "Credit Card", "PayPal"
}

type DeliveryAssignment struct {
	PartnerID uint `json:"partner_id"`
	OrderID   uint `json:"order_id"`
}

type DeliveryStatusUpdate struct {
	OrderID uint   `json:"order_id"`
	Status  string `json:"status"` // e.g., "Picked Up", "On the Way", "Delivered"
}
