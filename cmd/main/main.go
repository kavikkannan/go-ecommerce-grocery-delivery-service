package main

import (
	/* "log"
	"net/http" */


	"github.com/gofiber/fiber/v2"
	/* "github.com/rs/cors" */
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kavikkannan/go-ecommerce-grocery-delivery-service/pkg/config"
	"github.com/kavikkannan/go-ecommerce-grocery-delivery-service/pkg/routes"
	



	_ "github.com/mattn/go-sqlite3"
)

func main() {

    config.Connect()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:3000", // Use "http" if your frontend is on HTTP
	}))
	routes.Setup(app)

	
	app.Listen(":9000")
}




/* // Open the database
database, err := sql.Open("sqlite3", "./example.db")
if err != nil {
	log.Fatal(err)
}
defer database.Close()

// Create the "people" table
statement, err := database.Prepare(`
	CREATE TABLE IF NOT EXISTS people (
		id INTEGER PRIMARY KEY,
		firstname TEXT,
		lastname TEXT
	)
`)
if err != nil {
	log.Fatal(err)
}
_, err = statement.Exec()
if err != nil {
	log.Fatal(err)
}

// Create the "orders" table with a foreign key to the "people" table
statement, err = database.Prepare(`
	CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY,
		order_date TEXT,
		amount REAL,
		person_id INTEGER,
		FOREIGN KEY(person_id) REFERENCES people(id)
	)
`)
if err != nil {
	log.Fatal(err)
}
_, err = statement.Exec()
if err != nil {
	log.Fatal(err)
}

// Insert data into the "people" table
statement, err = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
if err != nil {
	log.Fatal(err)
}
_, err = statement.Exec("John", "Doe")
if err != nil {
	log.Fatal(err)
}
_, err = statement.Exec("Jane", "Smith")
if err != nil {
	log.Fatal(err)
}

// Insert data into the "orders" table with foreign key reference to "people"
statement, err = database.Prepare("INSERT INTO orders (order_date, amount, person_id) VALUES (?, ?, ?)")
if err != nil {
	log.Fatal(err)
}
_, err = statement.Exec("2024-11-01", 250.75, 1) // John Doe's order
if err != nil {
	log.Fatal(err)
}
_, err = statement.Exec("2024-11-02", 125.50, 2) // Jane Smith's order
if err != nil {
	log.Fatal(err)
}

// Query to fetch orders along with the customer name (using the foreign key relationship)
rows, err := database.Query(`
	SELECT orders.id, orders.order_date, orders.amount, people.firstname, people.lastname
	FROM orders
	INNER JOIN people ON orders.person_id = people.id
`)
if err != nil {
	log.Fatal(err)
}
defer rows.Close()

// Display the orders with customer names
for rows.Next() {
	var orderID int
	var orderDate string
	var amount float64
	var firstname, lastname string
	err = rows.Scan(&orderID, &orderDate, &amount, &firstname, &lastname)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Order ID: %d, Date: %s, Amount: %.2f, Customer: %s %s\n", orderID, orderDate, amount, firstname, lastname)
}

// Check for any errors during iteration
if err = rows.Err(); err != nil {
	log.Fatal(err)
} */