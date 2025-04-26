package model

import "time"

type Account struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	UserID        int64     `gorm:"not null" json:"user_id"`
	AccountNumber string    `gorm:"size:15;unique;not null" json:"account_number"`
	Balance       float64   `gorm:"default:0.00" json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
