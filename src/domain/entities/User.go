package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UUID     string `json:"uuid" gorm:"uuid"`
	Username string `json:"username" gorm:"username"`
	Password string `gorm:"password"`
}
