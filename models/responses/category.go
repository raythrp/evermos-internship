package responses

type Category struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	NamaCategory string    `json:"nama_category"`

}

func (Category) TableName() string {
	return "category" // Match your predefined table name exactly
}
