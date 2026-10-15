package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex; not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`
	Tasks        []Task `gorm:"foreignKey:OwnerID" json:"-"`
}
