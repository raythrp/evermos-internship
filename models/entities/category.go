package models

import "time"

type Category struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	NamaCategory string    `json:"nama_category"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Produk       []Produk  `gorm:"foreignKey:IDCategory" json:"produk"`
	LogProduk    []LogProduk `gorm:"foreignKey:IDCategory" json:"log_produk"`
}

func (Category) TableName() string {
	return "category" // Match your predefined table name exactly
}
