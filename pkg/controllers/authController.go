package controllers

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/kavikkannan/go-ecommerce-grocery-delivery-service/pkg/config"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"time"
	"fmt"
)

const SecretKey = "secret"

// Register a new user
func Register(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	isAdmin := data["is_admin"] == "true"

	_, err := config.DB.Exec("INSERT INTO Login (name, email, password, is_admin) VALUES (?, ?, ?, ?)", data["name"], data["email"], password, isAdmin)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to register user"})
	}

	return c.JSON(fiber.Map{"message": "User registered successfully"})
}

// Login an existing user
func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid input data"})
	}

	var id int
	var hashedPassword []byte
	var isAdmin bool

	// Query to retrieve user information
	err := config.DB.QueryRow("SELECT id, password, is_admin FROM Login WHERE email = ?", data["email"]).Scan(&id, &hashedPassword, &isAdmin)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	} else if err != nil {
		fmt.Println("Database error:", err) // Log error for debugging
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Database error"})
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(data["password"])); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Incorrect password"})
	}

	// Create JWT claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Issuer":   strconv.Itoa(id),
		"Expires":  time.Now().Add(time.Hour * 24).Unix(),
		"IsAdmin":  isAdmin,
	})
	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		fmt.Println("JWT signing error:", err) // Log JWT error for debugging
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not login"})
	}

	// Set cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		SameSite: "None",
	}
	if c.Protocol() == "https" {
		cookie.Secure = true
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{"message": "Login successful"})
}


// Get User details based on JWT
func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthenticated"})
	}
	claims := token.Claims.(*jwt.StandardClaims)

	var Login struct {
		ID       int
		Name     string
		Email    string
		IsAdmin  bool
	}
	err = config.DB.QueryRow("SELECT id, name, email, is_admin FROM Login WHERE id = ?", claims.Issuer).Scan(&Login.ID, &Login.Name, &Login.Email, &Login.IsAdmin)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Database error"})
	}

	return c.JSON(Login)
}

// Logout by clearing JWT cookie
func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		SameSite: "None",
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}

// Fetch all products
func GetProducts(c *fiber.Ctx) error {
	rows, err := config.DB.Query("SELECT id, name, category, price, stock, description FROM Products")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to retrieve products"})
	}
	defer rows.Close()

	var products []map[string]interface{}
	for rows.Next() {
		var product struct {
			ID       int
			Name     string
			Category string
			Price    float64
			Stock    int
			Description string 
		}
		if err := rows.Scan(&product.ID, &product.Name, &product.Category, &product.Price, &product.Stock, &product.Description); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error scanning product data"})
		}
		products = append(products, map[string]interface{}{
			"id":       product.ID,
			"name":     product.Name,
			"category": product.Category,
			"price":    product.Price,
			"stock":	product.Stock,
			"description": product.Description,
		})
	}

	return c.JSON(products)
}

// Fetch a product by ID
func GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var product struct {
		ID       int
		Name     string
		Category string
		Price    float64
		Stock    int
		Description string 
	}
	err := config.DB.QueryRow("SELECT id, name, category, price, stock, description FROM Products WHERE id = ?", id).Scan(&product.ID, &product.Name, &product.Category, &product.Price, &product.Stock, &product.Description)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Product not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Database error"})
	}

	return c.JSON(product)
}

// Search for products
func SearchProducts(c *fiber.Ctx) error {
	query := c.Query("query")
	rows, err := config.DB.Query("SELECT id, name, category, price, stock, description FROM Products WHERE LOWER(name) LIKE ? OR LOWER(category) LIKE ?", "%"+strings.ToLower(query)+"%", "%"+strings.ToLower(query)+"%")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to retrieve products"})
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var product struct {
			ID       int
			Name     string
			Category string
			Price    float64
			Stock    int
			Description string 
		}
		if err := rows.Scan(&product.ID, &product.Name, &product.Category, &product.Price); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error scanning product data"})
		}
		results = append(results, map[string]interface{}{
			"id":       product.ID,
			"name":     product.Name,
			"category": product.Category,
			"price":    product.Price,
			"stock":	product.Stock,
			"description": product.Description,
		})
	}

	return c.JSON(results)
}




// POST /products - Add new product (Admin only)
func AddProduct(c *fiber.Ctx) error {
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	query := "INSERT INTO Products (name, category, price, stock, description) VALUES (?, ?, ?, ?, ?)"
	_, err := config.DB.Exec(query, data["name"], data["category"], data["price"], data["stock"], data["description"])
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to add product"})
	}

	return c.JSON(fiber.Map{"message": "Product added successfully"})
}




// PUT /products/:id - Update product (Admin only)


// Update an existing product by ID
func UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid product ID"})
	}

	var productData map[string]interface{}
	if err := c.BodyParser(&productData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	// Update product
	_, err = config.DB.Exec("UPDATE Products SET name = ?, category = ?, price = ? WHERE id = ?",
		productData["name"], productData["category"], productData["price"], id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update product"})
	}

	return c.JSON(fiber.Map{"message": "Product updated successfully"})
}

