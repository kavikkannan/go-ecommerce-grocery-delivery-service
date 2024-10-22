package config

import(
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kavikkannan/go-ecommerce-grocery-delivery-service/pkg/models"
)

var(
	DB* gorm.DB
)

func Connect(){
	d, err:= gorm.Open("mysql","root:root@tcp(127.0.0.1:3306)/kavi?charset=utf8&parseTime=True&loc=Local")
	if err != nil{
		panic(err)
	}
	DB=d


    // Auto-Migrate the schema
    d.AutoMigrate(&models.Login{}, &models.Products{}, &models.CartItem{}, &models.Order{},
                   &models.OrderItem{}, &models.DeliveryPartner{}, &models.PaymentRequest{},
                   &models.DeliveryAssignment{}, &models.DeliveryStatusUpdate{})

}

func GetDB() *gorm.DB{
	return DB
}