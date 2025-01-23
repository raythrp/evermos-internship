package models

import "time"

type Trx struct {
	ID              uint         `gorm:"primaryKey" json:"id"`
	IDUser          uint         `json:"id_user"`
	AlamatPengiriman uint        `json:"alamat_pengiriman"`
	HargaTotal      int          `json:"harga_total"`
	KodeInvoice     string       `json:"kode_invoice"`
	MethodBayar     string       `json:"method_bayar"`
	UpdatedAt       time.Time    `json:"updated_at"`
	CreatedAt       time.Time    `json:"created_at"`
	DetailTrx       []DetailTrx  `gorm:"foreignKey:IDTrx" json:"detail_trx"`
}

func (Trx) TableName() string {
	return "trx" // Match your predefined table name exactly
}