package sqlite

import (
	"time"
)

type Model struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"created_at"`
}

type Server struct {
	Model

	Name        string `gorm:"unique;"`
	TSAURL      string
	Description string
}

type Session struct {
	Model

	Server   Server `gorm:"unique;"`
	ServerID uint
	Active   bool
	Token    string
}
