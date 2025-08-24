package entities

import (
	"gorm.io/gorm"
)

type UserRole string

const (
	AdminRole     UserRole = "admin"
	MemberRole    UserRole = "member"
	ModeratorRole UserRole = "moderator"
)

type User struct {
	gorm.Model
	Username     string   `gorm:"uniqueIndex;size:100" json:"username"`
	Email        string   `gorm:"uniqueIndex;size:100" json:"email"`
	Password     string   `json:"password"`
	Role         UserRole `gorm:"type:enum('admin', 'member', 'moderator');default:'member'" json:"role"`
	Birthdate    string   `gorm:"type:date" json:"birthdate"`
	IsFirstLogin bool     `gorm:"default:false" json:"is_first_login"`
}
