package routes

import (
	
	"github.com/gofiber/fiber/v2"
	"github.com/kavikkannan/go-jwt/pkg/controllers"
	"github.com/github.com/kavikkannan/go-jwt/pkg/"
)

func Setup(app *fiber.App) {

	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.User)
	app.Post("/api/logout", controllers.Logout)

	/* for products */

	app.Get("/products", GetProducts)
	app.Get("/products/:id", GetProductByID)
	app.Get("/products/search", SearchProducts)

	// Admin routes (protected by AdminMiddleware)
	app.Post("/products", AdminMiddleware, AddProduct)
	app.Put("/products/:id", AdminMiddleware, UpdateProduct)

}