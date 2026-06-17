package models

import "time"

type FaceProfile struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	UserID         uint       `gorm:"not null" json:"user_id"`
	User           User       `gorm:"foreignKey:UserID" json:"user"`
	FaceVector     string     `gorm:"type:text" json:"-"`
	ImagePath      string     `json:"image_path"`
	ImageExpiresAt *time.Time `json:"image_expires_at"`
	ImageDeletedAt *time.Time `json:"image_deleted_at"`
	IsVerified     bool       `gorm:"default:false" json:"is_verified"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
