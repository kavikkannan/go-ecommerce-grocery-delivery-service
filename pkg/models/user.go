package models

type Login struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password []byte `json:"-"`   // This field will not be returned in JSON responses
	IsAdmin  bool   `json:"is_admin"` // To identify if the user is an admin
}



type Products struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Description string  `json:"description"`
}



type CartItem struct {
	ID        uint     `json:"id" gorm:"primaryKey"`
	UserID    uint     `json:"user_id"` // Foreign key to Users table
	User      Login    `json:"user" gorm:"foreignKey:UserID"` // Define the relation with Login (Users table)
	ProductID uint     `json:"product_id"` // Foreign key to Products table
	Product   Products `json:"product" gorm:"foreignKey:ProductID"` // Define the relation with Products table
	Quantity  int      `json:"quantity"`
}


type Order struct {
	ID                uint            `json:"id" gorm:"primaryKey"`
	UserID            uint            `json:"user_id"` // Foreign key to Users table
	User              Login           `json:"user" gorm:"foreignKey:UserID"` // Relation to Login (Users)
	TotalAmount       float64         `json:"total_amount"`
	Status            string          `json:"status"`
	DeliveryStatus    string          `json:"delivery_status"`
	DeliveryPartnerID uint            `json:"delivery_partner_id"` // Foreign key to DeliveryPartner table
	DeliveryPartner   DeliveryPartner `json:"delivery_partner" gorm:"foreignKey:DeliveryPartnerID"` // Relation to DeliveryPartner
	Items             []OrderItem     `json:"items"` // One-to-many relationship with OrderItem
}

type DeliveryPartner struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}


type OrderItem struct {
	ID        uint     `json:"id" gorm:"primaryKey"`
	OrderID   uint     `json:"order_id"` // Foreign key to Orders table
	Order     Order    `json:"order" gorm:"foreignKey:OrderID"` // Relation to Order
	ProductID uint     `json:"product_id"` // Foreign key to Products table
	Product   Products `json:"product" gorm:"foreignKey:ProductID"` // Relation to Products
	Quantity  int      `json:"quantity"`
	Price     float64  `json:"price"`
}


type PaymentRequest struct {
	OrderID uint    `json:"order_id"` // Foreign key to Orders table
	Amount  float64 `json:"amount"`
	Method  string  `json:"method"` // e.g., "Credit Card", "PayPal"
}


type DeliveryAssignment struct {
	PartnerID uint `json:"partner_id"` // Foreign key to DeliveryPartner table
	OrderID   uint `json:"order_id"`   // Foreign key to Orders table
}


type DeliveryStatusUpdate struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	OrderID uint   `json:"order_id"` // Foreign key to Orders table
	Status  string `json:"status"`   // e.g., "Picked Up", "On the Way", "Delivered"
}

