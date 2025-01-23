package models

import "time"

type Produk struct {
	ID           uint         `gorm:"primaryKey" json:"id"`
	NamaProduk   string       `json:"nama_produk"`
	Slug         string       `json:"slug"`
	HargaReseller string      `json:"harga_reseller"`
	HargaKonsumen string      `json:"harga_konsumen"`
	Stok         int          `json:"stok"`
	Deskripsi    string       `json:"deskripsi"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	IDToko       uint         `json:"id_toko"`
	IDCategory   uint         `json:"id_category"`
	LogProduk    []LogProduk  `gorm:"foreignKey:IDProduk" json:"log_produk"`
	FotoProduk   []FotoProduk `gorm:"foreignKey:IDProduk" json:"foto_produk"`
}

func (Produk) TableName() string {
	return "produk" // Match your predefined table name exactly
}