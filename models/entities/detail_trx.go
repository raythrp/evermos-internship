package models

import "time"

type DetailTrx struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	IDTrx      uint      `json:"id_trx"`
	IDLogProduk uint     `json:"id_log_produk"`
	IDToko     uint      `json:"id_toko"`
	Kuantitas  int       `json:"kuantitas"`
	HargaTotal int       `json:"harga_total"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedAt  time.Time `json:"created_at"`
}

func (DetailTrx) TableName() string {
	return "detail_trx" // Match your predefined table name exactly
}