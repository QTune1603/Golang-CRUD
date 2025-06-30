package model

import (
	
)

type User struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Username  string `json:"username" gorm:"unique;not null"`
	Password  string `json:"-"` // Lưu password đã hash
	CreatedAt int64  `json:"created_at"`
}
