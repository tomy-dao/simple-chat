package model

import (
	"time"
)

type User struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserName    string         `json:"username" gorm:"column:username;unique;not null"`
	Password    string        `json:"-" gorm:"column:password;not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}
