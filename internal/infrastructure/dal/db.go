package dal

import (
	"fmt"

	"github.com/glebarez/sqlite"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"jcourse_go/internal/config"
)

func NewPostgresSQL(c config.PostgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if c.Debug {
		db = db.Debug()
	}
	return db, nil
}

func NewSqlite(c config.SqliteConfig) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(c.Path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
