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
		IsAdmin:  isAdmin, 
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
		"Expires":  time.Now().Add(time.Hour * 24).Unix(), 
		"IsAdmin":  user.IsAdmin,                          
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

	if err := config.DB.Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to retrieve products"})
	}

	return c.JSON(products)
}

// GET /products/:id - Fetch product by ID
func GetProductByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid product ID"})
	}

	var product models.Products
	if err := config.DB.First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Product not found"})
	}

	return c.JSON(product)
}

// GET /products/search?query=apple - Search for products
func SearchProducts(c *fiber.Ctx) error {
	query := c.Query("query")
	var products []models.Products
	var results []models.Products

	if err := config.DB.Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to retrieve products"})
	}

	for _, product := range products {
		if contains(product.Name, query) || contains(product.Category, query) {
			results = append(results, product)
		}
	}

	return c.JSON(results)
}


func contains(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr)) 
}

// POST /products - Add new product (Admin only)
func AddProduct(c *fiber.Ctx) error {
	product := new(models.Products)

	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	if err := config.DB.Create(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to add product"})
	}

	return c.JSON(product)
}



// PUT /products/:id - Update product (Admin only)
func UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid product ID"})
	}

	var product models.Products
	if err := config.DB.First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Product not found"})
	}

	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	if err := config.DB.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update product"})
	}

	return c.JSON(product)
}


/* Cart management 

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

/* Order management 

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

	// Preload "Items" and their associated "Product"
	if err := config.DB.Preload("Items.Product").First(&order, orderId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Order not found"})
	}

	return c.JSON(order)
}




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


Delivery management 

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
// POST /orders/:orderId/assign - Assign delivery partner
func AssignDeliveryPartner(c *fiber.Ctx) error {
	orderId, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid order ID"})
	}

	var assignment models.DeliveryAssignment
	if err := c.BodyParser(&assignment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	// Assign the delivery partner
	if err := config.DB.Model(&models.Order{}).Where("id = ?", orderId).Update("delivery_partner_id", assignment.PartnerID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to assign delivery partner"})
	}

	return c.JSON(fiber.Map{"message": "Delivery partner assigned successfully"})
}


// POST /delivery/update-status - Update delivery status
// POST /delivery/update-status - Update delivery status
func UpdateDeliveryStatus(c *fiber.Ctx) error {
	var statusUpdate models.DeliveryStatusUpdate

	if err := c.BodyParser(&statusUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	// Update the delivery status
	if err := config.DB.Model(&models.Order{}).Where("id = ?", statusUpdate.ID).Update("delivery_status", statusUpdate.Status).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update delivery status"})
	}

	return c.JSON(fiber.Map{"message": "Delivery status updated successfully"})
}

 */


// POST /cart/add - Add items to cart
func AddToCart(c *fiber.Ctx) error {
	var cartItem models.CartItem

	if err := c.BodyParser(&cartItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	var existingCartItem models.CartItem
	if err := config.DB.Where("user_id = ? AND product_id = ?", cartItem.UserID, cartItem.ProductID).First(&existingCartItem).Error; err == nil {
		existingCartItem.Quantity += cartItem.Quantity
		config.DB.Save(&existingCartItem)
	} else {
		config.DB.Create(&cartItem)
	}

	return c.JSON(fiber.Map{"message": "Item added to cart successfully"})
}
// GET /cart - Fetch cart contents
func GetCart(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	var cartItems []models.CartItem
	if err := config.DB.Preload("Product").Where("user_id = ?", userId).Find(&cartItems).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Cart is empty"})
	}

	return c.JSON(cartItems)
}

 // DELETE /cart/:productId - Remove a product from the cart
