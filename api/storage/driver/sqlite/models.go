package sqlite

import (
	"time"
)

type Model struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"created_at"`
}

type ServerConfig struct {
	Model

	Key   string
	Value string
}
