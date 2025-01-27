package responses

type DetailTrx struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	IDTrx      uint      `json:"id_trx"`
	IDLogProduk uint     `json:"-"`
	LogProduk  LogProduk    `gorm:"foreignKey:IDLogProduk;references:ID" json:"product"`
	IDToko     uint      `json:"-"`
	Toko 	     Toko	   `gorm:"foreignKey:IDToko;references:ID" json:"toko"`
	Kuantitas  int       `json:"kuantitas"`
	HargaTotal int       `json:"harga_total"`
}

func (DetailTrx) TableName() string {
	return "detail_trx" // Match your predefined table name exactly
}