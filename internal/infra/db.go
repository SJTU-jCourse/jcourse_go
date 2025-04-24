package infra

import (
	"fmt"

	"github.com/glebarez/sqlite"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"jcourse_go/config"
)

var dbClient *gorm.DB

func GetDBClient() *gorm.DB {
	return dbClient
}

func initPostgresql(conf *config.DB) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		conf.Host, conf.User, conf.Password, conf.DBName, conf.Port)
	var err error
	dbClient, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return err
}

func InitDBClient(conf *config.DB) {
	err := initPostgresql(conf)
	if err != nil {
		panic(err)
	}
}

func InitTestMemDBClient() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	dbClient = db
}
