package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kavikkannan/go-ecommerce-grocery-delivery-service/pkg/controllers"
	AdminMiddlewareAccess "github.com/kavikkannan/go-ecommerce-grocery-delivery-service/pkg/middleware"
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

	// Cart Routes
	app.Post("/cart/add", controllers.AddToCart)               // Add items to the user's shopping cart
	app.Get("/cart", controllers.GetCart)                      // Fetch the current contents of the user's shopping cart
	app.Delete("/cart/:productId", controllers.RemoveFromCart) // Remove a product from the cart

	// Order Routes
	app.Post("/orders/checkout", controllers.Checkout)       // Place an order for the items in the cart
	app.Get("/orders/:orderId", controllers.GetOrderDetails) // Fetch details of a specific order

	// Delivery Management Routes
	app.Get("/orders/:orderId/track", controllers.TrackOrder)                                                     // Track the real-time delivery status of an order
	app.Post("/orders/:orderId/assign", AdminMiddlewareAccess.AdminMiddleware, controllers.AssignDeliveryPartner) // Assign a delivery partner to an order
	app.Post("/delivery/update-status", controllers.UpdateDeliveryStatus)
}
