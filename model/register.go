package model

type Register struct {
	ID            uint   `gorm:"primary_key;auto_increment"`
	Name          string `json:"name"`
	Email         string `json:"email" binding:"required"`
	ContactNumber string `json:"contact_number" binding:"required"`
	Role          string `json:"role" binding:"required"`
	Password      string `json:"password" binding:"required"`
}
