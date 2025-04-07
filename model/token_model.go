package model

import "time"

type Token struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	UserID     uint      `gorm:"not null"`
	Token      string    `gorm:"not null"`
	Type       string    `gorm:"not null"`
	ExpriredAt time.Time `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
	User       *User     `gorm:"foreignKey:user_id;references:id"`
}
