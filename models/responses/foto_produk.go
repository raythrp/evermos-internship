package responses

type FotoProduk struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	IDProduk  uint      `json:"product_id"`
	Url       string    `json:"url"`
}

func (FotoProduk) TableName() string {
	return "foto_produk" // Match your predefined table name exactly
}