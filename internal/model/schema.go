package model

import (
	"time"
)

// seaotterms_com

type Account struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"NOT NULL unique"`
	Password  string    `gorm:"NOT NULL"`
	Email     string    `gorm:"NOT NULL unique"`
	CreatedAt time.Time // `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	UpdatedAt time.Time // `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP"`
}

// type Article struct {
// 	ID        uint           `gorm:"primaryKey"`
// 	Title     string         `gorm:"NOT NULL"`
// 	Username  string         `gorm:"NOT NULL"`
// 	Tags      pq.StringArray `gorm:"type:text[]"`
// 	Content   string         `gorm:"NOT NULL"`
// 	CreatedAt time.Time      // `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP"`
// 	UpdatedAt time.Time      // `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP"`
// }

// type Tag struct {
// 	ID   uint   `gorm:"primaryKey"`
// 	Name string `gorm:"unique;NOT NULL"`
// }

/* --------------------------------- */
/* --------------------------------- */

// A00_Blog

// all user info of the blog
type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Username   string    `gorm:"NOT NULL unique" json:"username"`
	Password   string    `gorm:"NOT NULL" json:"-"`
	Email      string    `gorm:"NOT NULL unique" json:"email"`
	Avatar     string    `gorm:"NOT NULL; default:''" json:"avatar"`
	Exp        int       `gorm:"default:0" json:"exp"`
	Management bool      `gorm:"default:false" json:"management"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
	CreateName string    `gorm:"NOT NULL" json:"createName"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	UpdateName string    `json:"updateName"`
}

// all article tags
// 不允許修改Name(PK)
type Tag struct {
	Name      string    `gorm:"primaryKey" json:"name"`
	IconName  string    `json:"iconName"`
	CreatedAt time.Time `gorm:"NOT NULL; autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// all article
type Article struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"NOT NULL" json:"title"`
	Content   string    `gorm:"NOT NULL" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Tags      []Tag     `gorm:"many2many:article_tags" json:"tags"`
}

// all user todos of the blog
type Todo struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	Owner      string     `gorm:"NOT NULL" json:"owner"`
	Topic      string     `gorm:"NOT NULL" json:"topic"`
	Title      string     `gorm:"NOT NULL" json:"title"`
	Status     uint       `gorm:"NOT NULL" json:"status"`
	Deadline   *time.Time `json:"deadline"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	CreateName string     `gorm:"NOT NULL" json:"createName"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	UpdateName string     `json:"updateName"`
}

// all todo topics of the blog
type TodoTopic struct {
	TopicName  string    `gorm:"primaryKey" json:"topicName"`
	TopicOwner string    `gorm:"primaryKey; default:'root'" json:"topicOwner"`
	UpdatedAt  time.Time `gorm:"autoCreateTime" json:"updatedAt"`
	UpdateName string    `json:"updateName"`
}

// An updated todo of all my public systems
type SystemTodo struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	SystemName  string     `gorm:"NOT NULL" json:"systemName"`
	Title       string     `gorm:"NOT NULL" json:"title"`
	Detail      string     `gorm:"NOT NULL" json:"detail"`
	Status      uint       `gorm:"NOT NULL" json:"status"`
	Deadline    *time.Time `json:"deadline"`
	Urgency     uint       `gorm:"NOT NULL" json:"urgency"`
	CreatedAt   time.Time  `gorm:"NOT NULL; autoCreateTime" json:"createdAt"`
	CreatedName string     `gorm:"NOT NULL" json:"createdName"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	UpdatedName string     `json:"updatedName"`
}

/* --------------------------------- */
/* --------------------------------- */

// A00_Galgame

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
