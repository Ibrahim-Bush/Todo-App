package model

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"`
	Role     string `json:"role" gorm:"default:'user';not null"`
}

type Auth_request struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
