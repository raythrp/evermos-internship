package responses

type Trx struct {
	ID              uint         `gorm:"primaryKey" json:"id"`
	IDAlamatPengiriman uint        `gorm:"column:alamat_pengiriman" json:"-"`
	AlamatPengiriman Alamat        `gorm:"foreignKey:IDAlamatPengiriman;references:ID" json:"alamat_pengiriman"`
	HargaTotal      int          `json:"harga_total"`
	KodeInvoice     string       `json:"kode_invoice"`
	MethodBayar     string       `json:"method_bayar"`
	DetailTrx       []DetailTrx  `gorm:"foreignKey:IDTrx;references:ID" json:"detail_trx"`
}

func (Trx) TableName() string {
	return "trx" // Match your predefined table name exactly
}