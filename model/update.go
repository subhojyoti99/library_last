package model

type Update struct {
	Email  string `json:"email" binding:"required"`
	RoleID uint   `gorm:"not null" json:"role_id"`
}
