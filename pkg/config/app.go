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
	d.AutoMigrate(&models.Login{})
	d.AutoMigrate(&models.Product{})
}

func GetDB() *gorm.DB{
	return DB
}