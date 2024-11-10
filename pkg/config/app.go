package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
	


// Connect opens the SQLite database and creates tables if they don't exist
func Connect() {
	var err error
	DB, err = sql.Open("sqlite3", "./example.db")
	if err != nil {
		panic(err)
	}


	// Create tables if they don't exist
	err = createTables(DB)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
}

// createTables runs SQL statements to create each table if it doesn't exist
func createTables(db *sql.DB) error {
	tableStatements := []string{
		`CREATE TABLE IF NOT EXISTS Products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			category TEXT NOT NULL,
			price REAL NOT NULL,
			stock INTEGER NOT NULL,
			description TEXT
		);`,







		
		`CREATE TABLE IF NOT EXISTS Login (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password BLOB NOT NULL,
			is_admin BOOLEAN NOT NULL DEFAULT 0
		);`,

		

		`CREATE TABLE IF NOT EXISTS CartItem (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			quantity INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES Login(id),
			FOREIGN KEY (product_id) REFERENCES Products(id)
		);`,

		`CREATE TABLE IF NOT EXISTS Orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			total_amount REAL NOT NULL,
			status TEXT NOT NULL,
			delivery_status TEXT NOT NULL,
			delivery_partner_id INTEGER,
			FOREIGN KEY (user_id) REFERENCES Login(id),
			FOREIGN KEY (delivery_partner_id) REFERENCES DeliveryPartner(id)
		);`,

		`CREATE TABLE IF NOT EXISTS DeliveryPartner (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS OrderItem (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			quantity INTEGER NOT NULL,
			price REAL NOT NULL,
			FOREIGN KEY (order_id) REFERENCES Orders(id),
			FOREIGN KEY (product_id) REFERENCES Products(id)
		);`,

		`CREATE TABLE IF NOT EXISTS PaymentRequest (
			order_id INTEGER NOT NULL,
			amount REAL NOT NULL,
			method TEXT NOT NULL,
			FOREIGN KEY (order_id) REFERENCES Orders(id)
		);`,

		`CREATE TABLE IF NOT EXISTS DeliveryAssignment (
			partner_id INTEGER NOT NULL,
			order_id INTEGER NOT NULL,
			FOREIGN KEY (partner_id) REFERENCES DeliveryPartner(id),
			FOREIGN KEY (order_id) REFERENCES Orders(id)
		);`,

		`CREATE TABLE IF NOT EXISTS DeliveryStatusUpdate (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_id INTEGER NOT NULL,
			status TEXT NOT NULL,
			FOREIGN KEY (order_id) REFERENCES Orders(id)
		);`,
	}

	// Execute each table creation statement
	for _, stmt := range tableStatements {
		_, err := db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
	}

	return nil
}

// GetDB returns a pointer to the database connection
func GetDB() *sql.DB {
	return DB
}
