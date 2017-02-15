package sqlite

import (
	"path"

	"github.com/juliengk/go-cert/ca/database"
	"github.com/juliengk/go-cert/ca/database/backend"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func init() {
	database.RegisterBackend("sqlite", New)
}

type Config struct {
	DB *gorm.DB
}

func New(config string) (backend.Backender, error) {
	debug := false

	file := path.Join(config, "index.db")

	db, err := gorm.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	db.LogMode(debug)

	db.AutoMigrate(&Certificate{})

	return &Config{DB: db}, nil
}

func (c *Config) End() {
	c.DB.Close()
}
