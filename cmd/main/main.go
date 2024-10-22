package main

import (
	/* "log"
	"net/http" */
    "github.com/gofiber/fiber/v2"
	/* "github.com/rs/cors" */ 
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kavikkannan/go-ecommerce-grocery-delivery-service/pkg/routes"
    "github.com/kavikkannan/go-ecommerce-grocery-delivery-service/pkg/config"
	"github.com/gofiber/fiber/v2/middleware/cors"
	
)

func main() {
    config.Connect()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:"https://busticketbooking-topaz.vercel.app",

	}))
	routes.Setup(app)

	
	app.Listen(":9000")
}
