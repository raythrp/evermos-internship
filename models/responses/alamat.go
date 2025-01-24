package responses


type Alamat struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	JudulAlamat  string    `json:"judul_alamat"`
	NamaPenerima string    `json:"nama_penerima"`
	NoTelp       string    `json:"no_telp"`
	DetailAlamat string    `json:"detail_alamat"`
}

func (Alamat) TableName() string {
	return "alamat" // Match your predefined table name exactly
}