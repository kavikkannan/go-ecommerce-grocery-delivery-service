package models

type Login struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password []byte `json:"-"`
	IsAdmin  bool   `json:"is_admin"` // New field to identify if the user is an admin
}


type Product struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Price       float64 `json:"price"`
	Stock       int    `json:"stock"`
	Description string `json:"description"`
}