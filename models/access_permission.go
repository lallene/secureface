package models

import "time"

type AccessPermission struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"not null" json:"user_id"`
	DoorID uint `gorm:"not null" json:"door_id"`

	User User `gorm:"foreignKey:UserID" json:"user"`
	Door Door `gorm:"foreignKey:DoorID" json:"door"`

	CanAccess bool   `gorm:"default:true" json:"can_access"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
