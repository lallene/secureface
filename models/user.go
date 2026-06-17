package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `gorm:"unique" json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
