package models

import "time"

type Toko struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	IDUser    uint      `json:"id_user"`
	NamaToko  string    `json:"nama_toko"`
	UrlFoto   string    `json:"url_foto"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	Produk    []Produk  `gorm:"foreignKey:IDToko" json:"produk"`
	LogProduk []LogProduk `gorm:"foreignKey:IDToko" json:"log_produk"`
	DetailTrx []DetailTrx `gorm:"foreignKey:IDToko" json:"detail_trx"`
}

func (Toko) TableName() string {
	return "toko" // Match your predefined table name exactly
}