// Add items to the cart
func AddToCart(c *fiber.Ctx) error {
	var cartItem map[string]interface{}
	if err := c.BodyParser(&cartItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	var existingQuantity int
	err := config.DB.QueryRow("SELECT quantity FROM CartItem WHERE user_id = ? AND product_id = ?", cartItem["user_id"], cartItem["product_id"]).Scan(&existingQuantity)

	if err == sql.ErrNoRows {
		// Insert new item if it doesn't exist
		_, err = config.DB.Exec("INSERT INTO CartItem (user_id, product_id, quantity) VALUES (?, ?, ?)",
			cartItem["user_id"], cartItem["product_id"], cartItem["quantity"])
	} else if err == nil {
		// Update quantity if item exists
		_, err = config.DB.Exec("UPDATE CartItem SET quantity = ? WHERE user_id = ? AND product_id = ?",
			existingQuantity+int(cartItem["quantity"].(float64)), cartItem["user_id"], cartItem["product_id"])
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to add item to cart"})
	}

	return c.JSON(fiber.Map{"message": "Item added to cart successfully"})
}

// Fetch cart contents
func GetCart(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("userId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	rows, err := config.DB.Query("SELECT ci.id, ci.user_id, ci.product_id, ci.quantity, p.name, p.price FROM CartItem ci JOIN Products p ON ci.product_id = p.id WHERE ci.user_id = ?", userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to retrieve cart contents"})
	}
	defer rows.Close()

	var cartItems []map[string]interface{}
	for rows.Next() {
		var itemId, userId, productId, quantity int
		var productName string
		var productPrice float64

		// Scan the values into appropriate variables
		err := rows.Scan(&itemId, &userId, &productId, &quantity, &productName, &productPrice)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error scanning cart data"})
		}

		// Assign the scanned values to the map
		item := map[string]interface{}{
			"id":       itemId,
			"user_id":  userId,
			"product_id": productId,
			"quantity": quantity,
			"product": map[string]interface{}{
				"name":  productName,
				"price": productPrice,
			},
		}

		cartItems = append(cartItems, item)
	}

	return c.JSON(cartItems)
}


// Remove a product from the cart
func RemoveFromCart(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	productId, err := strconv.Atoi(c.Params("productId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid product ID"})
	}

	_, err = config.DB.Exec("DELETE FROM CartItem WHERE user_id = ? AND product_id = ?", userId, productId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to remove product from cart"})
	}

	return c.JSON(fiber.Map{"message": "Product removed from cart successfully"})
}

// Place an order
func Checkout(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("userId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	rows, err := config.DB.Query("SELECT ci.product_id, ci.quantity, p.price FROM CartItem ci JOIN Products p ON ci.product_id = p.id WHERE ci.user_id = ?", userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to retrieve cart contents"})
	}
	defer rows.Close()

	var totalAmount float64
	var cartItems []map[string]interface{}
	for rows.Next() {
		var productId, quantity int
		var price float64

		// Scan values into variables
		err := rows.Scan(&productId, &quantity, &price)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error scanning cart data"})
		}

		// Create a map for the item and assign the scanned values
		item := map[string]interface{}{
			"product_id": productId,
			"quantity":   quantity,
			"price":      price,
		}

		// Calculate total amount
		totalAmount += price * float64(quantity)
		cartItems = append(cartItems, item)
	}

	// Insert new order
	res, err := config.DB.Exec("INSERT INTO Orders (user_id, total_amount, status, delivery_status) VALUES (?, ?, ?, ?)", userId, totalAmount, "Pending", "Not Assigned")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create order", "error": err.Error()})
	}

	orderId, err := res.LastInsertId()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to retrieve order ID"})
	}

	// Insert order items
	for _, item := range cartItems {
		_, err := config.DB.Exec("INSERT INTO OrderItem (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)",
			orderId, item["product_id"], item["quantity"], item["price"])
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to add order items"})
		}
	}

	// Clear user's cart
	_, err = config.DB.Exec("DELETE FROM CartItem WHERE user_id = ?", userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to clear cart"})
	}

	return c.JSON(fiber.Map{"message": "Order placed successfully", "order_id": orderId})
}

// Get order details by order ID
func GetOrderDetails(c *fiber.Ctx) error {
	orderId, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid order ID"})
	}

	row := config.DB.QueryRow("SELECT id, user_id, total_amount, status FROM Orders WHERE id = ?", orderId)
	var order struct {
		ID          int
		UserID      int
		TotalAmount float64
		Status      string
	}
	if err := row.Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.Status); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Order not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Database error"})
	}

	rows, err := config.DB.Query("SELECT product_id, quantity, price FROM OrderItem WHERE order_id = ?", orderId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to retrieve order items"})
	}
	defer rows.Close()

	var items []map[string]interface{}
	for rows.Next() {
		var productId, quantity int
		var price float64

		// Scan values into variables
		err := rows.Scan(&productId, &quantity, &price)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error scanning order items"})
		}

		// Create a map for the item and assign the scanned values
		item := map[string]interface{}{
			"product_id": productId,
			"quantity":   quantity,
			"price":      price,
		}

		items = append(items, item)
	}

	return c.JSON(fiber.Map{"order": order, "items": items})
}


