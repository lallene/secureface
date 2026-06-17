package models

import "time"

type AccessLog struct {
	ID     uint  `gorm:"primaryKey" json:"id"`
	UserID *uint `json:"user_id"`
	DoorID uint  `gorm:"not null" json:"door_id"`

	User *User `gorm:"foreignKey:UserID" json:"user"`
	Door Door  `gorm:"foreignKey:DoorID" json:"door"`

	Status          string   `json:"status"` // GRANTED, DENIED, UNKNOWN_FACE
	Reason          string   `json:"reason"`
	SourceIP        string   `json:"source_ip"`
	ConfidenceScore *float64 `json:"confidence_score"`

	CreatedAt time.Time `json:"created_at"`
}
