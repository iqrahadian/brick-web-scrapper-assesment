package model

import "time"

type Product struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	ProductUrl   string
	Description  string
	ImageUrl     string
	Price        float64
	Rating       float32
	MerchantName string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// type Model struct {
// 	ID        uint
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// 	DeletedAt gorm.DeletedAt `gorm:"index"`
// }
