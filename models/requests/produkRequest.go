package requests

import "time"

type CreateProduk struct {
	NamaProduk   string       `json:"nama_produk"`
	HargaReseller string      `json:"harga_reseller"`
	HargaKonsumen string      `json:"harga_konsumen"`
	Stok         int          `json:"stok"`
	Deskripsi    string       `json:"deskripsi"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	IDToko       uint         `json:"id_toko"`
	IDCategory   uint         `json:"category_id"`
}

func (CreateProduk) TableName() string {
	return "produk" // Match your predefined table name exactly
}