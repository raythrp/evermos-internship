package responses

type Produk struct {
	ID           uint         `gorm:"primaryKey" json:"id"`
	NamaProduk   string       `json:"nama_produk"`
	Slug         string       `json:"slug"`
	HargaReseller string      `json:"harga_reseller"`
	HargaKonsumen string      `json:"harga_konsumen"`
	Stok         int          `json:"stok"`
	Deskripsi    string       `json:"deskripsi"`
	IDToko       uint         `json:"-"`
	IDCategory   uint         `json:"-"`
	Toko       Toko         `gorm:"foreignKey:IDToko;references:ID" json:"toko"`
	Category   Category         `gorm:"foreignKey:IDCategory;references:ID" json:"category"`
	FotoProduk   []FotoProduk `gorm:"foreignKey:IDProduk;references:ID" json:"photos"`
}

func (Produk) TableName() string {
	return "produk" // Match your predefined table name exactly
}