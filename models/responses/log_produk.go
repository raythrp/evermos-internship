package responses

type LogProduk struct {
	ID           uint         `gorm:"primaryKey" json:"-"`
	IDProduk      uint      `json:"id"`
	NamaProduk   string       `json:"nama_produk"`
	Slug         string       `json:"slug"`
	HargaReseller string      `json:"harga_reseller"`
	HargaKonsumen string      `json:"harga_konsumen"`
	Stok         int          `json:"stok"`
	Deskripsi    string       `json:"deskripsi"`
	IDToko       uint         `json:"-"`
	IDCategory   uint         `json:"-"`
	Toko       TokoAtLogProduk    `gorm:"foreignKey:IDToko;references:ID" json:"toko"`
	Category   Category         `gorm:"foreignKey:IDCategory;references:ID" json:"category"`
	FotoProduk   []FotoProduk `gorm:"foreignKey:IDProduk;references:ID" json:"photos"`
}

func (LogProduk) TableName() string {
	return "log_produk" // Match your predefined table name exactly
}