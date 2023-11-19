package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	ProductUrl  string
	Description string
	ImageUrl    string
	Price       float64
	Rating      float64
	Merchant    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	p.ID = uuid.NewString()
	return
}
