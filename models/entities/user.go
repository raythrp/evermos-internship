package models

import "time"

type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Nama        string    `json:"nama"`
	KataSandi   string    `json:"kata_sandi"`
	NoTelp      string    `gorm:"unique" json:"notelp"`
	TanggalLahir time.Time `json:"tanggal_lahir"`
	JenisKelamin string    `json:"jenis_kelamin"`
	Tentang      string    `json:"tentang"`
	Pekerjaan    string    `json:"pekerjaan"`
	Email        string    `json:"email"`
	IDProvinsi   string    `json:"id_provinsi"`
	IDKota       string    `json:"id_kota"`
	IsAdmin      bool      `json:"is_admin"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
	Toko         Toko      `gorm:"foreignKey:IDUser" json:"toko"`
	Alamat       []Alamat  `gorm:"foreignKey:IDUser" json:"alamat"`
	Trx          []Trx     `gorm:"foreignKey:IDUser" json:"trx"`
}

func (User) TableName() string {
	return "user" // Match your predefined table name exactly
}