func RemoveFromCart(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	productId, err := strconv.Atoi(c.Params("productId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid product ID"})
	}

	if err := config.DB.Where("user_id = ? AND product_id = ?", userId, productId).Delete(&models.CartItem{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to remove product from cart"})
	}

	return c.JSON(fiber.Map{"message": "Product removed from cart successfully"})
}
// POST /orders/checkout - Place an order
func Checkout(c *fiber.Ctx) error {
	var order models.Order

	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	var cartItems []models.CartItem
	if err := config.DB.Preload("Product").Where("user_id = ?", userId).Find(&cartItems).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Cart is empty"})
	}

	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += item.Product.Price * float64(item.Quantity)
	}

	order.UserID = uint(userId)
	order.TotalAmount = totalAmount
	order.Status = "Pending"
	config.DB.Create(&order)

	for _, item := range cartItems {
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Product.Price,
		}
		config.DB.Create(&orderItem)
	}

	config.DB.Where("user_id = ?", userId).Delete(&models.CartItem{})

	return c.JSON(order)
}

func GetOrderDetails(c *fiber.Ctx) error {
	orderId, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid order ID"})
	}

	var order models.Order

	if err := config.DB.Preload("Items.Product").First(&order, orderId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Order not found"})
	}

	return c.JSON(order)
}




// POST /payments/initiate - Initiate payment
func InitiatePayment(c *fiber.Ctx) error {
	var paymentRequest models.PaymentRequest

	if err := c.BodyParser(&paymentRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	config.DB.Model(&models.Order{}).Where("id = ?", paymentRequest.OrderID).Update("status", "Paid")

	return c.JSON(fiber.Map{"message": "Payment initiated successfully"})
}



// GET /orders/:orderId/track - Track order delivery status
func TrackOrder(c *fiber.Ctx) error {
	orderId, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid order ID"})
	}

	var order models.Order

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

	var assignment models.DeliveryAssignment
	if err := c.BodyParser(&assignment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	if err := config.DB.Model(&models.Order{}).Where("id = ?", orderId).Update("delivery_partner_id", assignment.PartnerID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to assign delivery partner"})
	}

	return c.JSON(fiber.Map{"message": "Delivery partner assigned successfully"})
}
// POST /delivery/update-status - Update delivery status
func UpdateDeliveryStatus(c *fiber.Ctx) error {
	var statusUpdate models.DeliveryStatusUpdate

	if err := c.BodyParser(&statusUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	if err := config.DB.Model(&models.Order{}).Where("id = ?", statusUpdate.ID).Update("delivery_status", statusUpdate.Status).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update delivery status"})
	}

	return c.JSON(fiber.Map{"message": "Delivery status updated successfully"})
}

func AdminDashboard(c *fiber.Ctx) error {
    var totalOrders int64
    var activeUsers int64
    var topProducts []models.Products
    var avgDeliveryTime float64

    config.DB.Model(&models.Order{}).Count(&totalOrders)

    config.DB.Model(&models.Login{}).Where("is_admin = ?", false).Count(&activeUsers)

    config.DB.Model(&models.OrderItem{}).
        Select("product_id, sum(quantity) as total_sold").
        Joins("JOIN products on products.id = order_items.product_id").
        Group("product_id").
        Order("total_sold desc").
        Limit(5).
        Scan(&topProducts)

    config.DB.Model(&models.Order{}).
        Where("status = ?", "Delivered").
        Select("avg(timestampdiff(MINUTE, created_at, updated_at)) as avg_delivery_time").
        Row().Scan(&avgDeliveryTime)

    return c.JSON(fiber.Map{
        "total_orders":      totalOrders,
        "active_users":      activeUsers,
        "top_selling":       topProducts,
        "average_delivery":  avgDeliveryTime,
    })
}

func GetAllOrders(c *fiber.Ctx) error {
    var orders []models.Order
    status := c.Query("status") 

    query := config.DB.Preload("Items.Product").Preload("DeliveryPartner")

    if status != "" {
        query = query.Where("status = ?", status)
    }

    if err := query.Find(&orders).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Orders not found"})
    }

    return c.JSON(orders)
}

func CancelOrder(c *fiber.Ctx) error {
    orderId, err := strconv.Atoi(c.Params("orderId"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid order ID"})
    }

    var order models.Order
    if err := config.DB.First(&order, orderId).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Order not found"})
    }

    if err := config.DB.Model(&order).Update("status", "Canceled").Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to cancel order"})
    }

    return c.JSON(fiber.Map{"message": "Order canceled successfully"})
}
