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

var products []models.Product 



// GET /products - Fetch all products
func GetProducts(c *fiber.Ctx) error {
	return c.JSON(products)
}

// GET /products/:id - Fetch product by ID
func GetProductByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid product ID"})
	}

	for _, product := range products {
		if product.ID == uint(id) {
			return c.JSON(product)
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Product not found"})
}

// GET /products/search?query=apple - Search for products
func SearchProducts(c *fiber.Ctx) error {

	

	query := c.Params("query")
	var results []models.Product
	
	
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
	product := new(models.Product)

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
	var product models.Product
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


