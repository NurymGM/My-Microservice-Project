package models

import "gorm.io/gorm"

type Userr struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
}