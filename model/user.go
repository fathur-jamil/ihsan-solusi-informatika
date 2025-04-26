package model

import "time"

type User struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	NIK         string    `gorm:"size:16;unique;not null" json:"nik"`
	PhoneNumber string    `gorm:"size:15;not null" json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
