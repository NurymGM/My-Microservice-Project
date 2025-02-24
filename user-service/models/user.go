package models

import "gorm.io/gorm"

type Userr struct {
    gorm.Model
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" gorm:"unique;not null" binding:"required,email"`
    Password string `json:"password" gorm:"not null" binding:"required"`
    Role     string `json:"role" gorm:"not null" binding:"required"` 
}
