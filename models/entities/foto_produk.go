package models

import "time"

type FotoProduk struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	IDProduk  uint      `json:"id_produk"`
	Url       string    `json:"url"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (FotoProduk) TableName() string {
	return "foto_produk" // Match your predefined table name exactly
}