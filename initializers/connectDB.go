package initializers

import (
	"fmt"
	"log"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func ConnectDB(config *Config) {
	var err error
	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", config.DBUserName, config.DBUserPassword, config.DBHost, config.DBPort, config.DBName)
	DB,err = gorm.Open(config.DBDriver,DBURL)
	if err != nil{
		fmt.Println("Cannot connect to database ", config.DBDriver)
		log.Fatal("connection error:", err)
	}else {
		fmt.Println("We are connected to the database ", config.DBDriver)
	}

}
