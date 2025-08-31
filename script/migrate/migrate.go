package main

import (
	"github.com/joho/godotenv"

	"jcourse_go/internal/dal"

	dal2 "jcourse_go/internal/infrastructure/dal"
)

func main() {
	_ = godotenv.Load()
	dal.InitDBClient()
	db := dal.GetDBClient()
	err := dal2.Migrate(db)
	if err != nil {
		panic(err)
	}
}
