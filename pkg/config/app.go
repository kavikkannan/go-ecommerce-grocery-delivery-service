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
	d, err:= gorm.Open("mysql","admin:subi1234@tcp(database-1.c5akcsdv8dty.ap-south-1.rds.amazonaws.com:3306)/task?charset=utf8&parseTime=True&loc=Local")
	if err != nil{
		panic(err)
	}
	DB=d
	d.AutoMigrate(&models.Login{})
}

func GetDB() *gorm.DB{
	return DB
}