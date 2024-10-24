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
	app.Post("/cart/add", controllers.AddToCart)               
	app.Get("/cart", controllers.GetCart)                      
	app.Delete("/cart/:productId", controllers.RemoveFromCart) 

	// Order Routes
	app.Post("/orders/checkout", controllers.Checkout)       
	app.Get("/orders/:orderId", controllers.GetOrderDetails) 
	app.Post("/payment/initiate", controllers.InitiatePayment) 
	
	// Delivery Management Routes
	app.Get("/orders/:orderId/track", controllers.TrackOrder)                                                     
	app.Post("/orders/:orderId/assign", AdminMiddlewareAccess.AdminMiddleware, controllers.AssignDeliveryPartner) 
	app.Post("/delivery/update-status", controllers.UpdateDeliveryStatus)

	// Admin routes
	app.Get("/admin/dashboard", controllers.AdminDashboard)
	app.Get("/admin/orders", controllers.GetAllOrders)
	app.Put("/admin/orders/:orderId/cancel", controllers.CancelOrder)

}

