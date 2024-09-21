package main

import (
	"github.com/joho/godotenv"

	"jcourse_go/dal"
	"jcourse_go/repository"
)

func main() {
	_ = godotenv.Load()
	dal.InitDBClient()
	db := dal.GetDBClient()
	err := repository.Migrate(db)
	if err != nil {
		panic(err)
	}
}
