package models

import "time"

type Alamat struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	IDUser       uint      `json:"id_user"`
	JudulAlamat  string    `json:"judul_alamat"`
	NamaPenerima string    `json:"nama_penerima"`
	NoTelp       string    `json:"no_telp"`
	DetailAlamat string    `json:"detail_alamat"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
	Trx          []Trx     `gorm:"foreignKey:AlamatPengiriman" json:"trx"`
}

func (Alamat) TableName() string {
	return "alamat" // Match your predefined table name exactly
}