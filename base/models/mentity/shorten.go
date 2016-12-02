package mentity

// Redirect redirect entity mapping
type Redirect struct {
	ID          uint64 `gorm:"primary_key;column:redirect_id"`
	LongURL     string `gorm:"column:long_url"`
	ShortURL    string `gorm:"column:short_url"`
	LongCrc     uint64 `gorm:"column:long_crc"`
	ShortCrc    uint64 `gorm:"column:short_crc"`
	Status      uint8  `gorm:"column:status"`
	CreatedByIP uint64 `gorm:"column:created_by_ip"`
	UpdateByIP  uint64 `gorm:"column:updated_by_ip"`
	CreateAT    uint64 `gorm:"column:created_at"`
	UpdateAT    uint64 `gorm:"column:updated_at"`
}
