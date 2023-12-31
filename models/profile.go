package models

import "time"

type User struct {
	ID        int     `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Username  string  ` json:"username" binding:"required"`
	Email     string  `gorm:"uniqueIndex" json:"email" binding:"required"`
	Password  string  ` json:"password" binding:"required"`
	Photos    []Photo `gorm:"foreignKey:UserID"`
	Token     string
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Photo struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `gorm:"column:photoUrl" json:"photoUrl"`
	UserID   int    `gorm:"index" json:"user_id"`
}
