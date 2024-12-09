package model

import (
	"time"

	"github.com/lib/pq"
)

type Account struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"NOT NULL unique"`
	Password  string    `gorm:"NOT NULL"`
	Email     string    `gorm:"NOT NULL unique"`
	CreatedAt time.Time // `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	UpdatedAt time.Time // `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP"`
}

type Article struct {
	ID        uint           `gorm:"primaryKey"`
	Title     string         `gorm:"NOT NULL"`
	Username  string         `gorm:"NOT NULL"`
	Tags      pq.StringArray `gorm:"type:text[]"`
	Content   string         `gorm:"NOT NULL"`
	CreatedAt time.Time      // `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      // `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP"`
}

type Tag struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;NOT NULL"`
}
