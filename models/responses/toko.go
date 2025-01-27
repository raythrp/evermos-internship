package responses

type MyToko struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	IDUser    uint      `json:"id_user"`
	NamaToko  string    `json:"nama_toko"`
	UrlFoto   string    `json:"url_foto"`
}

func (MyToko) TableName() string {
	return "toko" // Match your predefined table name exactly
}

type Toko struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	NamaToko  string    `json:"nama_toko"`
	UrlFoto   string    `json:"url_foto"`
}

func (Toko) TableName() string {
	return "toko" // Match your predefined table name exactly
}

type TokoAtLogProduk struct {
	ID        uint      `gorm:"primaryKey" json:"-"`
	NamaToko  string    `json:"nama_toko"`
	UrlFoto   string    `json:"url_foto"`
}

func (TokoAtLogProduk) TableName() string {
	return "toko" // Match your predefined table name exactly
}