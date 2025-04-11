package model

import "time"

type User struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"not null" json:"name"`
	Email         string    `gorm:"unique;not null" json:"email"`
	Password      string    `gorm:"not null" json:"-"`
	Role          string    `gorm:"default:'user';not null" json:"role"`
	VerifiedEmail bool      `gorm:"default:false;not null" json:"verified_email"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"-"`
	Tokens        []Token   `gorm:"foreignKey:user_id;references:id" json:"-"`
}
