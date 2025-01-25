package responses

type Toko struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	IDUser    uint      `json:"id_user"`
	NamaToko  string    `json:"nama_toko"`
	UrlFoto   string    `json:"url_foto"`
}

func (Toko) TableName() string {
	return "toko" // Match your predefined table name exactly
}