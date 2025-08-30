package main

import (
	"github.com/joho/godotenv"

	"jcourse_go/internal/dal"
)

func main() {
	_ = godotenv.Load()
	dal.InitDBClient()
	db := dal.GetDBClient()
	err := dal.Migrate(db)
	if err != nil {
		panic(err)
	}
}
