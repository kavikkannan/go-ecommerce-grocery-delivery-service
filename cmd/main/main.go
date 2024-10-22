package main

import (
	/* "log"
	"net/http" */
	"log"

	"github.com/gofiber/fiber/v2"
	/* "github.com/rs/cors" */
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kavikkannan/go-ecommerce-grocery-delivery-service/pkg/config"
	"github.com/kavikkannan/go-ecommerce-grocery-delivery-service/pkg/routes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	InitDB()
    config.Connect()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:"https://busticketbooking-topaz.vercel.app",

	}))
	routes.Setup(app)

	
	app.Listen(":9000")
}


var db *gorm.DB

func InitDB() {
    var err error
    dsn := "root:root@tcp(127.0.0.1:3306)/kavi?charset=utf8&parseTime=True&loc=Local"
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("failed to connect to database:", err)
    }
}