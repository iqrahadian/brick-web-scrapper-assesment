package sqlite

import (
	"fmt"

	"github.com/iqrahadian/brick-web-scrapper-assesment/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	fmt.Println("INIT DB")
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("MIGRATE")
	err = db.AutoMigrate(&model.Product{})
	if err != nil {
		fmt.Println("ERROR MIGRATE")
		panic(err)
	}

	return db
}
