package sqlite

import (
	"fmt"

	"github.com/iqrahadian/brick-web-scrapper-assesment/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewClient() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// cleanup database
	db.Exec("delete from products")

	err = db.AutoMigrate(&model.Product{})
	if err != nil {
		fmt.Println("ERROR MIGRATE")
		panic(err)
	}

	return db
}
