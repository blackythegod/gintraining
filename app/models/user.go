package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email    string    `json:"email" gorm:"unique"`
	Username string    `json:"username" gorm:"unique"`
	Password string    `json:"password" gorm:""`
}
type Login struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}
