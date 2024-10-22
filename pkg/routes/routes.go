package routes

import (
	
	"github.com/gofiber/fiber/v2"
	"github.com/kavikkannan/go-ecommerce-grocery-delivery-service/pkg/controllers"
	"github.com/kavikkannan/go-ecommerce-grocery-delivery-service/pkg/middleware"
)

func Setup(app *fiber.App) {

	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.User)
	app.Post("/api/logout", controllers.Logout)

	/* for products */

	app.Get("/products", controllers.GetProducts)
	app.Get("/products/:id", controllers.GetProductByID)
	app.Get("/products/search/:query", controllers.SearchProducts)

	// Admin routes (protected by AdminMiddleware)
	app.Post("/products", AdminMiddlewareAccess.AdminMiddleware, controllers.AddProduct)
	app.Put("/products/:id", AdminMiddlewareAccess.AdminMiddleware, controllers.UpdateProduct)

}