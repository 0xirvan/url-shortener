package model

type ShortURL struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	UserID      uint   `gorm:"not null"`
	OriginalURL string `gorm:"not null"`
	ShortCode   string `gorm:"not null;unique"`
	Clicks      uint   `gorm:"default:0"`
	CreatedAt   string `gorm:"autoCreateTime"`
	UpdatedAt   string `gorm:"autoUpdateTime"`
	User        *User  `gorm:"foreignKey:user_id;references:id"`
}
