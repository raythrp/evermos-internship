package requests

type UpdateUser struct {
	Nama        string    `json:"nama"`
	KataSandi   string    `json:"kata_sandi"`
	NoTelp      string    `json:"no_telp"`
	TanggalLahir string   `json:"tanggal_lahir"`
	Tentang      string    `json:"tentang"`
	Pekerjaan    string    `json:"pekerjaan"`
	Email        string    `json:"email"`
	IDProvinsi   string    `json:"id_provinsi"`
	IDKota       string    `json:"id_kota"`
}

type CreateAlamat struct {
	JudulAlamat  string    `json:"judul_alamat"`
	NamaPenerima string    `json:"nama_penerima"`
	NoTelp       string    `json:"no_telp"`
	DetailAlamat string    `json:"detail_alamat"`
}

type UpdateAlamat struct {
	NamaPenerima string    `json:"nama_penerima"`
	NoTelp       string    `json:"no_telp"`
	DetailAlamat string    `json:"detail_alamat"`
}