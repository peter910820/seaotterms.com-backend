package model

import (
	"time"

	"github.com/lib/pq"
)

// blog schema
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

/* --------------------------------- */
/* --------------------------------- */

// galgame brand record schema
type BrandRecord struct {
	Brand       string    `gorm:"primaryKey" json:"brand"`          // PK
	Completed   int       `gorm:"not null" json:"completed"`        // Completed game amount
	Total       int       `gorm:"not null" json:"total"`            // Total game amount
	Annotation  string    `gorm:"not null" json:"annotation"`       // Annotation
	Dissolution bool      `gorm:"default:false" json:"dissolution"` // Dissolution
	InputTime   time.Time `gorm:"autoCreateTime" json:"inputTime"`  // InputTime
	InputName   string    `gorm:"not null" json:"inputName"`        // InputName
	UpdateTime  time.Time `gorm:"autoUpdateTime" json:"updateTime"` // UpdateTime
	UpdateName  string    `gorm:"not null" json:"updateName"`       // UpdateName
}

// galgame game record schema
type GameRecord struct {
	Name        string    `gorm:"primaryKey" json:"name"`           // PK
	Brand       string    `gorm:"not null" json:"brand"`            // Brand
	ReleaseDate time.Time `gorm:"not null" json:"releaseDate"`      // ReleaseDate
	AllAges     bool      `gorm:"not null" json:"allAges"`          // For all ages
	EndDate     time.Time `gorm:"not null" json:"endDate"`          // End date of play
	InputTime   time.Time `gorm:"autoCreateTime" json:"inputTime"`  // InputTime
	InputName   string    `gorm:"not null" json:"inputName"`        // InputName
	UpdateTime  time.Time `gorm:"autoUpdateTime" json:"updateTime"` // UpdateTime
	UpdateName  string    `gorm:"not null" json:"updateName"`       // UpdateName
}