// Initiate payment for an order
func InitiatePayment(c *fiber.Ctx) error {
	var paymentRequest map[string]int
	if err := c.BodyParser(&paymentRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	_, err := config.DB.Exec("UPDATE Orders SET status = 'Paid' WHERE id = ?", paymentRequest["order_id"])
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update order status"})
	}

	return c.JSON(fiber.Map{"message": "Payment initiated successfully"})
}




// GET /orders/:orderId/track - Track order delivery status


func TrackOrder(c *fiber.Ctx) error {
	orderId, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid order ID"})
	}

	var deliveryStatus string
	err = config.DB.QueryRow("SELECT delivery_status FROM Orders WHERE id = ?", orderId).Scan(&deliveryStatus)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Order not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Database error"})
	}

	return c.JSON(fiber.Map{"status": deliveryStatus})
}

// Assign delivery partner
func AssignDeliveryPartner(c *fiber.Ctx) error {
	orderId, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid order ID"})
	}

	var assignment struct {
		PartnerID int `json:"partner_id"`
	}
	if err := c.BodyParser(&assignment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	_, err = config.DB.Exec("UPDATE Orders SET delivery_partner_id = ? WHERE id = ?", assignment.PartnerID, orderId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to assign delivery partner"})
	}

	return c.JSON(fiber.Map{"message": "Delivery partner assigned successfully"})
}

// Update delivery status
func UpdateDeliveryStatus(c *fiber.Ctx) error {
	var statusUpdate struct {
		ID     int    `json:"id"`
		Status string `json:"status"`
	}

	if err := c.BodyParser(&statusUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	_, err := config.DB.Exec("UPDATE Orders SET delivery_status = ? WHERE id = ?", statusUpdate.Status, statusUpdate.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update delivery status"})
	}

	return c.JSON(fiber.Map{"message": "Delivery status updated successfully"})
}

// Admin Dashboard
func AdminDashboard(c *fiber.Ctx) error {
	var totalOrders int64
	var activeUsers int64
	/* var avgDeliveryTime float64 */

	// Total orders
	err := config.DB.QueryRow("SELECT COUNT(*) FROM orders").Scan(&totalOrders)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error retrieving total orders"})
	}

	// Active users (non-admins)
	err = config.DB.QueryRow("SELECT COUNT(*) FROM Login WHERE is_admin = 0").Scan(&activeUsers)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error retrieving active users"})
	}

	// Top-selling products
	var topProducts []struct {
		ProductID int
		TotalSold int
	}
	rows, err := config.DB.Query(`
		SELECT product_id, SUM(quantity) AS total_sold
		FROM OrderItem
		JOIN Products ON Products.id = OrderItem.product_id
		GROUP BY product_id
		ORDER BY total_sold DESC
		LIMIT 5
	`)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error retrieving top-selling products"})
	}
	defer rows.Close()
	for rows.Next() {
		var product struct {
			ProductID int
			TotalSold int
		}
		if err := rows.Scan(&product.ProductID, &product.TotalSold); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error scanning product data"})
		}
		topProducts = append(topProducts, product)
	}

	// Average delivery time
	/* err = config.DB.QueryRow(`
    SELECT AVG((julianday(updated_at) - julianday(created_at)) * 24 * 60)
    FROM Orders WHERE status = 'Delivered'
`).Scan(&avgDeliveryTime)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error calculating average delivery time"})
	} */

	return c.JSON(fiber.Map{
		"total_orders":     totalOrders,
		"active_users":     activeUsers,
		"top_selling":      topProducts,
		
	})
}

// Get All Orders
func GetAllOrders(c *fiber.Ctx) error {
	var orders []map[string]interface{}
	status := c.Query("status")

	query := "SELECT id, status, delivery_partner_id FROM Orders"
	if status != "" {
		query += " WHERE status = ?"
	}

	rows, err := config.DB.Query(query, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to retrieve orders"})
	}
	defer rows.Close()

	for rows.Next() {
		var order struct {
			ID               int
			Status           string
			DeliveryPartnerID int
		}
		if err := rows.Scan(&order.ID, &order.Status, &order.DeliveryPartnerID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error scanning order data"})
		}
		orders = append(orders, map[string]interface{}{
			"id":               order.ID,
			"status":           order.Status,
			"delivery_partner": order.DeliveryPartnerID,
		})
	}

	return c.JSON(orders)
}

// Cancel Order
func CancelOrder(c *fiber.Ctx) error {
	orderId, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid order ID"})
	}

	// Update order status to "Canceled"
	_, err = config.DB.Exec("UPDATE Orders SET status = 'Canceled' WHERE id = ?", orderId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to cancel order"})
	}

	return c.JSON(fiber.Map{"message": "Order canceled successfully"})
}
