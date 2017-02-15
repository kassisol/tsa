package sqlite

import (
	"github.com/kassisol/tsa/storage"
	"github.com/kassisol/tsa/storage/driver"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func init() {
	storage.RegisterDriver("sqlite", New)
}

type Config struct {
	DB *gorm.DB
}

func New(dbFilePath string) (driver.Storager, error) {
	debug := false

	db, err := gorm.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, err
	}

	db.LogMode(debug)

	db.AutoMigrate(&ServerConfig{})

	return &Config{DB: db}, nil
}

func (c *Config) End() {
	c.DB.Close()
}
