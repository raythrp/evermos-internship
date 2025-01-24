package requests

type Login struct {
	KataSandi   string    `json:"kata_sandi"`
	NoTelp      string    `json:"no_telp"`
}

type Register struct {
	Nama        string    `json:"nama"`
	KataSandi   string    `json:"kata_sandi"`
	NoTelp      string    `json:"no_telp"`
	TanggalLahir string `json:"tanggal_lahir"`
	Tentang      string    `json:"tentang"`
	Pekerjaan    string    `json:"pekerjaan"`
	Email        string    `json:"email"`
	IDProvinsi   string    `json:"id_provinsi"`
	IDKota       string    `json:"id_kota"`
}