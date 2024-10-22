package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kavikkannan/go-jwt/pkg/config"
	"github.com/kavikkannan/go-jwt/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	/* "net/http" */
	"strconv"
	"time"
)

const SecretKey = "secret"

var products = []Product{
	{ID: 1, Name: "Apple", Category: "Fruits", Price: 1.2, Stock: 50, Description: "Fresh apples"},
	{ID: 2, Name: "Milk", Category: "Dairy", Price: 1.5, Stock: 30, Description: "1 liter fresh milk"},
}

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

var products = []Product{
	{ID: 1, Name: "Apple", Category: "Fruits", Price: 1.2, Stock: 50, Description: "Fresh apples"},
	{ID: 2, Name: "Milk", Category: "Dairy", Price: 1.5, Stock: 30, Description: "1 liter fresh milk"},
}

// JWT Middleware for Admin Access
func AdminMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthenticated"})
	}

	return c.Next()
}

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
	query := c.Query("query")
	var results []Product

	for _, product := range products {
		if contains(product.Name, query) || contains(product.Category, query) {
			results = append(results, product)
		}
	}

	return c.JSON(results)
}

func contains(str, substr string) bool {
	return fiber.Ctx{}.Locals("compare")(str, substr) != ""
}

// POST /products - Add new product (Admin only)
func AddProduct(c *fiber.Ctx) error {
	product := new(Product)

	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	product.ID = uint(len(products) + 1)
	products = append(products, *product)
	return c.JSON(product)
}

// PUT /products/:id - Update product (Admin only)
func UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid product ID"})
	}

	product := new(Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot parse JSON"})
	}

	for i, p := range products {
		if p.ID == uint(id) {
			products[i] = *product
			products[i].ID = uint(id) // Retain original ID
			return c.JSON(products[i])
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Product not found"})
