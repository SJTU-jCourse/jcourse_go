package dal

import (
	"fmt"

	"jcourse_go/util"

	"github.com/glebarez/sqlite"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbClient *gorm.DB

func GetDBClient() *gorm.DB {
	return dbClient
}

func initPostgresql() error {
	host := util.GetPostgresHost()
	port := util.GetPostgresPort()
	user := util.GetPostgresUser()
	password := util.GetPostgresPassword()
	dbname := util.GetPostgresDBName()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, password, dbname, port)
	var err error
	dbClient, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return err
}

func InitDBClient() {
	err := initPostgresql()
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
