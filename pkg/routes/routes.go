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
	app.Get("/api/user/:userId", controllers.GetUserByID)
	app.Post("/api/logout", controllers.Logout)

	/* for products */

	app.Get("/products", controllers.GetProducts)
	app.Get("/products/:id", controllers.GetProductByID)
	app.Get("/products/search/:query", controllers.SearchProducts)

	// Admin routes (protected by AdminMiddleware)
	app.Post("/products",  controllers.AddProduct)
	app.Put("/products/:id", controllers.UpdateProduct)
	app.Post("/products/remove/:id", controllers.DeleteProduct)

	// Cart Routes
	app.Post("/cart/add", controllers.AddToCart)               
	app.Get("/cart/:userId", controllers.GetCart)                      
	app.Delete("/cart/:productId", controllers.RemoveFromCart) 

	// Order Routes
	app.Post("/orders/checkout/:userId", controllers.Checkout)    
	app.Get("/ordersIds/:userId", controllers.GetOrderIdsByUserID)   
	app.Get("/orders/:orderId", controllers.GetOrderDetails) 
	app.Get("/orders", controllers.GetAllOrderIds) 
	app.Post("/payment/initiate", controllers.InitiatePayment) 
	
	// Delivery Management Routes
	app.Get("/orders/:orderId/track", controllers.TrackOrder)                                                     
	app.Post("/orders/:orderId/assign", controllers.AssignDeliveryPartner) 
	app.Post("/delivery/update-status", controllers.UpdateDeliveryStatus)

	// Admin routes
	app.Get("/admin/dashboard", controllers.AdminDashboard)
	app.Get("/admin/orders", controllers.GetAllOrders)
	app.Put("/admin/orders/:orderId/cancel", controllers.CancelOrder)
	app.Put("/admin",AdminMiddlewareAccess.AdminMiddleware)

}

