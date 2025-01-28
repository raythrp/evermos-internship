package requests

type CreateTrx struct {
	MetodeBayar string `json:"method_bayar"`
	AlamatKirim uint `json:"alamat_kirim"` 
	DetailTrx []CreateDetailTrx `json:"detail_trx"`
}

type CreateDetailTrx struct {
	ProdukID uint `json:"product_id"`
	Kuantitas int `json:"kuantitas"`
}