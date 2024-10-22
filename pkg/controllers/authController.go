package controllers

import (


	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/kavikkannan/go-ecommerce-grocery-delivery-service/pkg/config"
	"github.com/kavikkannan/go-ecommerce-grocery-delivery-service/pkg/models"
	"golang.org/x/crypto/bcrypt"

	/* "net/http" */
	"strconv"
	"strings"
	"time"
	
)

const SecretKey = "secret"


func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	isAdmin := false
	if data["is_admin"] == "true" {
		isAdmin = true
	}

	user := models.Login{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
		IsAdmin:  isAdmin, // Set whether the user is an admin
	}

	config.DB.Create(&user)

	return c.JSON(user)
}


func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.Login

	config.DB.Where("email = ?", data["email"]).First(&user)

	if user.ID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Issuer":   strconv.Itoa(int(user.ID)),
		"Expires":  time.Now().Add(time.Hour * 24).Unix(), //1 day
		"IsAdmin":  user.IsAdmin,                          // Include IsAdmin flag in the token
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

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

	return c.JSON(fiber.Map{
		"message": "success",
	})
}


func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.Login

	config.DB.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		SameSite: "None",
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
		
	})
}


/* for product page: */

// GET /products - Fetch all products
func GetProducts(c *fiber.Ctx) error {
	var products []models.Products

	// Fetch all products from the database
	if err := config.DB.Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to retrieve products"})
	}

	return c.JSON(products)
}

// GET /products/:id - Fetch product by ID
func GetProductByID(c *fiber.Ctx) error {
	// Parse the product ID from the URL parameter
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid product ID"})
	}

	// Retrieve the product from the database by ID
	var product models.Products
	if err := config.DB.First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Product not found"})
	}

	return c.JSON(product)
}

// GET /products/search?query=apple - Search for products
func SearchProducts(c *fiber.Ctx) error {
	// Get the query parameter from the URL
	query := c.Query("query")
	var products []models.Products
	var results []models.Products

	// Fetch all products from the database
	if err := config.DB.Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to retrieve products"})
	}

	// Filter the products based on name or category (case-insensitive)
	for _, product := range products {
		if contains(product.Name, query) || contains(product.Category, query) {
			results = append(results, product)
		}
	}

	return c.JSON(results)
}


func contains(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr)) // case insensitive search
}

// POST /products - Add new product (Admin only)
func AddProduct(c *fiber.Ctx) error {
	// Create a new instance of the product model
	product := new(models.Products)

	// Parse the request body into the product struct
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	// Save the product to the database using GORM
	if err := config.DB.Create(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to add product"})
	}

	// Return the newly created product as a JSON response
	return c.JSON(product)
}



// PUT /products/:id - Update product (Admin only)
func UpdateProduct(c *fiber.Ctx) error {
	// Parse the product ID from the URL parameter
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid product ID"})
	}

	// Retrieve the product from the database
	var product models.Products
	if err := config.DB.First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Product not found"})
	}

	// Parse the request body into the product struct
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	// Update the product in the database
	if err := config.DB.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update product"})
	}

	return c.JSON(product)
}


/* Cart management */

// POST /cart/add - Add items to the cart
func AddToCart(c *fiber.Ctx) error {
	var cartItem models.CartItem

	// Parse the request body to get product ID and quantity
	if err := c.BodyParser(&cartItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	// Add the item to the user's cart in the database
	if err := config.DB.Create(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to add item to cart"})
	}

	return c.JSON(cartItem)
}

// GET /cart - Get current cart items
func GetCart(c *fiber.Ctx) error {
	var cart []models.CartItem

	// Retrieve the cart items from the database
	if err := config.DB.Find(&cart).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to retrieve cart items"})
	}

	return c.JSON(cart)
}

// DELETE /cart/:productId - Remove product from the cart
func RemoveFromCart(c *fiber.Ctx) error {
	productId, err := strconv.Atoi(c.Params("productId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid product ID"})
	}

	// Delete the item from the cart in the database
	if err := config.DB.Where("product_id = ?", productId).Delete(&models.CartItem{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to remove item from cart"})
	}

	return c.JSON(fiber.Map{"message": "Item removed from cart"})
}

/* Order management */

// POST /orders/checkout - Place an order
func Checkout(c *fiber.Ctx) error {
	var order models.Order

	// Parse the request body to get order details
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	// Save the order to the database
	if err := config.DB.Create(&order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to place order"})
	}

	// Clear the cart
	config.DB.Where("user_id = ?", order.UserID).Delete(&models.CartItem{})

	return c.JSON(order)
}

// GET /orders/:orderId - Get order details
func GetOrderDetails(c *fiber.Ctx) error {
	orderId, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid order ID"})
	}

	var order models.Order

	// Retrieve the order from the database
	if err := config.DB.Preload("Items").First(&order, orderId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Order not found"})
	}

	return c.JSON(order)
}

/* Payment Integration*/

// POST /payments/initiate - Initiate payment
func InitiatePayment(c *fiber.Ctx) error {
	var paymentRequest models.PaymentRequest

	// Parse the payment details
	if err := c.BodyParser(&paymentRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	// Call the payment gateway API (e.g., Stripe, PayPal)
	// Handle the payment logic here

	// Update order status in the database
	config.DB.Model(&models.Order{}).Where("id = ?", paymentRequest.OrderID).Update("status", "Paid")

	return c.JSON(fiber.Map{"message": "Payment initiated successfully"})
}


/* Delivery management */

// GET /orders/:orderId/track - Track order delivery status
func TrackOrder(c *fiber.Ctx) error {
	orderId, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid order ID"})
	}

	var order models.Order

	// Retrieve the order delivery status
	if err := config.DB.First(&order, orderId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Order not found"})
	}

	return c.JSON(fiber.Map{"status": order.DeliveryStatus})
}

// POST /orders/:orderId/assign - Assign delivery partner
func AssignDeliveryPartner(c *fiber.Ctx) error {
	orderId, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid order ID"})
	}

	// Parse the delivery partner assignment details
	var assignment models.DeliveryAssignment
	if err := c.BodyParser(&assignment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	// Assign the delivery partner to the order in the database
	if err := config.DB.Model(&models.Order{}).Where("id = ?", orderId).Update("delivery_partner_id", assignment.PartnerID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to assign delivery partner"})
	}

	return c.JSON(fiber.Map{"message": "Delivery partner assigned successfully"})
}

// POST /delivery/update-status - Update delivery status
func UpdateDeliveryStatus(c *fiber.Ctx) error {
	var statusUpdate models.DeliveryStatusUpdate

	// Parse the status update details
	if err := c.BodyParser(&statusUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	// Update the delivery status in the database
	if err := config.DB.Model(&models.Order{}).Where("id = ?", statusUpdate.OrderID).Update("delivery_status", statusUpdate.Status).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update delivery status"})
	}

	return c.JSON(fiber.Map{"message": "Delivery status updated successfully"})
}
